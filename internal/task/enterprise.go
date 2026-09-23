// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-18, by liasica

package task

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"auroraride.com/edocseal/internal/biz"
)

type EnterpriseTask struct{}

func NewEnterpriseTask() *EnterpriseTask {
	return &EnterpriseTask{}
}

// Run 定时任务：每天凌晨3点执行企业证书检查过期任务
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
	renewed, err := biz.RenewEnterpriseCertificate()
	if err != nil {
		zap.L().Error("更新证书失败", zap.Error(err))
		return
	}

	zap.L().Info("企业证书检查过期任务完成", zap.Bool("renewed", renewed))
}
