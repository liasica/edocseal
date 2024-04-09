// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package biz

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/signintech/pdft"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
)

const (
	FontSong = "Song"
)

var (
	// 模板配置列表
	templates sync.Map

	// 字体列表
	fontlist = []string{
		FontSong,
	}
)

var fonts sync.Map

// LoadFonts 加载所有字体
func LoadFonts(path string) error {
	for _, f := range fontlist {
		file := filepath.Join(path, f+".ttf")
		// 读取字体
		name := edocseal.FileNameWithoutExtension(file)
		data, err := pdft.PDFParseFont(file, name)
		if err != nil {
			return err
		}
		fonts.Store(name, data)
	}
	return nil
}

// AddFonts 添加所有字体
func AddFonts(ft *pdft.PDFt) {
	fonts.Range(func(key, value any) bool {
		_ = ft.AddFontFromData(key.(string), value.(*pdft.PDFFontData))
		return true
	})
}

// GetTemplate 根据ID返回模板路径以及配置
func GetTemplate(id string) (*model.TemplateData, error) {
	data, ok := templates.Load(id)
	if !ok {
		var err error
		data, err = loadTemplateConfig(id)
		if err != nil {
			return nil, err
		}
	}
	return data.(*model.TemplateData), nil
}

func loadTemplateConfig(id string) (*model.TemplateData, error) {
	cfgPath := filepath.Join(g.GetTemplateDir(), id+".json")
	path := filepath.Join(g.GetTemplateDir(), id+".pdf")
	if !edocseal.FileExists(cfgPath) || !edocseal.FileExists(path) {
		return nil, errors.New("模板未找到")
	}

	data := &model.TemplateData{
		Path: path,
	}
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	_ = jsoniter.Unmarshal(b, &data.Fields)

	return data, nil
}
