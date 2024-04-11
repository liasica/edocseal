// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"crypto/x509/pkix"
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
		Country:            []string{"中国"},   // 国家
		Province:           []string{province}, // 省份
		Locality:           []string{city},     // 城市
		Organization:       []string{name},     // 证书持有者组织名称
		OrganizationalUnit: []string{idcard},   // 证书持有者组织唯一标识
		CommonName:         idcard,             // 证书持有者通用名，需保持唯一，否则验证会失败
	})
	if err != nil {
		return
	}
	fmt.Println(address, phone, idcard)
	return
}

// 机构签发证书
// TODO: 待实现
func agencyIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
	return
}
