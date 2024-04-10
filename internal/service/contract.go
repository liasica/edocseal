// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/liasica/edocseal/pb"
)

type ContractService struct {
	pb.UnimplementedContractServer
}

// Create 创建合同
func (*ContractService) Create(ctx context.Context, req *pb.ContractCreateRequest) (*pb.ContractCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

// Sign 合同签署
func (*ContractService) Sign(context.Context, *pb.ContractSignRequest) (*pb.ContractSignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sign not implemented")
}
