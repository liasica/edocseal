// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package g

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
	"io"
)

var (
	//go:embed assets/check.png
	checkImageBytes []byte
	check           image.Image

	//go:embed assets/fonts/Song.ttf
	fontSongBytes []byte
	fontSong      io.Reader
)

// GetCheckImage 获取勾选图片
func GetCheckImage() image.Image {
	return check
}

func GetFontSong() (string, io.Reader) {
	return "Song", fontSong
}

func init() {
	var err error
	check, err = png.Decode(bytes.NewReader(checkImageBytes))
	if err != nil {
		panic(err)
	}

	fontSong = bytes.NewReader(fontSongBytes)
}
