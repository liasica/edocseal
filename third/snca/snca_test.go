// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-17, by liasica

package snca

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/g"
)

func TestSnca_RequestCACert(t *testing.T) {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	g.LoadConfig("./config/config.yaml")

	priKey := ca.GenerateRsaPrivateKey()
	key, _ := x509.MarshalPKCS8PrivateKey(priKey)
	_ = ca.SaveToFile("certificates/test_key.pem", key, ca.BlocTypePrivateKey)

	name := "王五"

	// 生成CSR请求
	csr, err := ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: name,
	})
	require.NoError(t, err)

	var b []byte
	b, err = NewSnca(g.GetSnca()).RequestCACert(
		CertTypePersonal,
		csr,
		name,
		name,
		"18555555555",
		"110101199003070000",
		"陕西省",
		"西安市",
		"",
	)

	require.NoError(t, err)
	t.Log(b)

	_ = ca.SaveToFile("certificates/test.pem", b, ca.BlocTypeCertificate)
}

func TestSnca_RequestEnterpriseCACert(t *testing.T) {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	g.LoadConfig("./config/config.yaml")

	priKey := ca.GenerateRsaPrivateKey()
	key, _ := x509.MarshalPKCS8PrivateKey(priKey)
	_ = ca.SaveToFile("certificates/enterprise_test_key.pem", key, ca.BlocTypePrivateKey)

	name := "西安时光驹新能源科技有限公司"

	// 生成CSR请求
	csr, err := ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: name,
	})
	require.NoError(t, err)

	var b []byte
	b, err = NewSnca(g.GetSnca()).RequestCACert(
		CertTypeEnterprise,
		csr,
		name,
		"李四",
		"18555555555",
		"110101199003070000",
		"陕西省",
		"西安市",
		"91610133MA712R0U00",
	)

	require.NoError(t, err)
	t.Log(b)

	_ = ca.SaveToFile("certificates/enterprise_test_cert.pem", b, ca.BlocTypeCertificate)
}
