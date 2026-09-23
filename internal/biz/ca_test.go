// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-12, by liasica

package biz

import (
	"testing"

	"github.com/stretchr/testify/require"

	"auroraride.com/edocseal/internal/g"
)

func TestAgencyIssueCertificate(t *testing.T) {
	g.LoadConfig("config/config.yaml")

	crt, key, err := agencyIssueCertificate("张三", "北京市", "北京市", "address", "17719646710", "110101199003070000")
	require.NoError(t, err)

	t.Logf("crt: %s", crt)
	t.Logf("key: %s", key)
}

func TestRequestFromUrl(t *testing.T) {
	crt, key, err := requestFromUrl("http://localhost:5000/enterprise/cert/test", "91610103MACUXL1W1A")
	require.NoError(t, err)
	t.Logf("cert: %x", crt)
	t.Logf("key: %x", key)
}
