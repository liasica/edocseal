// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-22, by liasica

package snca

import (
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

type UrlFailover struct {
	urls    [2]string
	current atomic.Uint32
}

func NewUrlFailover(url, urlFallback string) *UrlFailover {
	return &UrlFailover{
		urls: [2]string{url, urlFallback},
	}
}

// Switch 切换到另一个 URL
func (uf *UrlFailover) Switch() {
	uf.current.Store(1 - uf.current.Load())
}

// Current 获取当前使用的 URL
func (uf *UrlFailover) Current() string {
	return uf.urls[uf.current.Load()]
}

// Previous 获取备用 URL
func (uf *UrlFailover) Previous() string {
	return uf.urls[1-uf.current.Load()]
}

func (uf *UrlFailover) Next() (string, error) {
	return uf.Current(), nil
}

func (uf *UrlFailover) Feedback(_ *resty.RequestFeedback) {
}

func (uf *UrlFailover) Close() error {
	return nil
}

func createRestyClient(failover *UrlFailover) *resty.Client {
	return createRestyClientWithTimeout(failover, 3*time.Second)
}

func createRestyClientWithTimeout(failover *UrlFailover, timeout time.Duration) *resty.Client {
	client := resty.New().
		SetBaseURL(failover.Current()).
		SetTimeout(timeout).
		SetLoadBalancer(failover)

	client.AddRequestMiddleware(func(_ *resty.Client, request *resty.Request) error {
		zap.L().Info("准备发送 SNCA 请求",
			zap.String("url", requestURL(request)))
		return nil
	})

	client.OnSuccess(func(_ *resty.Client, response *resty.Response) {
		zap.L().Info("收到 SNCA 响应",
			zap.String("url", requestURL(response.Request)),
			zap.Int("statusCode", response.StatusCode()),
			zap.ByteString("response", response.Bytes()))
	})

	client.OnError(func(request *resty.Request, err error) {
		zap.L().Error("SNCA 请求失败",
			zap.String("url", requestURL(request)),
			zap.Error(err))
	})

	return client
}

func (s *Snca) request() *resty.Client {
	return s.client
}

func requestURL(request *resty.Request) string {
	if request == nil {
		return ""
	}

	if request.RawRequest != nil && request.RawRequest.URL != nil {
		return request.RawRequest.URL.String()
	}

	return request.URL
}

func shouldSwitchRequestURL(response *resty.Response, err error) bool {
	if err == nil || response == nil {
		return false
	}

	return response.RawResponse == nil
}

func (s *Snca) executeRequestWithURLFailover(sendRequest func() (*resty.Response, error)) (response *resty.Response, err error) {
	response, err = sendRequest()
	if s.urlFailover == nil || !shouldSwitchRequestURL(response, err) {
		return
	}

	failedURL := ""
	if response != nil {
		failedURL = requestURL(response.Request)
	}

	previousURL := s.urlFailover.Current()
	s.urlFailover.Switch()

	zap.L().Warn("SNCA 请求发生传输错误，切换到备用 URL",
		zap.String("failedURL", failedURL),
		zap.String("previousURL", previousURL),
		zap.String("currentURL", s.urlFailover.Current()),
		zap.Error(err))

	response, err = sendRequest()
	return
}
