// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package model

import (
	"fmt"
	"os"
)

// TemplateFieldType 模板字段类型
type TemplateFieldType string

const (
	TemplateFieldTypeText      TemplateFieldType = "text"      // 文本
	TemplateFieldTypeCheckbox  TemplateFieldType = "checkbox"  // 勾选
	TemplateFieldTypeSignature TemplateFieldType = "signature" // 签名
)

func (t TemplateFieldType) Code() string {
	switch t {
	case TemplateFieldTypeText:
		return "/Tx"
	case TemplateFieldTypeCheckbox:
		return "/Btn"
	case TemplateFieldTypeSignature:
		return "/Sig"
	}
	return ""
}

func NewTemplateFieldType(name string) (TemplateFieldType, error) {
	switch name {
	case "/Tx":
		return TemplateFieldTypeText, nil
	case "/Btn":
		return TemplateFieldTypeCheckbox, nil
	case "/Sig":
		return TemplateFieldTypeSignature, nil
	}
	return "", fmt.Errorf("未知的字段类型 ( %s )", name)
}

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// TemplateRectangle 模板矩形
// 用于描述PDF页面上的矩形区域，通常用于定位字段位置
// 注意：PDF坐标系的原点在左下角，Y轴向上，X轴向右
type TemplateRectangle struct {
	LeftBottom Coordinate `json:"lb,omitempty"` // 左下角坐标 (x,y), 实际可读计算时需要取页面高度减去该值
	RightTop   Coordinate `json:"rt,omitempty"` // 右上角坐标 (x,y), 实际可读计算时需要取页面高度减去该值
}

// Width 获取宽度
func (rect *TemplateRectangle) Width() float64 {
	return rect.RightTop.X - rect.LeftBottom.X
}

// Height 获取高度
func (rect *TemplateRectangle) Height() float64 {
	return rect.RightTop.Y - rect.LeftBottom.Y
}

// IntList 获取整数列表
func (rect *TemplateRectangle) IntList() []int {
	return []int{
		int(rect.LeftBottom.X),
		int(rect.LeftBottom.Y),
		int(rect.RightTop.X),
		int(rect.RightTop.Y),
	}
}

// Template 模板
type Template struct {
	ID     string                   `json:"id"`     // 模板ID
	File   string                   `json:"file"`   // 模板文件
	Fields map[string]TemplateField `json:"fields"` // 字段列表

	content []byte
}

func (t *Template) LoadContent() (err error) {
	t.content, err = os.ReadFile(t.File)
	return
}

func (t *Template) Content() []byte {
	return t.content
}

// TemplateField 模板字段
type TemplateField struct {
	Page        int                `json:"page"`                  // 页码，从1开始
	Description string             `json:"description,omitempty"` // 字段描述
	Type        TemplateFieldType  `json:"type,omitempty"`        // 字段类型
	Rectangle   *TemplateRectangle `json:"rectangle,omitempty"`   // 字段位置
}
