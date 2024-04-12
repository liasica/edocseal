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

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
)

var (
	// 模板配置列表
	templates sync.Map
)

// GetTemplate 根据ID返回模板路径以及配置
func GetTemplate(id string) (*model.Template, error) {
	data, ok := templates.Load(id)
	if !ok {
		var err error
		data, err = loadTemplateConfig(id)
		if err != nil {
			return nil, err
		}
	}
	return data.(*model.Template), nil
}

func loadTemplateConfig(id string) (*model.Template, error) {
	b, err := os.ReadFile(filepath.Join(g.GetTemplateDir(), id+".json"))
	if err != nil {
		return nil, err
	}

	template := new(model.Template)

	// 解析模板配置
	err = jsoniter.Unmarshal(b, &template)
	if err != nil {
		return nil, err
	}

	// 验证模板文档是否存在
	if !edocseal.FileExists(template.File) {
		return nil, errors.New("模板未找到")
	}

	return template, nil
}

func oss() (ao *edocseal.AliyunOss, url string, err error) {
	cfg := g.GetAliyunOss()
	ao, err = edocseal.NewAliyunOss(cfg.AccessKeyId, cfg.AccessKeySecret, cfg.Endpoint, cfg.Bucket)
	url = cfg.Url
	return
}
