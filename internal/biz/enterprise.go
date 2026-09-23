// Copyright (C) edocseal. 2024-present.
//
// Created at 2026-09-23, by liasica

package biz

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/ca"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/certification"
	"auroraride.com/edocseal/internal/ent/enterprise"
	"auroraride.com/edocseal/internal/ent/enterprisecertification"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/internal/model"
	"auroraride.com/edocseal/pb"
	"auroraride.com/edocseal/third/snca"
)

const (
	enterpriseRenewBefore     = 7 * 24 * time.Hour  // 企业证书剩余有效期不足该时长时重新签发
	enterpriseRootRenewBefore = 30 * 24 * time.Hour // 企业根证书剩余有效期不足该时长时重新生成
	enterpriseSelfSignYears   = 10                  // 自签企业证书有效期（年）
)

// 证书状态
const (
	certificateStatusValid    = "VALID"    // 有效
	certificateStatusRenewing = "RENEWING" // 待续签
	certificateStatusExpired  = "EXPIRED"  // 已过期
)

// enterpriseCache 签约企业与企业证书的内存缓存，启动时从数据库加载，变更时先写数据库再更新缓存
type enterpriseCache struct {
	sync.RWMutex

	items       map[string]*ent.Enterprise              // 统一社会信用代码 -> 企业
	certs       map[string]*ent.EnterpriseCertification // 统一社会信用代码与生成方式 -> 企业证书
	defaultCode string                                  // 当前签约企业
}

var (
	enterprises = &enterpriseCache{}

	// enterpriseWriteMutex 串行化企业、企业证书与根证书的写操作
	enterpriseWriteMutex sync.Mutex
)

func enterpriseCertKey(creditCode, issuer string) string {
	return creditCode + "|" + issuer
}

// LoadEnterprises 从数据库加载签约企业与企业证书到内存
func LoadEnterprises() (err error) {
	ctx := context.Background()

	var items []*ent.Enterprise
	items, err = ent.NewDatabase().Enterprise.Query().All(ctx)
	if err != nil {
		return
	}

	var certs []*ent.EnterpriseCertification
	certs, err = ent.NewDatabase().EnterpriseCertification.Query().All(ctx)
	if err != nil {
		return
	}

	itemMap := make(map[string]*ent.Enterprise, len(items))
	var defaultCode string
	for _, item := range items {
		itemMap[item.CreditCode] = item
		if item.IsDefault {
			defaultCode = item.CreditCode
		}
	}

	certMap := make(map[string]*ent.EnterpriseCertification, len(certs))
	for _, cert := range certs {
		certMap[enterpriseCertKey(cert.CreditCode, cert.Issuer)] = cert
	}

	enterprises.Lock()
	enterprises.items = itemMap
	enterprises.certs = certMap
	enterprises.defaultCode = defaultCode
	enterprises.Unlock()
	return
}

