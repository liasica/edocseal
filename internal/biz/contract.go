// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package biz

import (
	"bytes"
	"encoding/base64"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/signintech/pdft"
	xpdf "github.com/signintech/pdft/minigopdf"
	"go.uber.org/zap"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
	"github.com/liasica/edocseal/pb"
)

const (
	entSignField   = "entSign"
	riderSignField = "riderSign"
)

// CreateDocument 根据模板创建待签约文档
func CreateDocument(templateId string, fields map[string]*pb.ContractFromField) (b []byte, docId string, err error) {
	// 获取模板和配置
	var tmpl *model.TemplateData
	tmpl, err = GetTemplate(templateId)
	if err != nil {
		return
	}

	paths := NewDocumentPaths()

	// 复制模板
	err = edocseal.FileCopy(tmpl.Path, paths.UnSignedDocument)
	if err != nil {
		return
	}

	// 创建PDF
	pt := new(pdft.PDFt)

	// 打开PDF
	err = pt.Open(paths.UnSignedDocument)
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

	// 文档写入缓冲区
	buf := new(bytes.Buffer)
	err = pt.SaveTo(buf)
	if err != nil {
		return
	}

	// 保存文档
	b = buf.Bytes()
	err = os.WriteFile(paths.UnSignedDocument, b, os.ModePerm)
	if err != nil {
		return
	}

	// 保存配置
	sb, _ := jsoniter.Marshal(&model.Sign{
		TemplateID: templateId,
		InFile:     paths.UnSignedDocument,
		OutFile:    paths.SignedDocument,
		Signatures: []model.Signature{
			{
				Field: entSignField,
				Image: g.GetSeal(),
				Key:   g.GetPrivateKey(),
				Cert:  g.GetCertificate(),
				Rect:  tmpl.Fields[entSignField].Rectangle.IntList(),
			},
			{
				Field: riderSignField,
				Image: paths.Image,
				Key:   paths.Key,
				Cert:  paths.Cert,
				Rect:  tmpl.Fields[riderSignField].Rectangle.IntList(),
			},
		},
	})
	err = os.WriteFile(paths.Config, sb, os.ModePerm)

	docId = paths.ID
	return
}

// SignDocument 文档签约
func SignDocument(req *pb.ContractSignRequest) (err error) {
	// 获取文档链接
	var paths *model.DocumentPaths
	paths, err = GetDocumentPaths(req.DocId)
	if err != nil {
		return
	}

	// 保存签名
	var img []byte
	img, err = base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		return
	}
	err = os.WriteFile(paths.Image, img, os.ModePerm)
	if err != nil {
		return
	}

	// 获取证书
	err = RequestCertificae(paths, req.Name, req.Province, req.City, req.Address, req.Phone, req.Idcard)
	if err != nil {
		return
	}

	// 调用签名
	var out []byte
	out, err = edocseal.Exec(g.GetSigner(), "--config", paths.Config)
	if err != nil {
		zap.L().Error("签名失败", zap.Error(err), zap.Reflect("payload", req), zap.String("output", string(out)))
		return err
	}
	zap.L().Info("签名成功", zap.String("docId", req.DocId))
	return nil
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
