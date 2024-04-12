// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/tjfoc/gmsm/sm3"
	"go.uber.org/zap"

	"github.com/liasica/edocseal/ca"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
)

// RequestCertificae 申请证书
func RequestCertificae(paths *model.DocumentPaths, name, province, city, address, phone, idcard string) (err error) {
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
	// 保存证书
	err = ca.SaveToFile(paths.Cert, crt, ca.BlocTypeCertificate)
	if err != nil {
		return
	}

	// 保存私钥
	return ca.SaveToFile(paths.Key, key, ca.BlocTypePrivateKey)
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
	param := make(map[string]string)
	param["alg"] = "RSA"
	param["certType"] = "1"
	param["province"] = province
	param["city"] = city
	param["agent"] = name
	param["agentTel"] = phone
	param["agentNumber"] = idcard
	param["clientName"] = name
	param["legalNumber"] = idcard
	param["legalPersonTel"] = phone
	param["companyAdd"] = address
	param["prjId"] = "401"
	param["pkcs10"] = pkcs10

	h := sm3.New()
	h.Write([]byte(param["agentNumber"] + "_" + param["clientName"]))
	sum := h.Sum(nil)
	param["passwdDigest"] = hex.EncodeToString(sum)

	var resp *resty.Response
	resp, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(param).
		SetResult(&model.AgencyCertResponse{}).
		Post("http://117.32.132.76:28998/boss-eai/cert/getCert")
	if err != nil {
		zap.L().Fatal("机构签发证书失败", zap.Error(err))
		return
	}
	res := new(model.AgencyCertResponse)
	err = jsoniter.Unmarshal(resp.Body(), res)
	if err != nil {
		zap.L().Fatal("机构签发证书失败", zap.Error(err))
		return
	}
	if res.ErrorCode != "200" || res.Result == "false" {
		zap.L().Fatal("机构签发证书失败", zap.Error(errors.New(res.ErrorMessage)))
		return nil, nil, errors.New(res.ErrorMessage)
	}

	crt, err = base64.StdEncoding.DecodeString(res.JsonObj.SignCert)
	return
}
