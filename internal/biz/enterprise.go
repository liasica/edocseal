// Copyright (C) edocseal. 2024-present.
//
// Created at 2026-09-23, by liasica

package biz

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/ca"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/enterprise"
	"auroraride.com/edocseal/internal/ent/enterprisecertification"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/model"
	"auroraride.com/edocseal/third/snca"
)

const (
	enterpriseRenewBefore   = 7 * 24 * time.Hour // 企业证书剩余有效期不足该时长时重新签发
	enterpriseSelfSignYears = 10                 // 自签企业证书有效期（年）
)

// enterpriseIssueMutex 防止并发签约时重复签发企业证书
var enterpriseIssueMutex sync.Mutex

// QueryEnterprise 按统一社会信用代码查询企业，代码为空时返回默认企业
func QueryEnterprise(creditCode string) (e *ent.Enterprise, err error) {
	query := ent.NewDatabase().Enterprise.Query()
	if creditCode == "" {
		query.Where(enterprise.IsDefault(true))
	} else {
		query.Where(enterprise.CreditCode(creditCode))
	}

	e, err = query.First(context.Background())
	if ent.IsNotFound(err) {
		err = fmt.Errorf("未找到签约企业: %s", creditCode)
		if creditCode == "" {
			err = errors.New("未设置默认签约企业")
		}
	}
	return
}

// EnterpriseSealPath 返回企业签章路径，文件名为统一社会信用代码
func EnterpriseSealPath(creditCode string) (path string, err error) {
	path = filepath.Join(g.GetSealDir(), creditCode+".png")
	if !edocseal.FileExists(path) {
		err = fmt.Errorf("未找到企业签章: %s", path)
	}
	return
}

// EnterpriseCertificate 返回企业在指定生成方式下的证书，没有或即将过期时重新签发
func EnterpriseCertificate(e *ent.Enterprise, issuer string) (cert *ent.EnterpriseCertification, err error) {
	enterpriseIssueMutex.Lock()
	defer enterpriseIssueMutex.Unlock()

	cert, _ = ent.NewDatabase().EnterpriseCertification.Query().
		Where(
			enterprisecertification.CreditCode(e.CreditCode),
			enterprisecertification.Issuer(issuer),
		).
		First(context.Background())
	if cert != nil && cert.ExpiresAt.After(time.Now().Add(enterpriseRenewBefore)) {
		return
	}

	var issued *ent.EnterpriseCertification
	issued, err = issueEnterpriseCertificate(e, issuer)
	if err == nil {
		cert = issued
		return
	}

	// 签发失败时旧证书尚未过期则继续使用
	if cert != nil && cert.ExpiresAt.After(time.Now()) {
		zap.L().Error("企业证书续期失败，继续使用旧证书",
			zap.Error(err),
			zap.String("creditCode", e.CreditCode),
			zap.String("issuer", issuer),
			zap.Time("expiresAt", cert.ExpiresAt),
		)
		err = nil
		return
	}

	cert = nil
	return
}

// EnterpriseCertificatePEM 返回企业陕西CA 证书与私钥的 PEM 内容，供其他环境拉取
func EnterpriseCertificatePEM(creditCode string) (result *edocseal.EnterpriseCertificate, err error) {
	var e *ent.Enterprise
	e, err = QueryEnterprise(creditCode)
	if err != nil {
		return
	}

	var cert *ent.EnterpriseCertification
	cert, err = EnterpriseCertificate(e, model.CertificateIssuerSnca)
	if err != nil {
		return
	}

	var crt, key []byte
	crt, err = os.ReadFile(cert.CertPath)
	if err != nil {
		return
	}

	key, err = os.ReadFile(cert.PrivatePath)
	if err != nil {
		return
	}

	result = &edocseal.EnterpriseCertificate{Cert: string(crt), Key: string(key)}
	return
}

// SaveEnterprise 按统一社会信用代码新增或更新企业，设为默认时取消其他企业的默认标记
func SaveEnterprise(e *ent.Enterprise) error {
	ctx := context.Background()
	return ent.WithTx(ctx, func(tx *ent.Tx) (err error) {
		if e.IsDefault {
			err = tx.Enterprise.Update().
				Where(enterprise.CreditCodeNEQ(e.CreditCode)).
				SetIsDefault(false).
				Exec(ctx)
			if err != nil {
				return
			}
		}

		return tx.Enterprise.Create().
			SetCreditCode(e.CreditCode).
			SetName(e.Name).
			SetProvince(e.Province).
			SetCity(e.City).
			SetPersonName(e.PersonName).
			SetPhone(e.Phone).
			SetIdcard(e.Idcard).
			SetIsDefault(e.IsDefault).
			OnConflictColumns(enterprise.FieldCreditCode).
			Update(func(u *ent.EnterpriseUpsert) {
				u.UpdateName().UpdateProvince().UpdateCity().UpdatePersonName().UpdatePhone().UpdateIdcard()
				// 未要求设为默认时保留原有默认标记
				if e.IsDefault {
					u.UpdateIsDefault()
				}
			}).
			Exec(ctx)
	})
}

