// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

type AgencyCertResponse struct {
	TaskCode      string               `json:"taskCode,omitempty"`      // 业务类型 applySealCert
	ResultCode    string               `json:"resultCode,omitempty"`    // 返回值，0 表示成功，其他表示错误
	ResultCodeMsg string               `json:"resultCodeMsg,omitempty"` // 返回信息描述
	Data          AgencyCertDataStruct `json:"data,omitempty"`
}

type AgencyCertDataStruct struct {
	SignCert         string `json:"signCert"`         // Base64 编码签名证书
	EncCert          string `json:"encCert"`          // Base64 编码的加密证书
	EncKeyProtection string `json:"encKeyProtection"` // Base64 编码的加密密钥对保护结构，被编码的结构 ENVELOPEDKEYBLOB，遵从GM/T 0016—2012 标准
	TransactionID    string `json:"transactionID"`    // transactionID
}

type RandomBResponse struct {
	ResultCode    string `json:"resultCode,omitempty"`    // 返回值，0 表示成功，其他表示错误
	ResultCodeMsg string `json:"resultCodeMsg,omitempty"` // 返回信息描述
	RandomB       string `json:"randomB,omitempty"`       // 服务端随机数 16 个字节
}

type BusResponse struct {
	Result    string        `json:"result,omitempty"`    // 业务数据接口操作结果
	ResultMsg string        `json:"resultMsg,omitempty"` // 业务数据接口操作结果描述
	Data      BusDataStruct `json:"data,omitempty"`
}

type BusDataStruct struct {
	TrustId    string `json:"trustId,omitempty"`    // 客服信任号
	AppId      string `json:"appId,omitempty"`      // 业务单号
	ClientName string `json:"clientName,omitempty"` // 证书名称
}
