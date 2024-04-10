// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-10, by liasica

package biz

import (
	"path/filepath"
	"sync"

	"github.com/signintech/pdft"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
)

const (
	FontSong = "Song"
)

var (
	once sync.Once

	fonts sync.Map

	// 字体列表
	fontlist = []string{
		FontSong,
	}
)

// 添加所有字体
func addFonts(ft *pdft.PDFt) (err error) {
	// 加载所有字体
	once.Do(func() {
		for _, f := range fontlist {
			file := filepath.Join(g.GetConfigPath(), f+".ttf")
			// 读取字体
			name := edocseal.FileNameWithoutExtension(file)
			var data *pdft.PDFFontData
			data, err = pdft.PDFParseFont(file, name)
			if err != nil {
				return
			}
			fonts.Store(name, data)
		}
	})

	if err != nil {
		return
	}

	fonts.Range(func(key, value any) bool {
		_ = ft.AddFontFromData(key.(string), value.(*pdft.PDFFontData))
		return true
	})

	return
}
