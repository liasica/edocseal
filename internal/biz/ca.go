// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"path/filepath"
	"time"

	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/ent"
	"github.com/liasica/edocseal/internal/ent/certification"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/third/snca"
)

func certificatePaths(idcard string) (keypath string, capath string) {
	return filepath.Join(g.GetCertificateDir(), idcard+"_key.pem"), filepath.Join(g.GetCertificateDir(), idcard+"_cert.pem")
}

// RequestCertificae 申请证书
func RequestCertificae(name, province, city, address, phone, idcard string) (cert *ent.Certification, err error) {
	cert = queryCertification(idcard)
	if cert != nil {
		return
	}

	var crt, key []byte
	if g.IsSelfSign() {
		crt, key, err = selfIssueCertificate(name, province, city, address, phone, idcard)
		if err != nil {
			return
		}
	} else {
		crt, key, err = agencyIssueCertificate(name, province, city, address, phone, idcard)
		if err != nil {
			return
		}
	}

	kp, cp := certificatePaths(idcard)

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
		OnConflictColumns(certification.FieldIDCardNumber).
		UpdateNewValues().
		Save(context.Background())
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
	crt, err = snca.NewSnca(g.GetSnca()).RequestCACert(
		snca.CertTypePersonal,
		name,
		name,
		phone,
		idcard,
		province,
		city,
		csr,
		"",
	)

	return
}
