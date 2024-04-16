// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/ent"
	"github.com/liasica/edocseal/internal/ent/certification"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
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

// 机构签发证书
// func agencyIssueCertificate(name, province, city, address, phone, idcard string) (crt, key []byte, err error) {
// 	// 生成私钥
// 	priKey := ca.GenerateRsaPrivateKey()
// 	key, _ = x509.MarshalPKCS8PrivateKey(key)
// 	// 生成CSR请求
// 	var csr []byte
// 	csr, err = ca.GenerateRequest(priKey, pkix.Name{
// 		Country:      []string{"CN"},
// 		Organization: []string{name},
// 	})
//
// 	pkcs10 := base64.StdEncoding.EncodeToString(csr)
// 	certParam := make(map[string]string)
// 	certParam["alg"] = "RSA"
// 	certParam["certType"] = "1"
// 	certParam["province"] = province
// 	certParam["city"] = city
// 	certParam["agent"] = name
// 	certParam["agentTel"] = phone
// 	certParam["agentNumber"] = idcard
// 	certParam["clientName"] = name
// 	certParam["legalNumber"] = idcard
// 	certParam["legalPersonTel"] = phone
// 	certParam["companyAdd"] = address
// 	certParam["prjId"] = "401"
// 	certParam["pkcs10"] = pkcs10
//
// 	h := sm3.New()
// 	h.Write([]byte(certParam["agentNumber"] + "_" + certParam["clientName"]))
// 	sum := h.Sum(nil)
// 	certParam["passwdDigest"] = hex.EncodeToString(sum)
//
// 	var (
// 		resp   *resty.Response
// 		result *model.AgencyCertResponse
// 	)
//
// 	resp, err = resty.New().R().
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(certParam).
// 		SetResult(&result).
// 		Post(g.GetSnca().Url)
// 	if err != nil {
// 		zap.L().Error("机构签发证书失败", zap.Error(err))
// 		return
// 	}
// 	zap.L().Info("机构签发证书", zap.String("response", string(resp.Body())))
//
// 	if result == nil {
// 		return nil, nil, errors.New("机构签发证书失败，返回结果为空")
// 	}
//
// 	if result.ErrorCode != "200" || result.Result == "false" {
// 		zap.L().Error("机构签发证书失败", zap.Error(errors.New(result.ErrorMessage)))
// 		return nil, nil, errors.New(result.ErrorMessage)
// 	}
//
// 	crt, err = base64.StdEncoding.DecodeString(result.JsonObj.SignCert)
// 	return
// }

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
	key, _ = x509.MarshalPKCS8PrivateKey(key)
	// 生成CSR请求
	var csr []byte
	csr, err = ca.GenerateRequest(priKey, pkix.Name{
		Country:      []string{"CN"},
		Organization: []string{name},
	})
	pkcs10 := base64.StdEncoding.EncodeToString(csr)
	log.Println(pkcs10)

	// 1.陕西CA随机数接口
	randomParam := make(map[string]string)
	randomParam["source"] = "SCNX"
	var (
		randomResp   *resty.Response
		randomResult *model.RandomBResponse
	)
	randomResp, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(randomParam).
		SetResult(&randomResult).
		Post("http://111.20.164.183:8998/sealCertService/com/api/cert/applyServiceRandom")
	if err != nil {
		zap.L().Error("机构随机数请求失败", zap.Error(err))
		return
	}
	zap.L().Info("机构随机数请求", zap.String("response", string(randomResp.Body())))
	if randomResult == nil {
		return nil, nil, errors.New("机构随机数请求失败，返回结果为空")
	}
	if randomResult.ResultCode != "0" {
		zap.L().Error("机构随机数请求失败", zap.Error(errors.New(randomResult.ResultCodeMsg)))
		return nil, nil, errors.New(randomResult.ResultCodeMsg)
	}

	// 2.陕西CA组装业务数据接口
	busParam := make(map[string]string)
	busParam["appType"] = "1"
	busParam["certType"] = "9"
	busParam["subCert"] = "1"
	busParam["source"] = "SCNX"
	busParam["clientName"] = "西安时光驹新能源科技有限公司"
	busParam["countryName"] = "CN"
	busParam["agent"] = name
	busParam["agentTel"] = phone
	busParam["agentNumber"] = idcard

	busParam["certId"] = strings.ReplaceAll(uuid.New().String(), "-", "")
	busParam["customerType"] = "1"

	busParam["socialCreditCode"] = "91610133MA6U8RAJ1X"

	busParam["province"] = province
	busParam["city"] = city

	var (
		busResp   *resty.Response
		busResult *model.BusResponse
	)
	busResp, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(busParam).
		SetResult(&busResult).
		Post("http://111.20.164.183:8998/sealCertService/com/api/cert/businessDataFinish")
	if err != nil {
		zap.L().Error("组装数据请求失败", zap.Error(err))
		return
	}
	zap.L().Info("组装数据请求失败", zap.String("response", string(busResp.Body())))
	if busResult == nil {
		return nil, nil, errors.New("组装数据请求失败，返回结果为空")
	}
	if busResult.Result != "true" {
		zap.L().Error("组装数据请求失败", zap.Error(errors.New(busResult.ResultMsg)))
		return nil, nil, errors.New(busResult.ResultMsg)
	}

	// 3.请求申请证书

	certParam := make(map[string]string)
	certParam["tokenInfo"] = fmt.Sprintf("%d", time.Now().UnixMicro()) + randomResult.RandomB + "SNCA" + busResult.Data.AppId
	certParam["commonName"] = "证书测试" // 西安时光驹新能源科技有限公司证书
	certParam["p10"] = pkcs10

	// certParam["tokenInfo"] = "16685798921901441668579892190144SNCA1288087"
	// certParam["commonName "] = "证书测试"
	// certParam["p10"] = "MIIBDjCBsgIBADB**/D1H4vuA9J5g=="

	var (
		certResp   *resty.Response
		certResult *model.AgencyCertResponse
	)
	marshal, err := json.Marshal(certParam)
	if err != nil {
		return nil, nil, err
	}

	log.Println(string(marshal))
	certResp, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(certParam).
		SetResult(&certResult).
		Post("http://111.20.164.183:8998/sealCertService/com/api/cert/applySealCert")
	if err != nil {
		zap.L().Error("机构签发证书失败", zap.Error(err))
		return
	}
	zap.L().Info("机构签发证书", zap.String("response", string(certResp.Body())))

	if certResult == nil {
		return nil, nil, errors.New("机构签发证书失败，返回结果为空")
	}

	if certResult.ResultCode != "0" {
		zap.L().Error("机构签发证书失败", zap.Error(errors.New(certResult.ResultCodeMsg)))
		return nil, nil, errors.New(certResult.ResultCodeMsg)
	}

	crt, err = base64.StdEncoding.DecodeString(certResult.Data.SignCert)
	return
}
