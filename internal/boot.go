// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package internal

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"auroraride.com/edocseal"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/g"
	"auroraride.com/edocseal/third/snca"
)

func Boot() {
	// 替换全局日志
	zap.ReplaceGlobals(g.NewZap())

	// 验证signer version
	b, err := edocseal.Exec(g.GetSigner(), "--version")
	if err != nil {
		zap.L().Fatal("Python环境未安装或配置错误", zap.Error(err))
		os.Exit(1)
	}

	// 创建目录
	err = edocseal.CreateDirectory(g.GetTemplateDir())
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
	err = edocseal.CreateDirectory(g.GetCertificateDir())
	if err != nil {
		zap.L().Fatal("证书目录创建失败", zap.Error(err))
	}
	err = edocseal.CreateDirectory(filepath.Dir(g.GetBboltPath()))
	if err != nil {
		zap.L().Fatal("创建bbolt目录失败", zap.Error(err))
	}

	// 创建数据库
	dsn, dbDebug := g.GetPostgresConfig()
	err = ent.CreateDatabase(dsn, dbDebug)
	if err != nil {
		zap.L().Fatal("数据库连接失败", zap.Error(err))
	}

	// 初始化证书
	rootCrt := g.NewCertificate()

	// 初始化snca
	snca.Setup(g.GetSnca())

	zap.L().Info(
		"edocseal 初始化完成",
		zap.String("configFile", g.GetConfigFile()),
		zap.String("rpcBind", g.GetRPCBind()),
		zap.Bool("rootCrt", rootCrt.GetCertificate() != nil && rootCrt.GetPrivateKey() != nil),
		zap.String("signerVersion", string(b)),
	)
}
