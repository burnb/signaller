// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: signaller_api.proto

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

// SignallerApiClient is the client API for SignallerApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SignallerApiClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (SignallerApi_StreamClient, error)
}

type signallerApiClient struct {
	cc grpc.ClientConnInterface
}

func NewSignallerApiClient(cc grpc.ClientConnInterface) SignallerApiClient {
	return &signallerApiClient{cc}
}

func (c *signallerApiClient) Stream(ctx context.Context, opts ...grpc.CallOption) (SignallerApi_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &SignallerApi_ServiceDesc.Streams[0], "/signaller.SignallerApi/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &signallerApiStreamClient{stream}
	return x, nil
}

type SignallerApi_StreamClient interface {
	Send(*SubscriptionRequest) error
	Recv() (*PositionEvent, error)
	grpc.ClientStream
}

type signallerApiStreamClient struct {
	grpc.ClientStream
}

func (x *signallerApiStreamClient) Send(m *SubscriptionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *signallerApiStreamClient) Recv() (*PositionEvent, error) {
	m := new(PositionEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SignallerApiServer is the server API for SignallerApi service.
// All implementations must embed UnimplementedSignallerApiServer
// for forward compatibility
type SignallerApiServer interface {
	Stream(SignallerApi_StreamServer) error
	mustEmbedUnimplementedSignallerApiServer()
}

// UnimplementedSignallerApiServer must be embedded to have forward compatible implementations.
type UnimplementedSignallerApiServer struct {
}

func (UnimplementedSignallerApiServer) Stream(SignallerApi_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedSignallerApiServer) mustEmbedUnimplementedSignallerApiServer() {}

// UnsafeSignallerApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SignallerApiServer will
// result in compilation errors.
type UnsafeSignallerApiServer interface {
	mustEmbedUnimplementedSignallerApiServer()
}

func RegisterSignallerApiServer(s grpc.ServiceRegistrar, srv SignallerApiServer) {
	s.RegisterService(&SignallerApi_ServiceDesc, srv)
}

func _SignallerApi_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SignallerApiServer).Stream(&signallerApiStreamServer{stream})
}

type SignallerApi_StreamServer interface {
	Send(*PositionEvent) error
	Recv() (*SubscriptionRequest, error)
	grpc.ServerStream
}

type signallerApiStreamServer struct {
	grpc.ServerStream
}

func (x *signallerApiStreamServer) Send(m *PositionEvent) error {
	return x.ServerStream.SendMsg(m)
}

func (x *signallerApiStreamServer) Recv() (*SubscriptionRequest, error) {
	m := new(SubscriptionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SignallerApi_ServiceDesc is the grpc.ServiceDesc for SignallerApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SignallerApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "signaller.SignallerApi",
	HandlerType: (*SignallerApiServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _SignallerApi_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "signaller_api.proto",
}
