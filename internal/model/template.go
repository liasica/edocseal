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

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type TemplateRectangle struct {
	LeftBottom Coordinate `json:"lb,omitempty"` // 左下角坐标 (x,y)
	RightTop   Coordinate `json:"rt,omitempty"` // 右上角坐标 (x,y)
}

// Width 获取宽度
func (rect *TemplateRectangle) Width() float64 {
	return rect.RightTop.X - rect.LeftBottom.X
}

// Height 获取高度
func (rect *TemplateRectangle) Height() float64 {
	return rect.RightTop.Y - rect.LeftBottom.Y
}

// TemplateField 模板字段
type TemplateField struct {
	Page        int                `json:"page"`                  // 页码，从1开始
	Description string             `json:"description,omitempty"` // 字段描述
	Type        TemplateFieldType  `json:"type,omitempty"`        // 字段类型
	Rectangle   *TemplateRectangle `json:"rectangle,omitempty"`   // 字段位置
}

// TemplateData 模板数据
type TemplateData struct {
	Path   string                   `json:"path"`
	Fields map[string]TemplateField `json:"fields"`
}
