// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package g

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/liasica/edocseal/ca"
)

var certificate *Certificate

// Certificate 证书配置
type Certificate struct {
	certificate *x509.Certificate
	privateKey  *rsa.PrivateKey
}

// NewCertificate 初始化证书
func NewCertificate() *Certificate {
	if cfg.RootCertificate.Certificate == "" || cfg.RootCertificate.PrivateKey == "" {
		return nil
	}
	// 加载根证书和私钥
	priKey, _ := ca.LoadPrivateKeyFromFile(cfg.RootCertificate.PrivateKey)
	rootCa, _ := ca.LoadCertificateFromFile(cfg.RootCertificate.Certificate)
	certificate = &Certificate{
		certificate: rootCa,
		privateKey:  priKey,
	}
	return certificate
}

// GetRootCertificate 获取根证书
func GetRootCertificate() *x509.Certificate {
	if certificate == nil {
		return nil
	}
	return certificate.certificate
}

// GetRootPrivateKey 获取根证书私钥
func GetRootPrivateKey() *rsa.PrivateKey {
	if certificate == nil {
		return nil
	}
	return certificate.privateKey
}
