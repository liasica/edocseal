// Copyright (C) edocseal. 2025-present.
//
// Created at 2025-06-24, by liasica

package edocseal

import (
	"github.com/signintech/gopdf"
)

// 标题尺寸配置
// var titleRect = &model.TemplateRectangle{
// 	LeftBottom: model.Coordinate{
// 		X: 70,  // 左下角X坐标
// 		Y: 743, // 左下角Y坐标
// 	},
// 	RightTop: model.Coordinate{
// 		X: 524, // 右上角X坐标
// 		Y: 765, // 右上角Y坐标
// 	},
// }

// 页面边距配置
// 按照A4纸张尺寸设置边距
var pageMargin = &gopdf.Margins{
	Left:   70, // 2.5cm = 70pt, 合同实际是 2.5cm
	Top:    72, // 2.54cm = 72pt, 合同实际是 2.54cm
	Right:  70, // 2.54cm = 70pt, 合同实际是 1.91cm
	Bottom: 72, // 2.54cm = 72pt, 合同实际是 2.54cm
}

type Rect struct {
	X float64 // 开始X坐标
	Y float64 // 开始Y坐标
	W float64 // 矩形宽度
	H float64 // 矩形高度
}

func GetPageTitleRect(size *gopdf.Rect, fontSize float64) *Rect {
	return &Rect{
		X: pageMargin.Left,
		Y: pageMargin.Top,
		W: size.W - pageMargin.Left - pageMargin.Right,
		H: fontSize, // 标题高度等于字体大小
	}
}

// GetPageWithoutTitleRect 计算页面矩形，不包含标题区域
// 需要和标题高度有间距
func GetPageWithoutTitleRect(size *gopdf.Rect, fontSize float64, margin float64) *Rect {
	return &Rect{
		X: pageMargin.Left,
		Y: fontSize + pageMargin.Top + margin,
		W: size.W - pageMargin.Left - pageMargin.Right,
		H: size.H - pageMargin.Top - pageMargin.Bottom - fontSize - margin,
	}
}

// ScaleRect 缩放矩形适应容器
func ScaleRect(width, height, containerWidth, containerHeight float64) (scaledWidth, scaledHeight float64) {
	// 如果矩形已经在容器内，直接返回原始尺寸
	if width <= containerWidth && height <= containerHeight {
		return width, height
	}

	// 按宽度、高度分别计算缩放比例
	widthRatio := containerWidth / width
	heightRatio := containerHeight / height

	// 取较小的比例来保持长宽比
	scale := widthRatio
	if heightRatio < widthRatio {
		scale = heightRatio
	}

	scaledWidth = width * scale
	scaledHeight = height * scale
	return
}
