// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/ca"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/certification"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/model"
	"auroraride.com/edocseal/third/snca"
)

// enterpriseRenewMutex 防止定时任务与手动触发同时续期企业证书
var enterpriseRenewMutex sync.Mutex

func CertificatePaths(idcard string) (keypath string, capath string) {
	return filepath.Join(g.GetCertificateDir(), idcard+"_key.pem"), filepath.Join(g.GetCertificateDir(), idcard+"_cert.pem")
}

// RequestCertificae 按证书生成方式申请证书，issuer 为空时按服务配置
func RequestCertificae(
	issuer string,
	name string,
	province string,
	city string,
	address string,
	phone string,
	idcard string,
) (cert *ent.Certification, err error) {
	if issuer == "" {
		issuer = model.CertificateIssuerSnca
		if g.IsSelfSign() {
			issuer = model.CertificateIssuerSelf
		}
	}

	// 生成方式切换后不复用其他方式签发的证书
	cert = queryCertification(idcard, issuer)
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
func queryCertification(idcard, issuer string) *ent.Certification {
	cert, _ := ent.NewDatabase().Certification.Query().
		Where(
			certification.IDCardNumber(idcard),
			certification.Issuer(issuer),
			certification.ExpiresAtGT(time.Now()),
		).
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

// requestFromUrl 从企业证书地址获取证书和私钥，返回 DER 编码内容，证书与本机一致时均返回 nil
func requestFromUrl(url string, name, serial string) (keyBytes []byte, crtBytes []byte, err error) {
	// 若URL为空，则直接返回未找到企业证书错误
	if url == "" {
		return nil, nil, edocseal.ErrEnterpriseCertificateNotFound
	}

	client := resty.New().SetTimeout(30 * time.Second)
	defer func() {
		_ = client.Close()
	}()

	// 发送HTTP请求获取证书和私钥
	var (
		result edocseal.EnterpriseCertificate
		resp   *resty.Response
	)
	resp, err = client.R().
		SetResult(&result).
		SetQueryParam("name", name).
		SetQueryParam("serial", serial).
		Get(url)
	if err != nil {
		return
	}

	if !resp.IsSuccess() {
		err = fmt.Errorf("企业证书地址响应异常: %s", resp.Status())
		return
	}

	// 名称和序列号与本机一致时响应体为空
	if result.Key == "" && result.Cert == "" {
		return
	}

	keyBytes, err = decodePEM(result.Key)
	if err != nil {
		err = fmt.Errorf("企业私钥解析失败: %w", err)
		return
	}

	crtBytes, err = decodePEM(result.Cert)
	if err != nil {
		err = fmt.Errorf("企业证书解析失败: %w", err)
	}

	return
}

// decodePEM 返回 PEM 文本中第一个块的 DER 内容
func decodePEM(content string) (b []byte, err error) {
	block, _ := pem.Decode([]byte(content))
	if block == nil {
		err = errors.New("不是有效的 PEM 内容")
		return
	}

	b = block.Bytes
	return
}

// RenewEnterpriseCertificate 企业证书 7 天内到期时申请新证书并更新配置，返回是否已更换证书
func RenewEnterpriseCertificate() (renewed bool, err error) {
	if !enterpriseRenewMutex.TryLock() {
		err = errors.New("企业证书续期正在执行")
		return
	}
	defer enterpriseRenewMutex.Unlock()

	cert := g.GetEnterpriseConfig().GetCertificate()
	zap.L().Info("检查企业证书有效期",
		zap.String("subject", cert.Subject.String()),
		zap.String("serialNumber", cert.SerialNumber.String()),
		zap.Time("notBefore", cert.NotBefore),
		zap.Time("notAfter", cert.NotAfter),
	)

	if !cert.NotAfter.Before(time.Now().AddDate(0, 0, 7)) {
		return
	}

	return RequestEnterpriseCertAndUpdateConfig()
}

// RequestEnterpriseCertAndUpdateConfig 申请企业证书并更新配置，返回是否已更换证书
func RequestEnterpriseCertAndUpdateConfig() (renewed bool, err error) {
	cfg := g.GetEnterpriseConfig()

	var (
		keyBytes []byte
		crtBytes []byte
		crt      *x509.Certificate
	)

	// 从企业证书地址获取证书和私钥，未配置地址时向 SNCA 申请
	keyBytes, crtBytes, err = requestFromUrl(cfg.Url, cfg.Name, cfg.GetCertificate().SerialNumber.String())
	if errors.Is(err, edocseal.ErrEnterpriseCertificateNotFound) {
		keyBytes, crtBytes, err = requestFromSnca()
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

	// 校验私钥与证书是否匹配
	_, err = tls.X509KeyPair(ca.PEMEncoding(crtBytes, ca.BlocTypeCertificate), ca.PEMEncoding(keyBytes, ca.BlocTypePrivateKey))
	if err != nil {
		err = fmt.Errorf("企业证书与私钥不匹配: %w", err)
		return
	}

	zap.L().Info("申请企业证书成功",
		zap.String("name", cfg.Name),
		zap.String("serialNumber", crt.SerialNumber.String()),
		zap.Time("notBefore", crt.NotBefore),
		zap.Time("notAfter", crt.NotAfter),
	)

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

	// 改写配置文件并替换内存中的企业证书
	err = g.ReplaceEnterpriseCertificate(kf, cf)
	if err != nil {
		return
	}

	zap.L().Info("企业证书已更新", zap.String("certificate", cf), zap.String("privateKey", kf))

	renewed = true
	return
}
