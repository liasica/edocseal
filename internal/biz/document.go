// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package biz

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
)

// NewDocumentPaths 生成文档目录
func NewDocumentPaths() *model.DocumentPaths {
	docId := time.Now().Format("200601") + g.GetID()
	// 创建文件夹
	dir := getDocumentDir(docId)
	_ = edocseal.CreateDirectory(dir)
	return getDocumentPaths(dir, docId)
}

func getDocumentDir(docId string) string {
	dir, _ := filepath.Abs(filepath.Join(g.GetDocumentDir(), docId[:4], docId[4:6], docId[6:]))
	return dir
}

func getDocumentPaths(dir, docId string) *model.DocumentPaths {
	return &model.DocumentPaths{
		ID:                  docId,
		OssUnSignedDocument: filepath.Join("__contracts", docId[:4], docId[4:6], docId[6:], docId+"_unsigned.pdf"),
		OssSignedDocument:   filepath.Join("__contracts", docId[:4], docId[4:6], docId[6:], docId+".pdf"),
		Directory:           dir,
		UnSignedDocument:    filepath.Join(dir, docId+"_unsigned.pdf"),
		SignedDocument:      filepath.Join(dir, docId+".pdf"),
		Image:               filepath.Join(dir, "signature.png"),
		Config:              filepath.Join(dir, "config.json"),
		Cert:                filepath.Join(dir, "cert.pem"),
		Key:                 filepath.Join(dir, "key.pem"),
	}
}

// GetDocumentPaths 获取文档目录
func GetDocumentPaths(docId string) (*model.DocumentPaths, error) {
	dir, _ := filepath.Abs(filepath.Join(g.GetDocumentDir(), docId[:4], docId[4:6], docId[6:]))
	if !edocseal.FileExists(dir) {
		return nil, fmt.Errorf("文档不存在: %s", docId)
	}
	return getDocumentPaths(dir, docId), nil
}
