// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package biz

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/signintech/pdft"
	xpdf "github.com/signintech/pdft/minigopdf"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
	"github.com/liasica/edocseal/pb"
)

// GetDocumentFile 获取文档
func GetDocumentFile(docId string) (string, error) {
	path := filepath.Join(g.GetDocumentDir(), docId[:4], docId[4:6], docId+".pdf")
	if !edocseal.FileExists(path) {
		return "", edocseal.ErrDocumentNotFound(docId)
	}
	return path, nil
}

// GenerateDocumentId 生成文档目录，返回文档ID和目录地址
func GenerateDocumentId() (docId, path string) {
	prefix := time.Now().Format("200601")
	docId = g.GetID()
	// 创建文件夹
	path = filepath.Join(g.GetDocumentDir(), prefix[:4], prefix[4:], docId)
	_ = edocseal.CreateDirectory(path)
	return
}

// CreateDocument 根据模板创建待签约文档
func CreateDocument(templateId string, fields map[string]*pb.ContractFromField) (b []byte, docId string, err error) {
	// 获取模板和配置
	var tmpl *model.TemplateData
	tmpl, err = GetTemplate(templateId)
	if err != nil {
		return
	}

	var path string
	docId, path = GenerateDocumentId()

	// 生成文档ID
	dst := filepath.Join(path, docId+".pdf")

	// 复制模板
	err = edocseal.FileCopy(tmpl.Path, dst)
	if err != nil {
		return
	}

	// 创建PDF
	pt := new(pdft.PDFt)

	// 打开PDF
	err = pt.Open(dst)
	if err != nil {
		return
	}

	// 加载所有字体
	err = addFonts(pt)
	if err != nil {
		return
	}

	err = pt.SetFont(FontSong, "", 10)
	if err != nil {
		return
	}

	size := model.A4PageSize

	// 填充字段
	for k, v := range fields {
		c, ok := tmpl.Fields[k]
		if !ok {
			err = edocseal.ErrFieldNotFound(k)
			return
		}

		rect := c.Rectangle
		x := rect.LeftBottom.X
		y := rect.RightTop.Y

		switch val := v.Value.(type) {
		case *pb.ContractFromField_Text:
			err = pt.Insert(val.Text, c.Page, x, size.Height-y, rect.Width(), rect.Height(), xpdf.Left|xpdf.Middle)
		case *pb.ContractFromField_Checkbox:
			err = pt.InsertImg(g.GetCheckImage(), c.Page, x, size.Height-y, rect.Width(), rect.Height())
		}
		if err != nil {
			return
		}
	}

	// 创建目录
	err = edocseal.CreateDirectory(filepath.Dir(dst))
	if err != nil {
		return
	}

	// 文档写入缓冲区
	buf := new(bytes.Buffer)
	err = pt.SaveTo(buf)
	if err != nil {
		return
	}

	// 保存文档
	b = buf.Bytes()
	err = os.WriteFile(dst, b, os.ModePerm)
	return
}

// SignDocument 文档签约
func SignDocument(req *pb.ContractSignRequest) (err error) {
	var doc string
	doc, err = GetDocumentFile(req.DocId)

	// 获取证书
	RequestCertificae(filepath.Dir(doc))
}

// UploadDocument 上传至阿里云
func UploadDocument(path string, b []byte) (url string, err error) {
	// 获取OSS配置
	var ao *edocseal.AliyunOss
	ao, url, err = oss()
	if err != nil {
		return
	}

	// 上传至oss
	err = ao.UploadBytes(path, b)
	if err != nil {
		return
	}
	url = url + "/" + path
	return
}
