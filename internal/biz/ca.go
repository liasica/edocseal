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

	"go.uber.org/zap"
	"resty.dev/v3"

	"auroraride.com/edocseal"
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

// requestFromSnca 从SNCA申请证书
// 返回私钥、证书内容和证书对象
func requestFromSnca() (keyBytes []byte, crtBytes []byte, err error) {
	cfg := g.GetEnterpriseConfig()

	priKey := ca.GenerateRsaPrivateKey()
	keyBytes, _ = x509.MarshalPKCS8PrivateKey(priKey)

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
	crtBytes, err = snca.New().RequestCACert(
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

	return
}

// requestFromUrl 从指定URL获取证书和私钥
func requestFromUrl(url string, name, serial string) (keyBytes []byte, crtBytes []byte, err error) {
	// 若URL为空，则直接返回未找到企业证书错误
	if url == "" {
		return nil, nil, edocseal.ErrEnterpriseCertificateNotFound
	}

	// 发送HTTP请求获取证书和私钥
	var result edocseal.EnterpriseCertificate
	var resp *resty.Response
	resp, err = resty.New().R().
		SetResult(&result).
		SetQueryParam("name", name).
		SetQueryParam("serial", serial).
		Get(url)
	if err != nil {
		return
	}

	if resp.IsSuccess() && (result.Key == "" || result.Cert == "") {
		return
	}

	// 解析返回的证书和私钥
	keyBytes = []byte(result.Key)
	crtBytes = []byte(result.Cert)
	return
}

// RequestEnterpriseCertAndUpdateConfig 申请企业证书并更新配置
func RequestEnterpriseCertAndUpdateConfig() (err error) {
	cfg := g.GetEnterpriseConfig()

	var (
		keyBytes []byte
		crtBytes []byte
		crt      *x509.Certificate
	)

	// 尝试从url中获取证书和私钥
	keyBytes, crtBytes, err = requestFromUrl(cfg.Url, cfg.Name, cfg.GetCertificate().SerialNumber.String())
	if err != nil {
		// 如果返回错误是企业证书未找到，则尝试从SNCA申请
		if !errors.Is(err, edocseal.ErrEnterpriseCertificateNotFound) {
			keyBytes, crtBytes, err = requestFromSnca()
		}
	}

	// 如果申请证书失败，则返回错误
	if err != nil {
		return
	}

	// 如果keyBytes和crtBytes都为nil，则表示无需更新
	if keyBytes == nil && crtBytes == nil {
		zap.L().Info("无需更新企业证书", zap.String("name", cfg.Name))
		return
	}

	// 读取证书
	crt, err = x509.ParseCertificate(crtBytes)
	if err != nil {
		return
	}

	zap.L().Info("申请企业证书成功", zap.String("name", cfg.Name), zap.Int64("serialNumber", crt.SerialNumber.Int64()), zap.Time("notBefore", crt.NotBefore), zap.Time("notAfter", crt.NotAfter))

	dir := filepath.Dir(g.GetConfigFile())

	// 私钥路径
	kf := filepath.Join(dir, fmt.Sprintf("%s_%d_key.pem", cfg.Name, crt.SerialNumber))
	err = ca.SaveToFile(kf, keyBytes, ca.BlocTypePrivateKey)
	if err != nil {
		return
	}

	// 证书路径
	cf := filepath.Join(dir, fmt.Sprintf("%s_%d_cert.pem", cfg.Name, crt.SerialNumber))
	err = ca.SaveToFile(cf, crtBytes, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	c, _ := os.ReadFile(g.GetConfigFile())
	str := string(c)
	str = strings.ReplaceAll(str, cfg.PrivateKey, kf)
	str = strings.ReplaceAll(str, cfg.Certificate, cf)

	priKey, _ := ca.ParsePrivateKey(keyBytes)

	g.UpdateEnterpriseConfig(kf, cf, crt, crtBytes, priKey, keyBytes)

	return os.WriteFile(g.GetConfigFile(), []byte(str), os.ModePerm)
}
