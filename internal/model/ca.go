// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

type AgencyCertResponse struct {
	ErrorCode    string        `json:"errorCode"`    // 200为正常，500/ E300000为错误
	Result       string        `json:"result"`       // true为成功，false为失败
	ErrorMessage string        `json:"errorMessage"` // 200为正常，500/ E300000为错误
	JsonObj      JsonObjStruct `json:"jsonObj"`      // 签名公钥证书，cer格式 obj
}

type JsonObjStruct struct {
	SignCert string `json:"signCert"` // 签名公钥证书，cer格式
}
