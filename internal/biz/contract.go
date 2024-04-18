// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package biz

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/signintech/pdft"
	xpdf "github.com/signintech/pdft/minigopdf"
	"go.uber.org/zap"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/ent"
	"github.com/liasica/edocseal/internal/ent/document"
	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/model"
	"github.com/liasica/edocseal/pb"
)

// CreateDocument 根据模板创建待签约文档
func CreateDocument(req *pb.ContractCreateRequest, upload bool) (doc *ent.Document, err error) {
	expire := time.Unix(req.Expire, 0)

	// 查询文档防止重复创建
	jc := jsoniter.Config{SortMapKeys: true}.Froze()
	reqBytes, _ := jc.Marshal(req.Values)
	hasher := md5.New()
	// 写入模板ID
	hasher.Write([]byte(req.TemplateId))
	// 写入身份证号
	hasher.Write([]byte(req.Idcard))
	// 写入过期时间
	hasher.Write([]byte(strconv.FormatInt(req.Expire, 10)))
	// 写入字段
	hasher.Write(reqBytes)
	hash := hex.EncodeToString(hasher.Sum(nil))

	// 查找是否已存在未失效的未签约合同
	doc, _ = ent.NewDatabase().Document.Query().
		Where(
			document.Hash(hash),
			document.IDCardNumber(req.Idcard),
			document.StatusIn(document.StatusUnsigned),
		).
		First(context.Background())
	if doc != nil {
		return
	}

	// 获取模板和配置
	var tmpl *model.Template
	tmpl, err = GetTemplate(req.TemplateId)
	if err != nil {
		return
	}

	// 生成文档ID
	docId := time.Now().Format("20060102") + g.GetID()

	paths := NewPaths(docId)

	// 创建文件夹
	err = edocseal.CreateDirectory(filepath.Dir(paths.UnSigned))
	if err != nil {
		return
	}

	// 创建PDF
	pt := new(pdft.PDFt)

	// 打开PDF
	err = pt.Open(tmpl.File)
	if err != nil {
		return
	}

	// 加载字体
	var f fs.File
	f, err = g.GetFont(g.FontSong)
	if err != nil {
		return
	}
	err = pt.AddFontFrom(g.FontSong, f)
	if err != nil {
		return
	}
	err = pt.SetFont(g.FontSong, "", 10)
	if err != nil {
		return
	}

	size := model.A4PageSize

	// 填充字段
	for k, v := range req.Values {
		c, ok := tmpl.Fields[k]
		if !ok {
			zap.L().Warn("字段不存在", zap.String("field", k))
			continue
		}

		rect := c.Rectangle
		x := rect.LeftBottom.X
		y := rect.RightTop.Y

		switch val := v.Value.(type) {
		case *pb.ContractFromField_Text:
			var w float64
			w, err = pt.MeasureTextWidth(val.Text)
			if err != nil {
				return
			}
			err = pt.Insert(val.Text, c.Page, x, size.Height-y, w, rect.Height(), xpdf.Left|xpdf.Middle)
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
	b := buf.Bytes()
	err = os.WriteFile(paths.UnSigned, b, os.ModePerm)
	if err != nil {
		return
	}

	// 上传合同
	var url string
	if upload {
		url, err = UploadDocument(paths.OssUnSigned, b)
		if err != nil {
			return
		}
	}

	// 保存文档信息
	doc, err = ent.NewDatabase().Document.Create().
		SetID(docId).
		SetHash(hash).
		SetTemplateID(req.TemplateId).
		SetIDCardNumber(req.Idcard).
		SetStatus(document.StatusUnsigned).
		SetPaths(paths).
		SetExpiresAt(expire).
		SetUnsignedURL(url).
		SetCreateAt(time.Now()).
		Save(context.Background())

	return
}

// SignDocument 文档签约
func SignDocument(req *pb.ContractSignRequest, upload bool) (url string, err error) {
	doc, _ := QueryDocument(req.DocId)
	if doc == nil {
		return "", errors.New("未找到待签约文档")
	}

	if doc.Status == document.StatusSigned {
		return
	}

	// 获取模板
	var tmpl *model.Template
	tmpl, err = GetTemplate(doc.TemplateID)
	if err != nil {
		return
	}

	// 获取未签名文档
	if !edocseal.FileExists(doc.Paths.UnSigned) {
		return "", errors.New("未找到未签约文档")
	}

	// 保存签名
	var img []byte
	img, err = base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		return
	}
	err = os.WriteFile(doc.Paths.Image, img, os.ModePerm)
	if err != nil {
		return
	}

	// 获取证书
	var cert *ent.Certification
	cert, err = RequestCertificae(req.Name, req.Province, req.City, req.Address, req.Phone, req.Idcard)
	if err != nil {
		return
	}

	// 保存配置
	cfgPath, _ := filepath.Abs(filepath.Join(filepath.Dir(doc.Paths.UnSigned), "config.json"))
	kp, _ := filepath.Abs(cert.PrivatePath)
	cp, _ := filepath.Abs(cert.CertPath)

	signed, _ := filepath.Abs(doc.Paths.Signed)
	unsigned, _ := filepath.Abs(doc.Paths.UnSigned)
	imgPath, _ := filepath.Abs(doc.Paths.Image)

	ec := g.GetEnterpriseConfig()
	eSeal, _ := filepath.Abs(ec.Seal)
	eKey, _ := filepath.Abs(ec.PrivateKey)
	eCert, _ := filepath.Abs(ec.Certificate)

	sb, _ := jsoniter.Marshal(&model.Sign{
		TemplateID: doc.TemplateID,
		InFile:     unsigned,
		OutFile:    signed,
		Signatures: []model.Signature{
			{
				Field: model.EntSignField,
				Image: eSeal,
				Key:   eKey,
				Cert:  eCert,
				Rect:  tmpl.Fields[model.EntSignField].Rectangle.IntList(),
			},
			{
				Field: model.PersonalSignField,
				Image: imgPath,
				Key:   kp,
				Cert:  cp,
				Rect:  tmpl.Fields[model.PersonalSignField].Rectangle.IntList(),
			},
		},
	})
	err = os.WriteFile(cfgPath, sb, os.ModePerm)
	if err != nil {
		zap.L().Error("文档写入失败", zap.Error(err), zap.Reflect("payload", req), zap.Error(err))
		return
	}

	// 调用签名
	var out []byte
	out, err = edocseal.Exec(g.GetSigner(), "--config", cfgPath)
	if err != nil {
		zap.L().Error("签名失败", zap.Error(err), zap.Reflect("payload", req), zap.String("output", string(out)))
		return
	}

	// 上传至阿里云
	if upload {
		// 读取合同
		var b []byte
		b, err = os.ReadFile(doc.Paths.Signed)
		if err != nil {
			return
		}

		// 上传合同
		url, err = UploadDocument(doc.Paths.OssSigned, b)
		if err != nil {
			return
		}
	}

	zap.L().Info("签名成功", zap.String("docId", req.DocId), zap.String("url", url))

	// 更新数据库
	err = doc.Update().SetStatus(document.StatusSigned).SetSignedURL(url).Exec(context.Background())
	if err != nil {
		zap.L().Error("更新文档状态失败 → 签约成功", zap.Error(err), zap.String("docId", req.DocId), zap.String("url", url))
	}
	return
}

// UploadDocument 上传至阿里云
func UploadDocument(path string, b []byte) (url string, err error) {
	// 获取OSS配置
	var (
		ao *edocseal.AliyunOss
	)
	ao, err = oss()
	if err != nil {
		return
	}

	// 上传至oss
	err = ao.UploadBytes(path, b)
	if err != nil {
		return
	}
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	return CreateShortUrl(g.GetAliyunOss().GetUrl(path))
}
