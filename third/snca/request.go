// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-22, by liasica

package snca

import (
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

func (s *Snca) request() *resty.Client {
	return resty.New().
		SetBaseURL(s.url).
		SetTimeout(3 * time.Second).
		AddRetryHooks(func(res *resty.Response, err error) {
			zap.L().Info("触发重试", zap.String("url", s.url), zap.Error(err), zap.Reflect("response", res))
			if err != nil {
				res.Request.URL = s.urlFallback
			}
		}).
		AddRetryConditions(func(r *resty.Response, err error) bool {
			return r.IsError() || err != nil
		})
}
