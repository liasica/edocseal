// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

type ApplySealCertRequest struct {
	TokenInfo  string `json:"tokenInfo"`  // 客户端随机数（RandomA16位）与服务端 随机数（RandomB）字符串拼接，randomA+randomB+"SNCA"+appId（归档返回）
	CommonName string `json:"commonName"` // 证书主体的名称项，不大于 200 个字节
	P10        string `json:"p10"`        // P10
}

type ApplySealCertResponse struct {
	TaskCode      string            `json:"taskCode,omitempty"`      // 业务类型 applySealCert
	ResultCode    string            `json:"resultCode,omitempty"`    // 返回值，0 表示成功，其他表示错误
	ResultCodeMsg string            `json:"resultCodeMsg,omitempty"` // 返回信息描述
	Data          ApplySealCertData `json:"data,omitempty"`
}

type ApplySealCertData struct {
	SignCert         string `json:"signCert"`         // Base64 编码签名证书
	EncCert          string `json:"encCert"`          // Base64 编码的加密证书
	EncKeyProtection string `json:"encKeyProtection"` // Base64 编码的加密密钥对保护结构，被编码的结构 ENVELOPEDKEYBLOB，遵从GM/T 0016—2012 标准
	TransactionID    string `json:"transactionID"`    // transactionID
}

type ApplyServiceRandomResponse struct {
	ResultCode    string `json:"resultCode,omitempty"`    // 返回值，0 表示成功，其他表示错误
	ResultCodeMsg string `json:"resultCodeMsg,omitempty"` // 返回信息描述
	RandomB       string `json:"randomB,omitempty"`       // 服务端随机数 16 个字节
}

type BusinessDataFinishRequest struct {
	AppType          string `json:"appType"`          // 操作类型 新办：1，延期
	CertType         string `json:"certType"`         // 证书类型 企业：2 ,个人：9
	SubCert          string `json:"subCert"`          // 主子证标识 目前均为主证 主证：1 子证：0
	Source           string `json:"source"`           // 第三方来源标识（可在联调时约定）
	ClientName       string `json:"clientName"`       // 证书名称
	CountryName      string `json:"countryName"`      // 国家（证书C项）
	Agent            string `json:"agent"`            // 经办人
	AgentTel         string `json:"agentTel"`         // 经办人电话
	AgentNumber      string `json:"agentNumber"`      // 经办人证件号码
	CertId           string `json:"certId"`           // 证书序列号 延期、注销、变更时  除新办,丢补外均要传入
	CustomerType     string `json:"customerType"`     // 行业类型 联调对接时可分配
	SocialCreditCode string `json:"socialCreditCode"` // 社会信用代码（单位/业务证书必须传入）
	Province         string `json:"province"`         // 省份（证书L项）非必填
	City             string `json:"city"`             // 城市 （证书S项）非必填
}

type BusinessDataFinishResponse struct {
	Result    string                 `json:"result,omitempty"`    // 业务数据接口操作结果
	ResultMsg string                 `json:"resultMsg,omitempty"` // 业务数据接口操作结果描述
	Data      BusinessDataFinishData `json:"data,omitempty"`
}

type BusinessDataFinishData struct {
	TrustId    string `json:"trustId,omitempty"`    // 客服信任号
	AppId      string `json:"appId,omitempty"`      // 业务单号
	ClientName string `json:"clientName,omitempty"` // 证书名称
}
