// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package g

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/liasica/edocseal/ca"
)

var (
	rootCrt *Certificate
	entCrt  *Certificate
)

// Certificate 证书配置
type Certificate struct {
	certificate *x509.Certificate
	privateKey  *rsa.PrivateKey
}

// NewCertificate 获取证书
func NewCertificate() (root *Certificate, ent *Certificate) {
	if rootCrt == nil {
		rootCrt = loadCertificate(cfg.RootCertificate)
	}
	if entCrt == nil {
		entCrt = loadCertificate(cfg.EnterpriseCertificate)
	}

	return rootCrt, entCrt
}

// 加载证书
func loadCertificate(path CertificatePath) *Certificate {
	if path.Certificate == "" || path.PrivateKey == "" {
		return nil
	}
	// 加载根证书和私钥
	priKey, _ := ca.LoadPrivateKeyFromFile(path.PrivateKey)
	rootCa, _ := ca.LoadCertificateFromFile(path.Certificate)
	return &Certificate{
		certificate: rootCa,
		privateKey:  priKey,
	}
}

// IsValid 判断证书是否有效
func (crt *Certificate) IsValid() bool {
	return crt != nil && crt.certificate != nil && crt.privateKey != nil
}

// GetCertificate 获取证书
func (crt *Certificate) GetCertificate() *x509.Certificate {
	if crt == nil || crt.certificate == nil || crt.privateKey == nil {
		return nil
	}
	return crt.certificate
}

// GetPrivateKey 获取证书私钥
func (crt *Certificate) GetPrivateKey() *rsa.PrivateKey {
	if crt == nil || crt.certificate == nil || crt.privateKey == nil {
		return nil
	}
	return crt.privateKey
}
