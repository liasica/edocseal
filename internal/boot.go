// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package internal

import (
	"go.uber.org/zap"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/g"
)

func Boot() {
	// 替换全局日志
	zap.ReplaceGlobals(g.NewZap())

	// 创建目录
	err := edocseal.CreateDirectory(g.GetTemplateDir())
	if err != nil {
		zap.L().Fatal("创建模板目录失败", zap.Error(err))
	}
	err = edocseal.CreateDirectory(g.GetRuntimeDir())
	if err != nil {
		zap.L().Fatal("创建运行时目录失败", zap.Error(err))
	}
	err = edocseal.CreateDirectory(g.GetDocumentDir())
	if err != nil {
		zap.L().Fatal("文档目录创建失败", zap.Error(err))
	}

	// 初始化证书
	rootCrt, entCrt := g.NewCertificate()

	zap.L().Info(
		"edocseal 初始化完成",
		zap.String("cfgPath", g.GetConfigPath()),
		zap.String("rpcBind", g.GetRPCBind()),
		zap.Bool("rootCrt", rootCrt.GetCertificate() != nil && rootCrt.GetPrivateKey() != nil),
		zap.Bool("entCrt", entCrt.GetCertificate() != nil && entCrt.GetPrivateKey() != nil),
	)
}
