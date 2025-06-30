// Copyright (C) edocseal. 2025-present.
//
// Created at 2025-06-24, by liasica

package edocseal

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func GetImageSize(path string) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var img image.Config
	img, _, err = image.DecodeConfig(f) // DecodeConfig 只解码图片头部，不加载像素
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}
