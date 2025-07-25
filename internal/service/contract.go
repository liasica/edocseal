// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package service

import (
	"context"

	"go.uber.org/zap"

	"auroraride.com/edocseal/internal/biz"
	"auroraride.com/edocseal/internal/ent"
	"auroraride.com/edocseal/internal/task"
	"auroraride.com/edocseal/pb"
)

// GracefulStartGrpcServer TODO: 优雅启动和停止grpc服务
func GracefulStartGrpcServer() {
}

type ContractService struct {
	pb.UnimplementedContractServiceServer
}

// Create 创建合同
func (*ContractService) Create(_ context.Context, req *pb.ContractServiceCreateRequest) (*pb.ContractServiceCreateResponse, error) {
	res := &pb.ContractServiceCreateResponse{}

	err := <-task.DocumentTask().AddJob(func() (err error) {
		level := zap.InfoLevel
		defer func() {
			fields := []zap.Field{
				zap.Reflect("payload", req),
				zap.Reflect("response", res),
			}
			if err != nil {
				level = zap.ErrorLevel
				fields = append(fields, zap.Error(err))
			}
			zap.L().Check(level, "生成文档").Write(fields...)
		}()

		// 创建合同
		var doc *ent.Document
		doc, err = biz.CreateDocument(req, true)
		if err != nil {
			return err
		}

		res.Url = doc.UnsignedURL
		res.DocId = doc.ID
		return
	})

	return res, err
}

// Sign 合同签署
func (*ContractService) Sign(_ context.Context, req *pb.ContractServiceSignRequest) (*pb.ContractServiceSignResponse, error) {
	res := &pb.ContractServiceSignResponse{}

	w := task.SignTask().AddJob(func() (err error) {
		level := zap.InfoLevel
		defer func() {
			fields := []zap.Field{
				zap.String("docId", req.DocId),
				zap.Reflect("payload", map[string]string{
					"name":     req.Name,
					"province": req.Province,
					"city":     req.City,
					"address":  req.Address,
					"phone":    req.Phone,
					"idcard":   req.Idcard,
				}),
				zap.Reflect("response", res),
			}
			if err != nil {
				level = zap.ErrorLevel
				fields = append(fields, zap.Error(err))
			}
			zap.L().Check(level, "签署合同").Write(fields...)
		}()

		var url string
		// 签署合同
		url, err = biz.SignDocument(req, true)
		if err != nil {
			return
		}

		res.Status = pb.SignStatus_SIGN_STATUS_SUCCESS
		res.Url = url
		return
	})

	err := <-w
	if err != nil {
		res.Status = pb.SignStatus_SIGN_STATUS_FAIL
		res.Message = err.Error()
	}

	return res, nil
}
