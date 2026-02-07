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
	url         string
	urlFallback string
	useFallback atomic.Bool
}

func NewUrlFailover(url, urlFallback string) *UrlFailover {
	return &UrlFailover{
		url:         url,
		urlFallback: urlFallback,
	}
}

// ChangeCurrent 切换当前使用的 URL
func (uf *UrlFailover) ChangeCurrent() {
	// 切换到另一个 URL
	uf.useFallback.Store(!uf.useFallback.Load())
}

// GetURL 获取当前使用的 URL
func (uf *UrlFailover) GetURL() string {
	if uf.useFallback.Load() {
		return uf.urlFallback
	}
	return uf.url
}

// GetPreviousUrl 获取上一个 URL（当前 URL 的反向）
func (uf *UrlFailover) GetPreviousUrl() string {
	if uf.useFallback.Load() {
		return uf.url
	}
	return uf.urlFallback
}

func (uf *UrlFailover) Next() (string, error) {
	return uf.GetURL(), nil
}

func (uf *UrlFailover) Feedback(_ *resty.RequestFeedback) {
}

func (uf *UrlFailover) Close() error {
	return nil
}

func createRestyClient(failover *UrlFailover) *resty.Client {
	return resty.New().
		SetBaseURL(failover.GetURL()).
		SetTimeout(3 * time.Second).
		SetRetryCount(1).
		SetLoadBalancer(failover).
		AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
			zap.L().Info("准备发送请求", zap.String("baseUrl", c.BaseURL()), zap.String("uri", req.URL))
			return nil
		}).
		AddRetryHooks(func(res *resty.Response, err error) {
			var url string
			var statusCode int
			var hasRawResponse bool
			if res != nil {
				if res.Request != nil {
					url = res.Request.URL
				}
				statusCode = res.StatusCode()
				hasRawResponse = res.RawResponse != nil
			}
			zap.L().Info("触发重试",
				zap.String("originalURL", url),
				zap.Error(err),
				zap.Bool("hasError", err != nil),
				zap.Bool("hasRes", res != nil),
				zap.Bool("hasRequest", res != nil && res.Request != nil),
				zap.Int("statusCode", statusCode),
				zap.Bool("hasRawResponse", hasRawResponse))

			// 当没有 RawResponse 时，说明是网络错误（连接失败等）
			if res != nil && res.RawResponse == nil {
				failover.ChangeCurrent()
				zap.L().Info("触发切换URL", zap.String("newURL", failover.GetURL()), zap.String("previousURL", failover.GetPreviousUrl()))
			}
		}).
		AddRetryConditions(func(r *resty.Response, err error) bool {
			return err != nil || r == nil || (r != nil && r.IsError())
		}).
		AddResponseMiddleware(func(c *resty.Client, res *resty.Response) error {
			zap.L().Info("收到响应", zap.String("baseUrl", c.BaseURL()), zap.String("uri", res.Request.URL), zap.ByteString("response", res.Bytes()))
			return nil
		})
}

func (s *Snca) request() *resty.Client {
	return s.client
}
