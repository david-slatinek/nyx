// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: schema.proto

package __

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
	RecommendService_RecommendDialog_FullMethodName  = "/schema.RecommendService/RecommendDialog"
	RecommendService_RecommendSummary_FullMethodName = "/schema.RecommendService/RecommendSummary"
)

// RecommendServiceClient is the client API for RecommendService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecommendServiceClient interface {
	RecommendDialog(ctx context.Context, in *RecommendRequestDialog, opts ...grpc.CallOption) (*RecommendResponseList, error)
	RecommendSummary(ctx context.Context, in *RecommendRequestSummary, opts ...grpc.CallOption) (*RecommendResponse, error)
}

type recommendServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecommendServiceClient(cc grpc.ClientConnInterface) RecommendServiceClient {
	return &recommendServiceClient{cc}
}

func (c *recommendServiceClient) RecommendDialog(ctx context.Context, in *RecommendRequestDialog, opts ...grpc.CallOption) (*RecommendResponseList, error) {
	out := new(RecommendResponseList)
	err := c.cc.Invoke(ctx, RecommendService_RecommendDialog_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recommendServiceClient) RecommendSummary(ctx context.Context, in *RecommendRequestSummary, opts ...grpc.CallOption) (*RecommendResponse, error) {
	out := new(RecommendResponse)
	err := c.cc.Invoke(ctx, RecommendService_RecommendSummary_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecommendServiceServer is the server API for RecommendService service.
// All implementations must embed UnimplementedRecommendServiceServer
// for forward compatibility
type RecommendServiceServer interface {
	RecommendDialog(context.Context, *RecommendRequestDialog) (*RecommendResponseList, error)
	RecommendSummary(context.Context, *RecommendRequestSummary) (*RecommendResponse, error)
	mustEmbedUnimplementedRecommendServiceServer()
}

// UnimplementedRecommendServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRecommendServiceServer struct {
}

func (UnimplementedRecommendServiceServer) RecommendDialog(context.Context, *RecommendRequestDialog) (*RecommendResponseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecommendDialog not implemented")
}
func (UnimplementedRecommendServiceServer) RecommendSummary(context.Context, *RecommendRequestSummary) (*RecommendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecommendSummary not implemented")
}
func (UnimplementedRecommendServiceServer) mustEmbedUnimplementedRecommendServiceServer() {}

// UnsafeRecommendServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecommendServiceServer will
// result in compilation errors.
type UnsafeRecommendServiceServer interface {
	mustEmbedUnimplementedRecommendServiceServer()
}

func RegisterRecommendServiceServer(s grpc.ServiceRegistrar, srv RecommendServiceServer) {
	s.RegisterService(&RecommendService_ServiceDesc, srv)
}

func _RecommendService_RecommendDialog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecommendRequestDialog)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommendServiceServer).RecommendDialog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecommendService_RecommendDialog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommendServiceServer).RecommendDialog(ctx, req.(*RecommendRequestDialog))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecommendService_RecommendSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecommendRequestSummary)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommendServiceServer).RecommendSummary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecommendService_RecommendSummary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommendServiceServer).RecommendSummary(ctx, req.(*RecommendRequestSummary))
	}
	return interceptor(ctx, in, info, handler)
}

// RecommendService_ServiceDesc is the grpc.ServiceDesc for RecommendService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecommendService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "schema.RecommendService",
	HandlerType: (*RecommendServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RecommendDialog",
			Handler:    _RecommendService_RecommendDialog_Handler,
		},
		{
			MethodName: "RecommendSummary",
			Handler:    _RecommendService_RecommendSummary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "schema.proto",
}
