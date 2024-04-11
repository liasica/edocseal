// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/g"
)

func RequestCertificae(path string, name, province, city, address, phone, idcard string) {

}

// 自签发证书
func selfIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
	// 获取根证书
	rootCrt, _ := g.NewCertificate()
	if rootCrt == nil {
		return nil, nil, edocseal.ErrRootCertificateNotFound
	}

	// 签发证书
	crt, key, _, err = ca.CreateInterCertificate(rootCrt.GetPrivateKey(), rootCrt.GetCertificate(), pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{province},
		Locality:           []string{city},
		Organization:       []string{name},
		OrganizationalUnit: []string{idcard},
		CommonName:         idcard,
	})
	if err != nil {
		return
	}
	fmt.Println(address, phone, idcard)
	return
}

// 机构签发证书
func agencyIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
	// 生成私钥
	priKey := ca.GenerateRsaPrivateKey()

	// 生成CSR请求
	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:      []string{"CN"},
		Organization: []string{name},
	})

	pkcs10 := base64.StdEncoding.EncodeToString(csr)
	fmt.Println(pkcs10)

	// TODO: 待实现请求机构签发证书，model放到model/ca.go中
	return
}
