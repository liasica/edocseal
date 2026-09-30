// Copyright (C) edocseal. 2024-present.
//
// Created at 2026-09-24, by liasica

package task

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"auroraride.com/edocseal/internal/biz"
)

// CertificateTask 企业根证书与企业证书自动续签
type CertificateTask struct{}

func NewCertificateTask() *CertificateTask {
	return &CertificateTask{}
}

// Run 每天凌晨 3 点续签即将到期的企业根证书与企业证书
func (*CertificateTask) Run() {
	c := cron.New()
	_, err := c.AddFunc("0 3 * * *", biz.RenewCertificates)
	if err != nil {
		zap.L().Fatal("企业证书自动续签任务启动失败", zap.Error(err))
		return
	}
	c.Start()
}