// QueryEnterprise 从缓存查询企业，统一社会信用代码为空时返回当前签约企业
func QueryEnterprise(creditCode string) (e *ent.Enterprise, err error) {
	enterprises.RLock()
	defer enterprises.RUnlock()

	if creditCode == "" {
		creditCode = enterprises.defaultCode
		if creditCode == "" {
			err = errors.New("未设置签约企业")
			return
		}
	}

	e = enterprises.items[creditCode]
	if e == nil {
		err = fmt.Errorf("未找到签约企业: %s", creditCode)
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
	cert = cachedEnterpriseCertificate(e.CreditCode, issuer)
	if enterpriseCertificateUsable(cert) {
		return
	}

	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	// 等待期间其他请求可能已完成签发
	cert = cachedEnterpriseCertificate(e.CreditCode, issuer)
	if enterpriseCertificateUsable(cert) {
		return
	}

	var issued *ent.EnterpriseCertification
	issued, err = issueEnterpriseCertificate(e, issuer)
	if err == nil {
		enterprises.Lock()
		enterprises.certs[enterpriseCertKey(e.CreditCode, issuer)] = issued
		enterprises.Unlock()

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

// ListEnterprises 返回全部企业及其签章与证书，当前签约企业排在最前
func ListEnterprises() (items []*pb.Enterprise) {
	enterprises.RLock()
	defer enterprises.RUnlock()

	for _, e := range enterprises.items {
		item := &pb.Enterprise{
			CreditCode: e.CreditCode,
			Name:       e.Name,
			Province:   e.Province,
			City:       e.City,
			PersonName: e.PersonName,
			Phone:      e.Phone,
			Idcard:     e.Idcard,
			IsDefault:  e.IsDefault,
		}

		if path, err := EnterpriseSealPath(e.CreditCode); err == nil {
			item.Seal, _ = os.ReadFile(path)
		}

		if e.RootCertPath != "" {
			item.RootCertificate = certificateDetail(e.RootCertPath, enterpriseRootRenewBefore)
		}

		for _, issuer := range []string{model.CertificateIssuerSelf, model.CertificateIssuerSnca} {
			cert := enterprises.certs[enterpriseCertKey(e.CreditCode, issuer)]
			if cert == nil {
				continue
			}

			item.Certificates = append(item.Certificates, &pb.EnterpriseCertificate{
				Issuer: issuer,
				Detail: certificateDetail(cert.CertPath, enterpriseRenewBefore),
			})
		}

		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDefault != items[j].IsDefault {
			return items[i].IsDefault
		}
		return items[i].CreditCode < items[j].CreditCode
	})
	return
}

// SaveEnterprise 按统一社会信用代码新增或更新企业，签章为空时保留原签章；第一个企业自动设为签约企业
func SaveEnterprise(req *pb.EnterpriseSaveRequest) (err error) {
	if len(req.CreditCode) != 18 {
		return errors.New("统一社会信用代码应为 18 位")
	}
	if req.Name == "" || req.Province == "" || req.City == "" {
		return errors.New("企业名称、省份、城市不能为空")
	}

	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	old, _ := QueryEnterprise(req.CreditCode)
	if old == nil && len(req.Seal) == 0 {
		return errors.New("请上传签章图片")
	}

	if len(req.Seal) > 0 {
		err = saveEnterpriseSeal(req.CreditCode, req.Seal)
		if err != nil {
			return
		}
	}

	enterprises.RLock()
	first := enterprises.defaultCode == ""
	enterprises.RUnlock()

	ctx := context.Background()
	err = ent.NewDatabase().Enterprise.Create().
		SetCreditCode(req.CreditCode).
		SetName(req.Name).
		SetProvince(req.Province).
		SetCity(req.City).
		SetPersonName(req.PersonName).
		SetPhone(req.Phone).
		SetIdcard(req.Idcard).
		SetIsDefault(first).
		OnConflictColumns(enterprise.FieldCreditCode).
		Update(func(u *ent.EnterpriseUpsert) {
			u.UpdateName().UpdateProvince().UpdateCity().UpdatePersonName().UpdatePhone().UpdateIdcard()
		}).
		Exec(ctx)
	if err != nil {
		return
	}

	return LoadEnterprises()
}

// DeleteEnterprise 删除企业及其证书记录与签章，当前签约企业不能删除
func DeleteEnterprise(creditCode string) (err error) {
	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	var e *ent.Enterprise
	e, err = queryEnterpriseByCode(creditCode)
	if err != nil {
		return
	}

	if e.IsDefault {
		return errors.New("当前签约企业不能删除")
	}

	ctx := context.Background()
	err = ent.WithTx(ctx, func(tx *ent.Tx) (err error) {
		_, err = tx.EnterpriseCertification.Delete().Where(enterprisecertification.CreditCode(creditCode)).Exec(ctx)
		if err != nil {
			return
		}

		_, err = tx.Enterprise.Delete().Where(enterprise.CreditCode(creditCode)).Exec(ctx)
		return
	})
	if err != nil {
		return
	}

	if path, pathErr := EnterpriseSealPath(creditCode); pathErr == nil {
		_ = os.Remove(path)
	}

	return LoadEnterprises()
}

// SetDefaultEnterprise 设为签约企业，签约时立即生效
func SetDefaultEnterprise(creditCode string) (err error) {
	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	_, err = queryEnterpriseByCode(creditCode)
	if err != nil {
		return
	}

	_, err = EnterpriseSealPath(creditCode)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = ent.WithTx(ctx, func(tx *ent.Tx) (err error) {
		err = tx.Enterprise.Update().Where(enterprise.IsDefault(true)).SetIsDefault(false).Exec(ctx)
		if err != nil {
			return
		}

		return tx.Enterprise.Update().Where(enterprise.CreditCode(creditCode)).SetIsDefault(true).Exec(ctx)
	})
	if err != nil {
		return
	}

	return LoadEnterprises()
}

// 企业证书生成结果
const (
	generateStatusGenerated = "GENERATED" // 已生成
	generateStatusExists    = "EXISTS"    // 证书未过期无需生成
	generateStatusFailed    = "FAILED"    // 生成失败
)

// GenerateCertificates 按证书生成方式为缺少证书或证书已过期的企业生成企业证书，未过期的证书跳过
func GenerateCertificates(issuer string) (results []*pb.EnterpriseGenerateResult, err error) {
	if issuer != model.CertificateIssuerSelf && issuer != model.CertificateIssuerSnca {
		err = fmt.Errorf("不支持的证书生成方式: %s", issuer)
		return
	}

	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	enterprises.RLock()
	items := make([]*ent.Enterprise, 0, len(enterprises.items))
	for _, e := range enterprises.items {
		items = append(items, e)
	}
	enterprises.RUnlock()

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreditCode < items[j].CreditCode
	})

	for _, e := range items {
		result := &pb.EnterpriseGenerateResult{CreditCode: e.CreditCode, Name: e.Name}
		results = append(results, result)

		if cert := cachedEnterpriseCertificate(e.CreditCode, issuer); cert != nil && cert.ExpiresAt.After(time.Now()) {
			result.Status = generateStatusExists
			continue
		}

		issued, issueErr := issueEnterpriseCertificate(e, issuer)
		if issueErr != nil {
			zap.L().Error("企业证书生成失败", zap.Error(issueErr), zap.String("creditCode", e.CreditCode), zap.String("issuer", issuer))
			result.Status = generateStatusFailed
			result.Message = issueErr.Error()
			continue
		}

		enterprises.Lock()
		enterprises.certs[enterpriseCertKey(e.CreditCode, issuer)] = issued
		enterprises.Unlock()

		result.Status = generateStatusGenerated
	}
	return
}

