// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

type AgencyCertResponse struct {
	ErrorCode    string                  `json:"errorCode,omitempty"`    // 200为正常，500/ E300000为错误
	Result       string                  `json:"result,omitempty"`       // true为成功，false为失败
	ErrorMessage string                  `json:"errorMessage,omitempty"` // 200为正常，500/ E300000为错误
	JsonObj      AgencyCertJsonObjStruct `json:"jsonObj,omitempty"`      // 签名公钥证书，cer格式 obj
}

type AgencyCertJsonObjStruct struct {
	SignCert string `json:"signCert"` // 签名公钥证书，cer格式
}
