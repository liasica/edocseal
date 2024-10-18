// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package model

import (
	"fmt"
	"io"
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

	rs io.ReadSeeker // 模板内容
}

func (t *Template) LoadContent() (err error) {
	t.rs, err = os.Open(t.File)
	return
}

func (t *Template) ReadSeeker() io.ReadSeeker {
	return t.rs
}

// TemplateField 模板字段
type TemplateField struct {
	Page        int                `json:"page"`                  // 页码，从1开始
	Description string             `json:"description,omitempty"` // 字段描述
	Type        TemplateFieldType  `json:"type,omitempty"`        // 字段类型
	Rectangle   *TemplateRectangle `json:"rectangle,omitempty"`   // 字段位置
}
