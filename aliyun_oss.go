// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-10, by liasica

package edocseal

import (
	"bytes"
	"encoding/base64"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"resty.dev/v3"
)

type AliyunOss struct {
	*oss.Client

	Bucket *oss.Bucket
}

func NewAliyunOss(accessKeyId, accessKeySecret, endpoint, bucket string) (*AliyunOss, error) {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	var bk *oss.Bucket
	bk, err = client.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	return &AliyunOss{
		Client: client,
		Bucket: bk,
	}, nil
}

// UploadUrlFile 从URL获取资源并上传
func (c *AliyunOss) UploadUrlFile(name string, url string) (err error) {
	var res *resty.Response
	res, err = resty.New().R().Get(url)
	if err != nil {
		return err
	}
	return c.UploadBytes(name, res.Bytes())
}

// UploadBase64 上传jpg图片
func (c *AliyunOss) UploadBase64(name string, b64 string) (err error) {
	var b []byte
	b, err = base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return
	}
	return c.UploadBytes(name, b)
}

// UploadBytes 上传文件
func (c *AliyunOss) UploadBytes(name string, b []byte) (err error) {
	err = c.Bucket.PutObject(name, bytes.NewReader(b))
	if err != nil {
		return
	}
	return
}

// SingleDelete 单文件删除
func (c *AliyunOss) SingleDelete(name string) (err error) {
	err = c.Bucket.DeleteObject(name)
	if err != nil {
		return
	}
	return
}

// MultiDelete 多文件删除
func (c *AliyunOss) MultiDelete(names []string) (err error) {
	// 将oss.DeleteObjectsQuiet设置为true，表示不返回删除结果。
	_, err = c.Bucket.DeleteObjects(names, oss.DeleteObjectsQuiet(true))
	if err != nil {
		return
	}
	return
}
