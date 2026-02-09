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
	client := resty.New().
		SetBaseURL(failover.Current()).
		SetTimeout(3 * time.Second).
		SetRetryCount(1).
		SetLoadBalancer(failover).
		AddRetryConditions(func(r *resty.Response, err error) bool {
			return err != nil || r == nil || (r != nil && r.IsError())
		})

	// 使用 AddRequestMiddleware 记录请求
	client.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		zap.L().Info("准备发送请求",
			zap.String("baseURL", c.BaseURL()),
			zap.String("uri", req.URL))
		return nil
	})

	// 使用 OnSuccess 钩子记录成功响应
	client.OnSuccess(func(c *resty.Client, res *resty.Response) {
		zap.L().Info("收到响应",
			zap.String("baseURL", c.BaseURL()),
			zap.String("uri", res.Request.URL),
			zap.Int("statusCode", res.StatusCode()),
			zap.ByteString("response", res.Bytes()))
	})

	// 使用 OnError 钩子记录错误
	client.OnError(func(req *resty.Request, err error) {
		zap.L().Error("请求失败",
			zap.String("url", req.URL),
			zap.Error(err))
	})

	// 添加重试钩子
	client.AddRetryHooks(func(res *resty.Response, err error) {
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
			failover.Switch()
			zap.L().Info("触发切换URL",
				zap.String("newURL", failover.Current()),
				zap.String("previousURL", failover.Previous()))
		}
	})

	return client
}

func (s *Snca) request() *resty.Client {
	return s.client
}
