// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	mr "math/rand"
	"os"
	"time"
)

const DefaultPrivateKeyBits = 2048

type BlocType int

const (
	BlocTypePublicKey BlocType = iota
	BlocTypePrivateKey
	BlocTypeCertificateRequest
	BlocTypeCertificate
	BlocTypeX509CRL
)

// GenerateRsaPrivateKey 生成RSA私钥
func GenerateRsaPrivateKey() (key *rsa.PrivateKey) {
	key, _ = rsa.GenerateKey(rand.Reader, DefaultPrivateKeyBits)
	return
}

// GenerateRequest 生成CSR请求
func GenerateRequest(priKey *rsa.PrivateKey, subject pkix.Name) ([]byte, error) {
	return x509.CreateCertificateRequest(
		rand.Reader,
		&x509.CertificateRequest{
			Subject:            subject,
			SignatureAlgorithm: x509.SHA256WithRSA,
		},
		priKey,
	)
}

func ParseCertificateRequest(csr []byte) (caRequest *x509.CertificateRequest, err error) {
	return x509.ParseCertificateRequest(csr)
}

// GenerateRootCertificate 生成根证书
func GenerateRootCertificate(priKey *rsa.PrivateKey, subject pkix.Name) ([]byte, error) {
	pubKey := priKey.Public()
	serial, _ := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	ca := x509.Certificate{
		Subject:               subject,
		SerialNumber:          serial,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(99, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageContentCommitment |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDataEncipherment |
			x509.KeyUsageKeyAgreement |
			x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign,
	}
	return x509.CreateCertificate(rand.Reader, &ca, &ca, pubKey, priKey)
}

// PEMEncoding PEM编码
func PEMEncoding(key []byte, blockType BlocType) []byte {
	var keyType string
	switch blockType {
	case BlocTypePublicKey:
		keyType = "PUBLIC KEY"
	case BlocTypePrivateKey:
		keyType = "PRIVATE KEY"
	case BlocTypeCertificateRequest:
		keyType = "CERTIFICATE REQUEST"
	case BlocTypeCertificate:
		keyType = "CERTIFICATE"
	case BlocTypeX509CRL:
		keyType = "X509 CRL"
	}

	block := &pem.Block{
		Type:  keyType,
		Bytes: key,
	}
	return pem.EncodeToMemory(block)
}

// SaveToFile 保存到文件
func SaveToFile(path string, b []byte, blockType BlocType) error {
	// 进行PEM编码
	data := PEMEncoding(b, blockType)
	return os.WriteFile(path, data, 0644)
}

// ParsePrivateKey 解析私钥
func ParsePrivateKey(b []byte) (priKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(b)
	var k any
	k, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	priKey = k.(*rsa.PrivateKey)
	return
}

// LoadPrivateKeyFromFile 从文件加载私钥
func LoadPrivateKeyFromFile(path string) (priKey *rsa.PrivateKey, err error) {
	var b []byte
	b, err = os.ReadFile(path)
	if err != nil {
		return
	}
	return ParsePrivateKey(b)
}

// ParseCertificate 解析pem编码证书
func ParseCertificate(b []byte) (crt *x509.Certificate, err error) {
	block, _ := pem.Decode(b)
	return x509.ParseCertificate(block.Bytes)
}

// LoadCertificateFromFile 从pem文件加载证书
func LoadCertificateFromFile(path string) (crt *x509.Certificate, err error) {
	var b []byte
	b, err = os.ReadFile(path)
	if err != nil {
		return
	}
	return ParseCertificate(b)
}

// CreateInterCertificate 签发中级证书
func CreateInterCertificate(priKey *rsa.PrivateKey, ca *x509.Certificate, subject pkix.Name) (interCrt, interKey []byte, serial *big.Int, err error) {
	// 生成中级证书密钥对
	key := GenerateRsaPrivateKey()

	// 生成随机数
	serial = big.NewInt(mr.Int63())

	// 证书请求
	certificate := x509.Certificate{
		Subject:               subject,
		SerialNumber:          serial,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageContentCommitment |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDataEncipherment |
			x509.KeyUsageKeyAgreement |
			x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign,
	}

	// 使用根证书签发中级证书
	interCrt, err = x509.CreateCertificate(rand.Reader, &certificate, ca, key.Public(), priKey)
	if err != nil {
		return
	}

	interKey, _ = x509.MarshalPKCS8PrivateKey(key)
	return
}
