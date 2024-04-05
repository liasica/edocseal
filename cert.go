// Copyright (C) conseal. 2024-present.
//
// Created at 2024-04-03, by liasica

package edocseal

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

type CACertificate struct {
	PrivateKey string
	PublicKey  string
	Csr        string
}

func CreateCertificateAuthority(names pkix.Name, expiration time.Duration, size int) (*CACertificate, error) {
	// 生成密钥对
	keys, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, fmt.Errorf("unable to genarate private keys, error: %s", err)
	}

	// 创建CSR模板
	var csrTemplate = x509.CertificateRequest{
		Subject:            names,
		SignatureAlgorithm: x509.SHA512WithRSA,
	}

	// 生成CSR请求
	var csrCertificate []byte
	csrCertificate, err = x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keys)
	if err != nil {
		return nil, err
	}
	csr := base64.StdEncoding.EncodeToString(csrCertificate)

	// 生成序列号
	var serial *big.Int
	serial, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber:          serial,
		Subject:               names,
		NotBefore:             now.Add(-10 * time.Minute).UTC(),
		NotAfter:              now.Add(expiration).UTC(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	// 签发证书
	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &keys.PublicKey, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate, error: %s", err)
	}

	var request bytes.Buffer
	var privateKey bytes.Buffer
	if err = pem.Encode(&request, &pem.Block{Type: "CERTIFICATE", Bytes: certificate}); err != nil {
		return nil, err
	}
	if err = pem.Encode(&privateKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keys)}); err != nil {
		return nil, err
	}

	// var x bytes.Buffer
	// err = pem.Encode(&x, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: x509.MarshalPKCS1PublicKey(&keys.PublicKey)})
	// fmt.Println(x.String())

	_ = os.WriteFile("cert.cert", request.Bytes(), 0644)

	return &CACertificate{
		PrivateKey: privateKey.String(),
		PublicKey:  request.String(),
		Csr:        csr,
	}, nil
}
