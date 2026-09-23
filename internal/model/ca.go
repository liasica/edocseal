// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

import "crypto/x509/pkix"

// 证书生成方式
const (
	CertificateIssuerSelf = "SELF" // 自签
	CertificateIssuerSnca = "SNCA" // 陕西CA
)

// RootCertificateSubject 自签根证书主题，CN 需保持唯一，否则验证会失败
var RootCertificateSubject = pkix.Name{
	Country:            []string{"CN"},
	Province:           []string{"陕西省"},
	Locality:           []string{"西安市"},
	Organization:       []string{"陕西极光换电科技有限责任公司"},
	OrganizationalUnit: []string{"91610103MACUXL1W1A"},
	CommonName:         "陕西极光换电科技有限责任公司 Root CA",
}
