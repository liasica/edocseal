// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package g

import (
	"embed"
	_ "embed"
	"io/fs"
)

const (
	FontSong = "Song"
)

var (
	//go:embed assets/check.png
	check []byte

	//go:embed assets/fonts/*
	fonts embed.FS
)

// GetCheckImage 获取勾选图片
func GetCheckImage() []byte {
	return check
}

func GetFont(name string) (fs.File, error) {
	return fonts.Open("assets/fonts/" + name + ".ttf")
}
