// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package ca

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRootCertificate(t *testing.T) {
	priKey := GenerateRsaPrivateKey()
	pb := PEMEncoding(x509.MarshalPKCS1PrivateKey(priKey), BlocTypePrivateKey)
	rootCertificate, err := GenerateRootCertificate(priKey, pkix.Name{
		Country:            []string{"中国"},          // 国家
		Province:           []string{"北京市"},         // 省份
		Locality:           []string{"东城区"},         // 城市
		Organization:       []string{"NonEntiaLtd"}, // 证书持有者组织名称
		OrganizationalUnit: []string{"NonEntiaCA"},  // 证书持有者组织唯一标识
		CommonName:         "NonEntia Root CA",      // 证书持有者通用名，需保持唯一，否则验证会失败
	})
	require.NoError(t, err)
	require.Equal(t, len(rootCertificate), 1514)

	root := PEMEncoding(rootCertificate, BlocTypeCertificate)
	require.Equal(t, len(root), 2106)

	var csr []byte
	subject := pkix.Name{
		Country:            []string{"中国"},                 // 国家
		Province:           []string{"北京市"},                // 省份
		Locality:           []string{"东城区"},                // 城市
		Organization:       []string{"张三"},                 // 证书持有者组织名称
		OrganizationalUnit: []string{"110101199008087218"}, // 证书持有者组织唯一标识
		CommonName:         "Joash",                        // 证书持有者通用名，需保持唯一，否则验证会失败
	}
	csr, err = GenerateRequest(priKey, subject)
	require.NoError(t, err)

	var req *x509.CertificateRequest
	req, err = ParseCertificateRequest(csr)
	require.NoError(t, err)

	require.Equal(t, subject.Country, req.Subject.Country)
	require.Equal(t, subject.Province, req.Subject.Province)
	require.Equal(t, subject.Locality, req.Subject.Locality)
	require.Equal(t, subject.Organization, req.Subject.Organization)
	require.Equal(t, subject.OrganizationalUnit, req.Subject.OrganizationalUnit)
	require.Equal(t, subject.CommonName, req.Subject.CommonName)

	var p2 *rsa.PrivateKey
	p2, err = ParsePrivateKey(pb)
	require.NoError(t, err)
	require.Equal(t, priKey, p2)
}
