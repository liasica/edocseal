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

func oss() (ao *edocseal.AliyunOss, url string, err error) {
	cfg := g.GetAliyunOss()
	ao, err = edocseal.NewAliyunOss(cfg.AccessKeyId, cfg.AccessKeySecret, cfg.Endpoint, cfg.Bucket)
	url = cfg.Url
	return
}