// RegenerateRootCertificate 重新生成企业自签根证书，该企业的自签证书随之作废；根证书未过期时不能重新生成
func RegenerateRootCertificate(creditCode string) (err error) {
	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	var e *ent.Enterprise
	e, err = queryEnterpriseByCode(creditCode)
	if err != nil {
		return
	}

	if e.RootExpiresAt != nil && e.RootExpiresAt.After(time.Now()) {
		return errors.New("根证书未过期，不能重新生成")
	}

	_, _, err = generateEnterpriseRootLocked(e)
	return
}

// EnterpriseRoot 返回企业自签根证书，没有或即将过期时生成
func EnterpriseRoot(e *ent.Enterprise) (crt *x509.Certificate, key *rsa.PrivateKey, err error) {
	crt, key, err = loadEnterpriseRoot(e)
	if err == nil {
		return
	}

	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	return ensureEnterpriseRootLocked(e)
}

// RenewCertificates 续签即将到期的企业根证书与企业证书，供定时任务调用
func RenewCertificates() {
	enterpriseWriteMutex.Lock()
	defer enterpriseWriteMutex.Unlock()

	enterprises.RLock()
	items := make([]*ent.Enterprise, 0, len(enterprises.items))
	for _, e := range enterprises.items {
		items = append(items, e)
	}
	certs := make([]*ent.EnterpriseCertification, 0, len(enterprises.certs))
	for _, cert := range enterprises.certs {
		certs = append(certs, cert)
	}
	enterprises.RUnlock()

	for _, e := range items {
		if e.RootExpiresAt == nil || e.RootExpiresAt.After(time.Now().Add(enterpriseRootRenewBefore)) {
			continue
		}

		if _, _, err := generateEnterpriseRootLocked(e); err != nil {
			zap.L().Error("企业根证书续签失败", zap.Error(err), zap.String("creditCode", e.CreditCode))
		}
	}

	for _, cert := range certs {
		if enterpriseCertificateUsable(cert) {
			continue
		}

		e, err := QueryEnterprise(cert.CreditCode)
		if err != nil {
			continue
		}

		var issued *ent.EnterpriseCertification
		issued, err = issueEnterpriseCertificate(e, cert.Issuer)
		if err != nil {
			zap.L().Error("企业证书续签失败", zap.Error(err), zap.String("creditCode", cert.CreditCode), zap.String("issuer", cert.Issuer))
			continue
		}

		enterprises.Lock()
		enterprises.certs[enterpriseCertKey(cert.CreditCode, cert.Issuer)] = issued
		enterprises.Unlock()
	}

	zap.L().Info("企业证书自动续签检查完成", zap.Int("enterprises", len(items)), zap.Int("certificates", len(certs)))
}

