// grpc server is Kubelet and grpc client is Api Server

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.9
// source: proto/kubelet_for_apiserver.proto

package proto

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
	KubeletApiServerService_SayHello_FullMethodName  = "/kubelet_for_apiserver.kubeletApiServerService/SayHello"
	KubeletApiServerService_CreatePod_FullMethodName = "/kubelet_for_apiserver.kubeletApiServerService/CreatePod"
	KubeletApiServerService_DeletePod_FullMethodName = "/kubelet_for_apiserver.kubeletApiServerService/DeletePod"
	KubeletApiServerService_GetPod_FullMethodName    = "/kubelet_for_apiserver.kubeletApiServerService/GetPod"
)

// KubeletApiServerServiceClient is the client API for KubeletApiServerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KubeletApiServerServiceClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	CreatePod(ctx context.Context, in *ApplyPodRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	DeletePod(ctx context.Context, in *DeletePodRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	GetPod(ctx context.Context, in *GetPodRequest, opts ...grpc.CallOption) (*GetPodResponse, error)
}

type kubeletApiServerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKubeletApiServerServiceClient(cc grpc.ClientConnInterface) KubeletApiServerServiceClient {
	return &kubeletApiServerServiceClient{cc}
}

func (c *kubeletApiServerServiceClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, KubeletApiServerService_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubeletApiServerServiceClient) CreatePod(ctx context.Context, in *ApplyPodRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, KubeletApiServerService_CreatePod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubeletApiServerServiceClient) DeletePod(ctx context.Context, in *DeletePodRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, KubeletApiServerService_DeletePod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubeletApiServerServiceClient) GetPod(ctx context.Context, in *GetPodRequest, opts ...grpc.CallOption) (*GetPodResponse, error) {
	out := new(GetPodResponse)
	err := c.cc.Invoke(ctx, KubeletApiServerService_GetPod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KubeletApiServerServiceServer is the server API for KubeletApiServerService service.
// All implementations must embed UnimplementedKubeletApiServerServiceServer
// for forward compatibility
type KubeletApiServerServiceServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	CreatePod(context.Context, *ApplyPodRequest) (*StatusResponse, error)
	DeletePod(context.Context, *DeletePodRequest) (*StatusResponse, error)
	GetPod(context.Context, *GetPodRequest) (*GetPodResponse, error)
	mustEmbedUnimplementedKubeletApiServerServiceServer()
}

// UnimplementedKubeletApiServerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKubeletApiServerServiceServer struct {
}

func (UnimplementedKubeletApiServerServiceServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedKubeletApiServerServiceServer) CreatePod(context.Context, *ApplyPodRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePod not implemented")
}
func (UnimplementedKubeletApiServerServiceServer) DeletePod(context.Context, *DeletePodRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePod not implemented")
}
func (UnimplementedKubeletApiServerServiceServer) GetPod(context.Context, *GetPodRequest) (*GetPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPod not implemented")
}
func (UnimplementedKubeletApiServerServiceServer) mustEmbedUnimplementedKubeletApiServerServiceServer() {
}

// UnsafeKubeletApiServerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KubeletApiServerServiceServer will
// result in compilation errors.
type UnsafeKubeletApiServerServiceServer interface {
	mustEmbedUnimplementedKubeletApiServerServiceServer()
}

func RegisterKubeletApiServerServiceServer(s grpc.ServiceRegistrar, srv KubeletApiServerServiceServer) {
	s.RegisterService(&KubeletApiServerService_ServiceDesc, srv)
}

func _KubeletApiServerService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KubeletApiServerServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KubeletApiServerService_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KubeletApiServerServiceServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KubeletApiServerService_CreatePod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KubeletApiServerServiceServer).CreatePod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KubeletApiServerService_CreatePod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KubeletApiServerServiceServer).CreatePod(ctx, req.(*ApplyPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KubeletApiServerService_DeletePod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KubeletApiServerServiceServer).DeletePod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KubeletApiServerService_DeletePod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KubeletApiServerServiceServer).DeletePod(ctx, req.(*DeletePodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KubeletApiServerService_GetPod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KubeletApiServerServiceServer).GetPod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KubeletApiServerService_GetPod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KubeletApiServerServiceServer).GetPod(ctx, req.(*GetPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KubeletApiServerService_ServiceDesc is the grpc.ServiceDesc for KubeletApiServerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KubeletApiServerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kubelet_for_apiserver.kubeletApiServerService",
	HandlerType: (*KubeletApiServerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _KubeletApiServerService_SayHello_Handler,
		},
		{
			MethodName: "CreatePod",
			Handler:    _KubeletApiServerService_CreatePod_Handler,
		},
		{
			MethodName: "DeletePod",
			Handler:    _KubeletApiServerService_DeletePod_Handler,
		},
		{
			MethodName: "GetPod",
			Handler:    _KubeletApiServerService_GetPod_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/kubelet_for_apiserver.proto",
}
