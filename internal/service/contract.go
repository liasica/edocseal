// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package service

import (
	"context"
	"os"

	"github.com/liasica/edocseal/internal/biz"
	"github.com/liasica/edocseal/internal/model"
	"github.com/liasica/edocseal/internal/task"
	"github.com/liasica/edocseal/pb"
)

// GracefulStartGrpcServer TODO: 优雅启动和停止grpc服务
func GracefulStartGrpcServer() {
}

type ContractService struct {
	pb.UnimplementedContractServer
}

// Create 创建合同
func (*ContractService) Create(_ context.Context, req *pb.ContractCreateRequest) (*pb.ContractCreateResponse, error) {
	res := &pb.ContractCreateResponse{}

	err := <-task.DocumentTask().AddJob(func() error {
		// 创建合同
		doc, paths, err := biz.CreateDocument(req.TemplateId, req.Values)
		if err != nil {
			return err
		}

		// 上传合同
		var url string
		url, err = biz.UploadDocument(paths.OssUnSignedDocument, doc)
		if err != nil {
			return err
		}

		res.Url = url
		res.DocId = paths.ID
		return nil
	})

	return res, err
}

// Sign 合同签署
func (*ContractService) Sign(_ context.Context, req *pb.ContractSignRequest) (*pb.ContractSignResponse, error) {
	res := &pb.ContractSignResponse{}

	w := task.SignTask().AddJob(func() (err error) {
		// 签署合同
		var paths *model.DocumentPaths
		paths, err = biz.SignDocument(req)
		if err != nil {
			return
		}

		// 读取合同
		var b []byte
		b, err = os.ReadFile(paths.SignedDocument)
		if err != nil {
			return
		}

		// 上传合同
		var url string
		url, err = biz.UploadDocument(paths.OssSignedDocument, b)
		if err != nil {
			return
		}

		res.Status = pb.ContractSignStatus_SUCCESS
		res.Url = url
		return
	})

	err := <-w
	if err != nil {
		res.Status = pb.ContractSignStatus_FAIL
		res.Message = err.Error()
	}

	return res, nil
}