// 读取企业根证书，不存在或即将过期时返回错误
func loadEnterpriseRoot(e *ent.Enterprise) (crt *x509.Certificate, key *rsa.PrivateKey, err error) {
	if e.RootCertPath == "" || e.RootExpiresAt == nil || !e.RootExpiresAt.After(time.Now().Add(enterpriseRootRenewBefore)) {
		err = fmt.Errorf("企业根证书不可用: %s", e.CreditCode)
		return
	}

	crt, err = ca.LoadCertificateFromFile(e.RootCertPath)
	if err != nil {
		return
	}

	key, err = ca.LoadPrivateKeyFromFile(e.RootPrivatePath)
	return
}

// 返回可用的企业根证书，没有时生成，调用方需持有 enterpriseWriteMutex
func ensureEnterpriseRootLocked(e *ent.Enterprise) (crt *x509.Certificate, key *rsa.PrivateKey, err error) {
	// 以缓存中的最新记录为准，等待锁期间可能已生成
	if latest, queryErr := QueryEnterprise(e.CreditCode); queryErr == nil {
		e = latest
	}

	crt, key, err = loadEnterpriseRoot(e)
	if err == nil {
		return
	}

	return generateEnterpriseRootLocked(e)
}

// 生成企业根证书并写入企业记录，作废该企业的自签企业证书；为签约企业时一并作废自签个人证书。调用方需持有 enterpriseWriteMutex
func generateEnterpriseRootLocked(e *ent.Enterprise) (crt *x509.Certificate, key *rsa.PrivateKey, err error) {
	key = ca.GenerateRsaPrivateKey()

	var der []byte
	der, err = ca.GenerateRootCertificate(key, pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{e.Province},
		Locality:           []string{e.City},
		Organization:       []string{e.Name},
		OrganizationalUnit: []string{e.CreditCode},
		CommonName:         e.Name + " Root CA",
	})
	if err != nil {
		return
	}

	crt, err = x509.ParseCertificate(der)
	if err != nil {
		return
	}

	dir := filepath.Join(g.GetCertificateDir(), "root")
	err = edocseal.CreateDirectory(dir)
	if err != nil {
		return
	}

	keyBytes, _ := x509.MarshalPKCS8PrivateKey(key)
	base := filepath.Join(dir, fmt.Sprintf("%s_%s", e.CreditCode, crt.SerialNumber))
	kp, cp := base+"_key.pem", base+"_cert.pem"

	err = ca.SaveToFile(kp, keyBytes, ca.BlocTypePrivateKey)
	if err != nil {
		return
	}

	err = ca.SaveToFile(cp, der, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = ent.NewDatabase().Enterprise.Update().
		Where(enterprise.CreditCode(e.CreditCode)).
		SetRootCertPath(cp).
		SetRootPrivatePath(kp).
		SetRootExpiresAt(crt.NotAfter).
		Exec(ctx)
	if err != nil {
		return
	}

	_, err = ent.NewDatabase().EnterpriseCertification.Delete().
		Where(
			enterprisecertification.CreditCode(e.CreditCode),
			enterprisecertification.Issuer(model.CertificateIssuerSelf),
		).
		Exec(ctx)
	if err != nil {
		return
	}

	if e.IsDefault {
		_, err = ent.NewDatabase().Certification.Delete().
			Where(certification.Issuer(model.CertificateIssuerSelf)).
			Exec(ctx)
		if err != nil {
			return
		}
	}

	zap.L().Info("企业根证书已生成",
		zap.String("creditCode", e.CreditCode),
		zap.String("serialNumber", crt.SerialNumber.String()),
		zap.Time("notAfter", crt.NotAfter),
	)

	err = LoadEnterprises()
	return
}

