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
	bolt "go.etcd.io/bbolt"

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

func oss() (ao *edocseal.AliyunOss, err error) {
	cfg := g.GetAliyunOss()
	ao, err = edocseal.NewAliyunOss(cfg.AccessKeyId, cfg.AccessKeySecret, cfg.Endpoint, cfg.Bucket)
	return
}

// CreateShortUrl 生成短链接
func CreateShortUrl(url string) (short string, err error) {
	var id string
	id, err = g.NewShortId().Generate()
	if err != nil {
		return
	}
	// 保存
	err = g.NewBolt().Update(func(tx *bolt.Tx) error {
		return tx.Bucket(g.ShortUrlBucket).Put([]byte(id), []byte(url))
	})
	if err == nil {
		short = g.GetShortUrlPrefix() + id
	}
	return
}

// GetShortUrl 获取短链接
func GetShortUrl(id string) (url string, err error) {
	err = g.NewBolt().View(func(tx *bolt.Tx) error {
		v := tx.Bucket(g.ShortUrlBucket).Get([]byte(id))
		if v == nil {
			return errors.New("短链接未找到")
		}
		url = string(v)
		return nil
	})
	return
}
