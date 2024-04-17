// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-17, by liasica

package snca

const (
	UrlApplyServiceRandom = "/sealCertService/com/api/cert/applyServiceRandom"
	UrlBusinessDataFinish = "/sealCertService/com/api/cert/businessDataFinish"
	UrlApplySealCert      = "/sealCertService/com/api/cert/applySealCert"
)

// NewUrl 组装api url
func (s *Snca) NewUrl(path string) string {
	return s.url + path
}