// 解析证书文件为证书详情，文件无法解析时返回 nil
func certificateDetail(path string, renewBefore time.Duration) *pb.CertificateDetail {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var crt *x509.Certificate
	crt, err = ca.ParseCertificate(b)
	if err != nil {
		return nil
	}

	publicKey := crt.PublicKeyAlgorithm.String()
	if rsaKey, ok := crt.PublicKey.(*rsa.PublicKey); ok {
		publicKey = fmt.Sprintf("RSA %d 位", rsaKey.N.BitLen())
	}

	sum := sha256.Sum256(crt.Raw)
	fingerprint := make([]string, len(sum))
	for i, v := range sum {
		fingerprint[i] = fmt.Sprintf("%02X", v)
	}

	status := certificateStatusValid
	switch now := time.Now(); {
	case !crt.NotAfter.After(now):
		status = certificateStatusExpired
	case crt.NotAfter.Before(now.Add(renewBefore)):
		status = certificateStatusRenewing
	}

	return &pb.CertificateDetail{
		Subject:            crt.Subject.String(),
		Issuer:             crt.Issuer.String(),
		Serial:             fmt.Sprintf("%X", crt.SerialNumber),
		NotBefore:          crt.NotBefore.Unix(),
		NotAfter:           crt.NotAfter.Unix(),
		Status:             status,
		PublicKey:          publicKey,
		SignatureAlgorithm: crt.SignatureAlgorithm.String(),
		KeyUsage:           keyUsageNames(crt.KeyUsage),
		IsCa:               crt.IsCA,
		Fingerprint:        strings.Join(fingerprint, ":"),
		Pem:                string(b),
	}
}

// 密钥用法名称
func keyUsageNames(usage x509.KeyUsage) (names []string) {
	for _, item := range []struct {
		usage x509.KeyUsage
		name  string
	}{
		{x509.KeyUsageDigitalSignature, "数字签名"},
		{x509.KeyUsageContentCommitment, "不可否认"},
		{x509.KeyUsageKeyEncipherment, "密钥加密"},
		{x509.KeyUsageDataEncipherment, "数据加密"},
		{x509.KeyUsageKeyAgreement, "密钥协商"},
		{x509.KeyUsageCertSign, "证书签发"},
		{x509.KeyUsageCRLSign, "CRL 签发"},
	} {
		if usage&item.usage != 0 {
			names = append(names, item.name)
		}
	}
	return
}

// 按非空的统一社会信用代码从缓存查询企业
func queryEnterpriseByCode(creditCode string) (*ent.Enterprise, error) {
	if creditCode == "" {
		return nil, errors.New("统一社会信用代码不能为空")
	}
	return QueryEnterprise(creditCode)
}

func cachedEnterpriseCertificate(creditCode, issuer string) *ent.EnterpriseCertification {
	enterprises.RLock()
	defer enterprises.RUnlock()

	return enterprises.certs[enterpriseCertKey(creditCode, issuer)]
}

func enterpriseCertificateUsable(cert *ent.EnterpriseCertification) bool {
	return cert != nil && cert.ExpiresAt.After(time.Now().Add(enterpriseRenewBefore))
}

// 校验签章为 PNG 图片后写入签章目录
func saveEnterpriseSeal(creditCode string, seal []byte) (err error) {
	var format string
	_, format, err = image.DecodeConfig(bytes.NewReader(seal))
	if err != nil || format != "png" {
		return errors.New("签章图片须为 PNG 格式")
	}

	path := filepath.Join(g.GetSealDir(), creditCode+".png")
	err = os.WriteFile(path+".tmp", seal, 0o644)
	if err != nil {
		return
	}

	return os.Rename(path+".tmp", path)
}

// 签发企业证书，保存证书文件并写入证书记录
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

	// 签发或拉取到的证书须未过期且属于该企业
	if !crt.NotAfter.After(time.Now()) {
		err = fmt.Errorf("获取到的证书已过期，到期时间 %s", crt.NotAfter.Local().Format(time.DateTime))
		return
	}

	if crt.Subject.CommonName != e.Name {
		err = fmt.Errorf("获取到的证书不属于该企业，证书主题为 %s", crt.Subject.CommonName)
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

// 使用企业根证书签发企业证书，调用方需持有 enterpriseWriteMutex
func selfIssueEnterpriseCertificate(e *ent.Enterprise) (crt, key []byte, err error) {
	var (
		rootCrt *x509.Certificate
		rootKey *rsa.PrivateKey
	)

	rootCrt, rootKey, err = ensureEnterpriseRootLocked(e)
	if err != nil {
		return
	}

	crt, key, _, err = ca.CreateSigningCertificate(rootKey, rootCrt, pkix.Name{
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
