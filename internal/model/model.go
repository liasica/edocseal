// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-10, by liasica

package model

var (
	A4PageSize = PageSize{Width: 595, Height: 842}
)

type PageSize struct {
	Width  float64 `json:"width,omitempty"`  // 宽度
	Height float64 `json:"height,omitempty"` // 高度
}
