// Copyright (C) edocseal. 2024-present.
//
// Created at 2026-09-23, by liasica

package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/pb"
)

type EnterpriseService struct {
	pb.UnimplementedEnterpriseServiceServer
}

// 将业务错误转为 gRPC 错误，错误信息直接展示给后台
func enterpriseError(action string, err error) error {
	if err == nil {
		return nil
	}

	zap.L().Error(action+"失败", zap.Error(err))
	return status.Error(codes.FailedPrecondition, err.Error())
}

// List 签约企业列表
func (*EnterpriseService) List(_ context.Context, _ *pb.EnterpriseListRequest) (*pb.EnterpriseListResponse, error) {
	return &pb.EnterpriseListResponse{Items: biz.ListEnterprises()}, nil
}

// Save 新增或更新签约企业
func (*EnterpriseService) Save(_ context.Context, req *pb.EnterpriseSaveRequest) (*pb.EnterpriseEmptyResponse, error) {
	zap.L().Info("保存签约企业", zap.String("creditCode", req.CreditCode), zap.String("name", req.Name), zap.Bool("seal", len(req.Seal) > 0))
	return &pb.EnterpriseEmptyResponse{}, enterpriseError("保存签约企业", biz.SaveEnterprise(req))
}

// Delete 删除签约企业
func (*EnterpriseService) Delete(_ context.Context, req *pb.EnterpriseCodeRequest) (*pb.EnterpriseEmptyResponse, error) {
	zap.L().Info("删除签约企业", zap.String("creditCode", req.CreditCode))
	return &pb.EnterpriseEmptyResponse{}, enterpriseError("删除签约企业", biz.DeleteEnterprise(req.CreditCode))
}

// SetDefault 设为签约企业
func (*EnterpriseService) SetDefault(_ context.Context, req *pb.EnterpriseCodeRequest) (*pb.EnterpriseEmptyResponse, error) {
	zap.L().Info("设置签约企业", zap.String("creditCode", req.CreditCode))
	return &pb.EnterpriseEmptyResponse{}, enterpriseError("设置签约企业", biz.SetDefaultEnterprise(req.CreditCode))
}

// RevokeCertificates 作废企业证书
func (*EnterpriseService) RevokeCertificates(_ context.Context, req *pb.EnterpriseCodeRequest) (*pb.EnterpriseEmptyResponse, error) {
	zap.L().Info("作废企业证书", zap.String("creditCode", req.CreditCode))
	return &pb.EnterpriseEmptyResponse{}, enterpriseError("作废企业证书", biz.RevokeEnterpriseCertificates(req.CreditCode))
}

// RegenerateRootCertificate 重新生成企业自签根证书
func (*EnterpriseService) RegenerateRootCertificate(_ context.Context, req *pb.EnterpriseCodeRequest) (*pb.EnterpriseEmptyResponse, error) {
	zap.L().Info("重新生成企业根证书", zap.String("creditCode", req.CreditCode))
	return &pb.EnterpriseEmptyResponse{}, enterpriseError("重新生成企业根证书", biz.RegenerateRootCertificate(req.CreditCode))
}
