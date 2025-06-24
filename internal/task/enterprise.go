// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-18, by liasica

package task

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"auroraride.com/edocseal/ca"
	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/g"
)

type EnterpriseTask struct{}

func NewEnterpriseTask() *EnterpriseTask {
	return &EnterpriseTask{}
}

func (e *EnterpriseTask) Run() {
	c := cron.New()
	_, err := c.AddFunc("0 3 * * *", func() {
		zap.L().Info("开始执行企业证书检查过期任务")
		e.do()
	})
	if err != nil {
		zap.L().Fatal("企业证书检查过期任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

func (e *EnterpriseTask) do() {
	cfg := g.GetEnterpriseConfig()

	// 加载证书
	cert, err := ca.LoadCertificateFromFile(cfg.Certificate)
	if err != nil {
		zap.L().Error("加载企业证书失败", zap.Error(err))
		return
	}

	zap.L().Info("证书加载成功", zap.String("subject", cert.Subject.String()), zap.Time("notBefore", cert.NotBefore), zap.Time("notAfter", cert.NotAfter))

	// 证书过期检查（7天内过期）
	if cert.NotAfter.Before(time.Now().AddDate(0, 0, 7)) {
		err = biz.RequestEnterpriseCertAndUpdateConfig()
		if err != nil {
			zap.L().Error("更新证书失败", zap.Error(err))
		}
	}
}
