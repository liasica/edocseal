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
	"os"
	"path/filepath"
	"strings"
	"time"

	"auroraride.com/edocseal/ca"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/certification"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/third/snca"
)

func CertificatePaths(idcard string) (keypath string, capath string) {
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

// RequestEnterpriseCertAndUpdateConfig 申请企业证书并更新配置
func RequestEnterpriseCertAndUpdateConfig() (err error) {
	cfg := g.GetEnterpriseConfig()

	priKey := ca.GenerateRsaPrivateKey()
	key, _ := x509.MarshalPKCS8PrivateKey(priKey)

	// 生成CSR请求
	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: cfg.Name,
	})
	if err != nil {
		return
	}

	// 签发证书
	var b []byte
	b, err = snca.NewSnca(g.GetSnca()).RequestCACert(
		snca.CertTypeEnterprise,
		csr,
		cfg.Name,
		cfg.PersonName,
		cfg.Phone,
		cfg.Idcard,
		cfg.Province,
		cfg.City,
		cfg.CreditCode,
	)
	if err != nil {
		return
	}

	// 读取证书
	var crt *x509.Certificate
	crt, err = x509.ParseCertificate(b)
	if err != nil {
		return
	}

	dir := filepath.Dir(g.GetConfigFile())

	// 私钥路径
	kf := filepath.Join(dir, fmt.Sprintf("%s_%d_key.pem", cfg.Name, crt.SerialNumber))
	err = ca.SaveToFile(kf, key, ca.BlocTypePrivateKey)
	if err != nil {
		return
	}

	// 证书路径
	cf := filepath.Join(dir, fmt.Sprintf("%s_%d_cert.pem", cfg.Name, crt.SerialNumber))
	err = ca.SaveToFile(cf, b, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	c, _ := os.ReadFile(g.GetConfigFile())
	str := string(c)
	str = strings.ReplaceAll(str, cfg.PrivateKey, kf)
	str = strings.ReplaceAll(str, cfg.Certificate, cf)

	g.UpdateEnterpriseConfig(kf, cf)

	return os.WriteFile(g.GetConfigFile(), []byte(str), os.ModePerm)
}
