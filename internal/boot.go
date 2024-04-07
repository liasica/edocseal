// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package internal

import (
	"go.uber.org/zap"

	"github.com/liasica/edocseal/internal/g"
)

func Boot() {
	// 替换全局日志
	zap.ReplaceGlobals(g.NewZap())

	// 初始化证书
	g.NewCertificate()

	zap.L().Info(
		"edocseal 初始化完成成功",
		zap.String("cfgPath", g.GetConfigPath()),
		zap.String("rpcBind", g.GetRPCBind()),
		zap.Bool("rootCA", g.GetRootCertificate() != nil),
		zap.Bool("rootPrivateKey", g.GetRootPrivateKey() != nil),
	)
}