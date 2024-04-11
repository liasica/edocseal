// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: contract.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Contract_Create_FullMethodName = "/pb.Contract/create"
	Contract_Sign_FullMethodName   = "/pb.Contract/sign"
)

// ContractClient is the client API for Contract service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContractClient interface {
	// 创建合同
	Create(ctx context.Context, in *ContractCreateRequest, opts ...grpc.CallOption) (*ContractCreateResponse, error)
	// 合同签署
	Sign(ctx context.Context, in *ContractSignRequest, opts ...grpc.CallOption) (*ContractSignResponse, error)
}

type contractClient struct {
	cc grpc.ClientConnInterface
}

func NewContractClient(cc grpc.ClientConnInterface) ContractClient {
	return &contractClient{cc}
}

func (c *contractClient) Create(ctx context.Context, in *ContractCreateRequest, opts ...grpc.CallOption) (*ContractCreateResponse, error) {
	out := new(ContractCreateResponse)
	err := c.cc.Invoke(ctx, Contract_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contractClient) Sign(ctx context.Context, in *ContractSignRequest, opts ...grpc.CallOption) (*ContractSignResponse, error) {
	out := new(ContractSignResponse)
	err := c.cc.Invoke(ctx, Contract_Sign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContractServer is the server API for Contract service.
// All implementations must embed UnimplementedContractServer
// for forward compatibility
type ContractServer interface {
	// 创建合同
	Create(context.Context, *ContractCreateRequest) (*ContractCreateResponse, error)
	// 合同签署
	Sign(context.Context, *ContractSignRequest) (*ContractSignResponse, error)
	mustEmbedUnimplementedContractServer()
}

// UnimplementedContractServer must be embedded to have forward compatible implementations.
type UnimplementedContractServer struct {
}

func (UnimplementedContractServer) Create(context.Context, *ContractCreateRequest) (*ContractCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedContractServer) Sign(context.Context, *ContractSignRequest) (*ContractSignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sign not implemented")
}
func (UnimplementedContractServer) mustEmbedUnimplementedContractServer() {}

// UnsafeContractServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContractServer will
// result in compilation errors.
type UnsafeContractServer interface {
	mustEmbedUnimplementedContractServer()
}

func RegisterContractServer(s grpc.ServiceRegistrar, srv ContractServer) {
	s.RegisterService(&Contract_ServiceDesc, srv)
}

func _Contract_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContractCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContractServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Contract_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContractServer).Create(ctx, req.(*ContractCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Contract_Sign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContractSignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContractServer).Sign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Contract_Sign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContractServer).Sign(ctx, req.(*ContractSignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Contract_ServiceDesc is the grpc.ServiceDesc for Contract service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Contract_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Contract",
	HandlerType: (*ContractServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "create",
			Handler:    _Contract_Create_Handler,
		},
		{
			MethodName: "sign",
			Handler:    _Contract_Sign_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contract.proto",
}
