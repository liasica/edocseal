// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package model

// TemplateFieldType 模板字段类型
type TemplateFieldType string

const (
	TemplateFieldTypeText      TemplateFieldType = "text"      // 文本
	TemplateFieldTypeCheckbox  TemplateFieldType = "checkbox"  // 勾选
	TemplateFieldTypeSignature TemplateFieldType = "signature" // 签名
)

type TemplateRectangle struct {
	LeftBottom [2]float64 `json:"lb,omitempty"` // 左下角坐标 (x,y)
	RightTop   [2]float64 `json:"rt,omitempty"` // 右上角坐标 (x,y)
}

// TemplateField 模板字段
type TemplateField struct {
	Name        string            `json:"name,omitempty"`        // 字段名称
	Description string            `json:"description,omitempty"` // 字段描述
	Type        TemplateFieldType `json:"type,omitempty"`        // 字段类型
	Rectangle   TemplateRectangle `json:"rectangle,omitempty"`   // 字段位置
}

// TemplateData 模板数据
type TemplateData struct {
	Path   string          `json:"path"`
	Fields []TemplateField `json:"fields"`
}
