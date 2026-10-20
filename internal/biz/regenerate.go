// Copyright (C) edocseal. 2026-present.
//
// Created at 2026-10-20, by liasica

package biz

import (
	"context"
	"os"

	"go.uber.org/zap"

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/ent/document"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/pb"
)

// RegenerateDocument 按模板重新生成文档并以指定签约企业签署，企业证书按 enterpriseIssuer 选取
// 签署后的文档地址直接使用 OSS 地址，短链接库由服务进程独占，独立进程无法写入
func RegenerateDocument(
	creditCode string,
	enterpriseIssuer string,
	create *pb.ContractServiceCreateRequest,
	sign *pb.ContractServiceSignRequest,
) (doc *ent.Document, err error) {
	var ep *ent.Enterprise
	ep, err = QueryEnterprise(creditCode)
	if err != nil {
		return
	}

	doc, err = CreateDocument(create, false)
	if err != nil {
		return
	}

	sign.DocId = doc.ID

	var img []byte
	img, err = signWithEnterprise(doc, sign, ep, enterpriseIssuer)
	if err != nil {
		return
	}

	// 上传合同与签名图片
	var b []byte
	b, err = os.ReadFile(doc.Paths.Signed)
	if err != nil {
		return
	}

	var ao *edocseal.AliyunOss

	ao, err = oss()
	if err != nil {
		return
	}

	err = ao.UploadBytes(doc.Paths.OssSigned, b)
	if err != nil {
		return
	}

	_ = ao.UploadBytes(doc.Paths.OssImage, img)

	url := g.GetAliyunOss().GetUrl(doc.Paths.OssSigned)

	doc, err = doc.Update().SetStatus(document.StatusSigned).SetSignedURL(url).Save(context.Background())
	if err != nil {
		zap.L().Error("更新文档状态失败 → 重新生成", zap.Error(err), zap.String("docId", sign.DocId), zap.String("url", url))
	}
	return
}
