// Copyright (C) edocseal. 2025-present.
//
// Created at 2025-08-16, by liasica

package edocseal

type EnterpriseCertificate struct {
	Cert string `json:"cert"` // 企业证书内容
	Key  string `json:"key"`  // 企业私钥内容
}
