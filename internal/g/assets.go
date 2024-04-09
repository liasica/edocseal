// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package g

import (
	_ "embed"
)

var (
	//go:embed assets/check.png
	check []byte
)

// GetCheckImage 获取勾选图片
func GetCheckImage() []byte {
	return check
}
