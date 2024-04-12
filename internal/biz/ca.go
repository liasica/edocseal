// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
)

// RequestCertificae 申请证书
func RequestCertificae(paths *model.DocumentPaths, name, province, city, address, phone, idcard string) (err error) {
	var crt, key []byte
	if g.IsSelfSign() {
		crt, key, err = selfIssueCertificate(name, province, city, address, phone, idcard)
		if err != nil {
			return
		}
	}
	// 保存证书
	err = ca.SaveToFile(paths.Cert, crt, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	// 保存私钥
	return ca.SaveToFile(paths.Key, key, ca.BlocTypePrivateKey)
}

// 自签发证书
func selfIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
	// 获取根证书
	rootCrt := g.NewCertificate()
	if rootCrt == nil {
		return nil, nil, errors.New("根证书不存在")
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
