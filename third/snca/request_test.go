// Copyright (C) edocseal. 2026-present.
//
// Created at 2026-02-07, by liasica

package snca

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSncaApplyServiceRanSwitchesURLOnTimeout(t *testing.T) {
	var primaryRequestCount atomic.Int32
	primaryServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		primaryRequestCount.Add(1)
		require.Equal(t, http.MethodPost, request.Method)
		require.Equal(t, UrlApplyServiceRandom, request.URL.Path)

		time.Sleep(120 * time.Millisecond)
	}))
	defer primaryServer.Close()

	var fallbackRequestCount atomic.Int32
	fallbackServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fallbackRequestCount.Add(1)
		require.Equal(t, http.MethodPost, request.Method)
		require.Equal(t, UrlApplyServiceRandom, request.URL.Path)

		writer.Header().Set("Content-Type", "application/json")

		_, err := writer.Write([]byte(`{"resultCode":"0","resultCodeMsg":"ok","randomB":"random-b"}`))
		require.NoError(t, err)
	}))
	defer fallbackServer.Close()

	failover := NewUrlFailover(primaryServer.URL, fallbackServer.URL)
	sncaClient := &Snca{
		source:       "aurora",
		client:       createRestyClientWithTimeout(failover, 40*time.Millisecond),
		urlFailover:  failover,
		customerType: "ride",
	}

	randomB, err := sncaClient.ApplyServiceRan()
	require.NoError(t, err)
	require.Equal(t, "random-b", randomB)
	require.Equal(t, int32(1), primaryRequestCount.Load())
	require.Equal(t, int32(1), fallbackRequestCount.Load())
	require.Equal(t, fallbackServer.URL, failover.Current())
}

func TestSncaApplyServiceRanDoesNotSwitchURLOnBusinessError(t *testing.T) {
	var primaryRequestCount atomic.Int32
	primaryServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		primaryRequestCount.Add(1)
		require.Equal(t, http.MethodPost, request.Method)
		require.Equal(t, UrlApplyServiceRandom, request.URL.Path)

		writer.Header().Set("Content-Type", "application/json")

		_, err := writer.Write([]byte(`{"resultCode":"1","resultCodeMsg":"business failed"}`))
		require.NoError(t, err)
	}))
	defer primaryServer.Close()

	var fallbackRequestCount atomic.Int32
	fallbackServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fallbackRequestCount.Add(1)
		t.Fatal("fallback URL should not be called on business error")
	}))
	defer fallbackServer.Close()

	failover := NewUrlFailover(primaryServer.URL, fallbackServer.URL)
	sncaClient := &Snca{
		source:       "aurora",
		client:       createRestyClientWithTimeout(failover, 40*time.Millisecond),
		urlFailover:  failover,
		customerType: "ride",
	}

	randomB, err := sncaClient.ApplyServiceRan()
	require.Empty(t, randomB)
	require.EqualError(t, err, "business failed")
	require.Equal(t, int32(1), primaryRequestCount.Load())
	require.Equal(t, int32(0), fallbackRequestCount.Load())
	require.Equal(t, primaryServer.URL, failover.Current())
}
