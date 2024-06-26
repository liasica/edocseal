// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

// // DocumentPaths 文档路径，均为完整路径
// type DocumentPaths struct {
// 	ID                  string // 文档ID
// 	Directory           string // 文档目录
// 	OssUnSignedDocument string // Oss未签名文档路径
// 	OssSignedDocument   string // Oss签名文档路径
// 	UnSignedDocument    string // 未签名文档路径
// 	SignedDocument      string // 签名文档路径
// 	Image               string // 手写签名路径
// 	Config              string // 签名字段配置文件路径
// 	Cert                string // 证书路径
// 	Key                 string // 私钥路径
// }

// Paths 文档路径
type Paths struct {
	UnSigned    string `json:"unSigned"`    // 待签约文档路径
	Signed      string `json:"signed"`      // 已签约文档路径
	Image       string `json:"image"`       // 手写签名路径
	OssUnSigned string `json:"ossUnSigned"` // 待签约文档OSS路径
	OssSigned   string `json:"ossSigned"`   // 已签约文档OSS路径
}
