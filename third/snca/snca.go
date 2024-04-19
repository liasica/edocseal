// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-17, by liasica

package snca

import (
	"encoding/base64"
	"errors"
	"math/big"
	"math/rand"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/liasica/edocseal"
)

type Snca struct {
	url          string
	source       string
	customerType string
}

func NewSnca(url, source, customerType string) *Snca {
	return &Snca{
		url:          strings.TrimSuffix(url, "/"),
		source:       source,
		customerType: customerType,
	}
}

// ApplyServiceRan 获取服务端随机数
// 流程第一步
func (s *Snca) ApplyServiceRan() (randomB string, err error) {
	var (
		res    *resty.Response
		result ApplyServiceRandomResponse
	)

	res, err = resty.New().R().
		SetBody(&ApplyServiceRandomRequest{
			Source: s.source,
		}).
		SetResult(&result).
		Post(s.NewUrl(UrlApplyServiceRandom))

	if err != nil {
		return
	}

	zap.L().Info("获取服务端随机数返回结果", zap.String("response", string(res.Body())))

	if result.ResultCode != "0" {
		err = errors.New(result.ResultCodeMsg)
		return
	}

	randomB = result.RandomB
	return
}

// BusinessDataFinish 组装业务数据
// 流程第二步
func (s *Snca) BusinessDataFinish(typ CertType, name, personName, phone, idcard, province, city, creditCode string) (appId string, err error) {
	serial := big.NewInt(rand.Int63()).String()

	req := BusinessDataFinishRequest{
		AppType:          OpTypeNewly, // 暂无延期，固定为新办
		CertType:         typ,
		SubCert:          SubCertMaster, // 固定为主证
		Source:           s.source,
		PaymentMode:      PaymentModeCash,
		ClientName:       name,
		Agent:            personName,
		AgentTel:         phone,
		AgentNumber:      idcard,
		CertId:           serial,
		CustomerType:     s.customerType,
		SocialCreditCode: creditCode,
		CountryName:      "CN",
		Province:         province,
		City:             city,
	}

	if typ == CertTypeEnterprise {
		req.OrganizationName = "西安时光驹新能源科技有限公司"
		req.Department = "AUR"
	}

	if typ == CertTypePersonal {
		req.LegalPerson = personName
		req.LegalNumber = idcard
		req.LegalPersonTel = phone
	}

	var (
		res    *resty.Response
		result BusinessDataFinishResponse
	)
	res, err = resty.New().R().
		SetBody(req).
		SetResult(&result).
		Post(s.NewUrl(UrlBusinessDataFinish))
	if err != nil {
		return
	}
	zap.L().Info("组装数据请求结果", zap.String("response", string(res.Body())))

	if result.Result != "true" {
		err = errors.New(result.ResultMsg)
		return
	}

	appId = result.Data.AppId
	return
}

// ApplySealCert 申请证书
// 流程第三步
func (s *Snca) ApplySealCert(randomB, appId, name, csr string) (b []byte, err error) {
	var (
		res    *resty.Response
		result ApplySealCertResponse
	)

	res, err = resty.New().R().
		SetBody(ApplySealCertRequest{
			TokenInfo:  edocseal.RandStr(16) + randomB + "SNCA" + appId,
			CommonName: name,
			P10:        csr,
		}).
		SetResult(&result).
		Post(s.NewUrl(UrlApplySealCert))
	if err != nil {
		return
	}

	zap.L().Info("机构签发证书返回", zap.String("response", string(res.Body())))

	if result.ResultCode != "0" {
		err = errors.New(result.ResultCodeMsg)
		return
	}

	return base64.StdEncoding.DecodeString(result.Data.SignCert)
}

// RequestCACert 请求CA证书
func (s *Snca) RequestCACert(typ CertType, csrBytes []byte, name, personName, phone, idcard, province, city, creditCode string) (b []byte, err error) {
	var (
		randomB string
		appId   string
	)

	randomB, err = s.ApplyServiceRan()
	if err != nil {
		return
	}

	appId, err = s.BusinessDataFinish(typ, name, personName, phone, idcard, province, city, creditCode)
	if err != nil {
		return
	}

	csr := base64.StdEncoding.EncodeToString(csrBytes)
	return s.ApplySealCert(randomB, appId, name, csr)
}