// ListEnterprises 返回全部企业及其证书记录
func ListEnterprises() (items []*ent.Enterprise, certs []*ent.EnterpriseCertification, err error) {
	ctx := context.Background()

	items, err = ent.NewDatabase().Enterprise.Query().Order(ent.Asc(enterprise.FieldID)).All(ctx)
	if err != nil {
		return
	}

	certs, err = ent.NewDatabase().EnterpriseCertification.Query().Order(ent.Asc(enterprisecertification.FieldID)).All(ctx)
	return
}

// RevokeEnterpriseCertificates 作废企业的全部证书记录，下次签约时重新签发
func RevokeEnterpriseCertificates(creditCode string) (int, error) {
	return ent.NewDatabase().EnterpriseCertification.Delete().
		Where(enterprisecertification.CreditCode(creditCode)).
		Exec(context.Background())
}

// 签发企业证书，保存证书文件并更新证书记录
func issueEnterpriseCertificate(e *ent.Enterprise, issuer string) (cert *ent.EnterpriseCertification, err error) {
	var crtBytes, keyBytes []byte
	switch issuer {
	case model.CertificateIssuerSelf:
		crtBytes, keyBytes, err = selfIssueEnterpriseCertificate(e)
	case model.CertificateIssuerSnca:
		if url := g.GetEnterpriseConfig().Url; url != "" {
			crtBytes, keyBytes, err = requestFromUrl(url, e.CreditCode)
		} else {
			crtBytes, keyBytes, err = requestFromSnca(e)
		}
	default:
		err = fmt.Errorf("不支持的证书生成方式: %s", issuer)
	}
	if err != nil {
		return
	}

	var crt *x509.Certificate
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

	dir := filepath.Join(g.GetCertificateDir(), "enterprise")
	err = edocseal.CreateDirectory(dir)
	if err != nil {
		return
	}

	base := filepath.Join(dir, fmt.Sprintf("%s_%s_%s", e.CreditCode, issuer, crt.SerialNumber))
	kp, cp := base+"_key.pem", base+"_cert.pem"

	err = ca.SaveToFile(kp, keyBytes, ca.BlocTypePrivateKey)
	if err != nil {
		return
	}

	err = ca.SaveToFile(cp, crtBytes, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	zap.L().Info("企业证书签发成功",
		zap.String("creditCode", e.CreditCode),
		zap.String("issuer", issuer),
		zap.String("serialNumber", crt.SerialNumber.String()),
		zap.Time("notAfter", crt.NotAfter),
	)

	return ent.NewDatabase().EnterpriseCertification.Create().
		SetCreditCode(e.CreditCode).
		SetIssuer(issuer).
		SetPrivatePath(kp).
		SetCertPath(cp).
		SetExpiresAt(crt.NotAfter).
		OnConflictColumns(enterprisecertification.FieldCreditCode, enterprisecertification.FieldIssuer).
		UpdateNewValues().
		Save(context.Background())
}

// 使用根证书签发企业证书
func selfIssueEnterpriseCertificate(e *ent.Enterprise) (crt, key []byte, err error) {
	root := g.LoadRootCertificate()
	if !root.IsValid() {
		err = errors.New("根证书不存在")
		return
	}

	crt, key, _, err = ca.CreateSigningCertificate(root.GetPrivateKey(), root.GetCertificate(), pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{e.Province},
		Locality:           []string{e.City},
		Organization:       []string{e.Name},
		OrganizationalUnit: []string{e.CreditCode},
		CommonName:         e.Name,
	}, enterpriseSelfSignYears)
	return
}

// 向陕西CA 申请企业证书
func requestFromSnca(e *ent.Enterprise) (crt, key []byte, err error) {
	if e.PersonName == "" || e.Phone == "" || e.Idcard == "" {
		err = fmt.Errorf("企业代办人信息不完整: %s", e.CreditCode)
		return
	}

	priKey := ca.GenerateRsaPrivateKey()
	key, _ = x509.MarshalPKCS8PrivateKey(priKey)

	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:    []string{"CN"},
		CommonName: e.Name,
	})
	if err != nil {
		return
	}

	crt, err = snca.New().RequestCACert(
		snca.CertTypeEnterprise,
		csr,
		e.Name,
		e.PersonName,
		e.Phone,
		e.Idcard,
		e.Province,
		e.City,
		e.CreditCode,
	)
	return
}

// 从企业证书地址拉取陕西CA 企业证书和私钥，返回 DER 编码内容
func requestFromUrl(url, creditCode string) (crt, key []byte, err error) {
	client := resty.New().SetTimeout(30 * time.Second)
	defer func() {
		_ = client.Close()
	}()

	var (
		result edocseal.EnterpriseCertificate
		resp   *resty.Response
	)
	resp, err = client.R().
		SetResult(&result).
		SetQueryParam("code", creditCode).
		Get(url)
	if err != nil {
		return
	}

	if !resp.IsSuccess() {
		err = fmt.Errorf("企业证书地址响应异常: %s", resp.Status())
		return
	}

	crt, err = decodePEM(result.Cert)
	if err != nil {
		err = fmt.Errorf("企业证书解析失败: %w", err)
		return
	}

	key, err = decodePEM(result.Key)
	if err != nil {
		err = fmt.Errorf("企业私钥解析失败: %w", err)
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
