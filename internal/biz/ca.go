// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"auroraride.com/edocseal/ca"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/certification"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/model"
	"auroraride.com/edocseal/third/snca"
)

func CertificatePaths(idcard string) (keypath string, capath string) {
	return filepath.Join(g.GetCertificateDir(), idcard+"_key.pem"), filepath.Join(g.GetCertificateDir(), idcard+"_cert.pem")
}

// ResolveCertificateIssuer 证书生成方式为空时按配置 selfSign 决定
func ResolveCertificateIssuer(issuer string) string {
	if issuer != "" {
		return issuer
	}

	if g.IsSelfSign() {
		return model.CertificateIssuerSelf
	}
	return model.CertificateIssuerSnca
}

// RequestCertificae 按证书生成方式申请个人证书
func RequestCertificae(
	issuer string,
	name string,
	province string,
	city string,
	address string,
	phone string,
	idcard string,
) (cert *ent.Certification, err error) {
	// 未过期的证书直接复用，不区分签发方式
	cert = queryCertification(idcard)
	if cert != nil {
		return
	}

	var crt, key []byte
	switch issuer {
	case model.CertificateIssuerSelf:
		crt, key, err = selfIssueCertificate(name, province, city, address, phone, idcard)
	case model.CertificateIssuerSnca:
		crt, key, err = agencyIssueCertificate(name, province, city, address, phone, idcard)
	default:
		err = fmt.Errorf("不支持的证书生成方式: %s", issuer)
	}
	if err != nil {
		return
	}

	kp, cp := CertificatePaths(idcard)

	// 保存私钥
	err = ca.SaveToFile(kp, key, ca.BlocTypePrivateKey)
	if err != nil {
		return
	}

	// 保存证书
	err = ca.SaveToFile(cp, crt, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	// 保存证书信息
	return ent.NewDatabase().Certification.Create().
		SetIDCardNumber(idcard).
		SetPrivatePath(kp).
		SetCertPath(cp).
		SetExpiresAt(time.Now().Add(time.Hour*24 - time.Minute*10)). // 防止证书过期，有效期减少10分钟
		SetIssuer(issuer).
		OnConflictColumns(certification.FieldIDCardNumber).
		UpdateNewValues().
		Save(context.Background())
}

// 自签发证书
func selfIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
	// 获取根证书
	rootCrt := g.LoadRootCertificate()
	if !rootCrt.IsValid() {
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

// 查询证书
func queryCertification(idcard string) *ent.Certification {
	cert, _ := ent.NewDatabase().Certification.Query().
		Where(certification.IDCardNumber(idcard), certification.ExpiresAtGT(time.Now())).
		First(context.Background())

	return cert
}

// 机构签发证书
func agencyIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
	// 生成私钥
	priKey := ca.GenerateRsaPrivateKey()
	key, _ = x509.MarshalPKCS8PrivateKey(priKey)

	// 生成CSR请求
	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: name,
	})
	if err != nil {
		return
	}

	// 申请证书
	crt, err = snca.New().RequestCACert(
		snca.CertTypePersonal,
		csr,
		name,
		name,
		phone,
		idcard,
		province,
		city,
		"",
	)

	return
}
