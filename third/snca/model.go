// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-17, by liasica

package snca

// OpType 操作类型
type OpType string

const (
	OpTypeNewly OpType = "1" // 新办
)

// CertType 证书类型
type CertType string

const (
	CertTypeEnterprise CertType = "2" // 企业
	CertTypePersonal   CertType = "9" // 个人
)

// SubCert 主子证标识
type SubCert string

const (
	SubCertMaster SubCert = "1" // 主证
	SubCertSlave  SubCert = "0" // 子证
)

type PaymentMode string

const (
	PaymentModeCash   PaymentMode = "1" // 现金
	PaymentModeWeChat PaymentMode = "6" // 微信
	PaymentModeAlipay PaymentMode = "7" // 支付宝
)

type ApplyServiceRandomRequest struct {
	Source string `json:"source,omitempty"` // 约定标记
}

type ApplyServiceRandomResponse struct {
	ResultCode    string `json:"resultCode,omitempty"`    // 返回值，0 表示成功，其他表示错误
	ResultCodeMsg string `json:"resultCodeMsg,omitempty"` // 返回信息描述
	RandomB       string `json:"randomB,omitempty"`       // 服务端随机数 16 个字节
}

type BusinessDataFinishRequest struct {
	AppType          OpType      `json:"appType,omitempty"`          // 操作类型 新办：1，延期
	CertType         CertType    `json:"certType,omitempty"`         // 证书类型 企业：2 ,个人：9
	SubCert          SubCert     `json:"subCert,omitempty"`          // 主子证标识 目前均为主证 主证：1 子证：0
	Source           string      `json:"source,omitempty"`           // 第三方来源标识（可在联调时约定）
	PaymentMode      PaymentMode `json:"paymentMode,omitempty"`      // 支付方式 付款方式 1现金 6微信 7支付宝，其他情况默认现金
	ClientName       string      `json:"clientName,omitempty"`       // 证书名称
	LegalPerson      string      `json:"legalPerson,omitempty"`      // 法人，个人必填，企业非必填
	LegalNumber      string      `json:"legalNumber,omitempty"`      // 法人证件号码，个人必填，企业非必填
	LegalPersonTel   string      `json:"legalPersonTel,omitempty"`   // 法人电话，个人必填，企业非必填
	Agent            string      `json:"agent,omitempty"`            // 经办人
	AgentTel         string      `json:"agentTel,omitempty"`         // 经办人电话
	AgentNumber      string      `json:"agentNumber,omitempty"`      // 经办人证件号码
	CertId           string      `json:"certId,omitempty"`           // 证书序列号，延期、注销、变更时必填，除新办，丢补外均要传入
	CustomerType     string      `json:"customerType,omitempty"`     // 行业类型 联调对接时可分配
	SocialCreditCode string      `json:"socialCreditCode,omitempty"` // 社会信用代码（单位/业务证书必须传入）
	CountryName      string      `json:"countryName,omitempty"`      // 国家（证书C项）
	OrganizationName string      `json:"organizationName,omitempty"` // 组织机构名称（证书O项）
	Province         string      `json:"province,omitempty"`         // 省份（证书L项）
	City             string      `json:"city,omitempty"`             // 城市（证书S项）
	Department       string      `json:"department,omitempty"`       // 部门（证书OU项）
	DelayType        string      `json:"delayType,omitempty"`        // 延期时间选择标识 0 延期顺延（从之前证书截止日期开始），其它情况：延期非顺延（过期从当前时间延期，不过期顺延）
}

type BusinessDataFinishResponse struct {
	Result    string                         `json:"result,omitempty"`    // 业务数据接口操作结果
	ResultMsg string                         `json:"resultMsg,omitempty"` // 业务数据接口操作结果描述
	Data      BusinessDataFinishResponseData `json:"data,omitempty"`
}

type BusinessDataFinishResponseData struct {
	TrustId    string `json:"trustId,omitempty"`    // 客服信任号
	AppId      string `json:"appId,omitempty"`      // 业务单号
	ClientName string `json:"clientName,omitempty"` // 证书名称
}

type ApplySealCertRequest struct {
	TokenInfo  string `json:"tokenInfo,omitempty"`  // 客户端随机数（RandomA16位）与服务端 随机数（RandomB）字符串拼接，randomA+randomB+"SNCA"+appId（归档返回）
	CommonName string `json:"commonName,omitempty"` // 证书主体的名称项，不大于 200 个字节
	P10        string `json:"p10,omitempty"`        // P10编码的证书请求
}

type ApplySealCertResponse struct {
	TaskCode      string                    `json:"taskCode,omitempty"`      // 业务类型 applySealCert
	ResultCode    string                    `json:"resultCode,omitempty"`    // 返回值，0 表示成功，其他表示错误
	ResultCodeMsg string                    `json:"resultCodeMsg,omitempty"` // 返回信息描述
	Data          ApplySealCertResponseData `json:"data,omitempty"`
}

type ApplySealCertResponseData struct {
	SignCert         string `json:"signCert,omitempty"`         // Base64 编码签名证书
	EncCert          string `json:"encCert,omitempty"`          // Base64 编码的加密证书
	EncKeyProtection string `json:"encKeyProtection,omitempty"` // Base64 编码的加密密钥对保护结构，被编码的结构 ENVELOPEDKEYBLOB，遵从GM/T 0016—2012 标准
	TransactionID    string `json:"transactionID,omitempty"`    // transactionID
}
