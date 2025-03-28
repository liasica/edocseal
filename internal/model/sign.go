// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package model

const (
	EntSignField      = "aurSeal"
	PersonalSignField = "riderSeal"
)

type Sign struct {
	TemplateID string      `json:"templateId"`
	InFile     string      `json:"in_file"`
	OutFile    string      `json:"out_file"`
	Signatures []Signature `json:"signatures"`
}

type Signature struct {
	Field string `json:"field"` // 字段名
	Image string `json:"image"` // 签名图片
	Key   string `json:"key"`   // 私钥
	Cert  string `json:"cert"`  // 证书
	Rect  []int  `json:"rect"`  // 签名区域
	Page  int    `json:"page"`  // 签名所在页面，从0开始
}
