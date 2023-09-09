// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: basic.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BasicDataInterface_CreateLongConnection_FullMethodName = "/toc_python_forwarder.BasicDataInterface/CreateLongConnection"
	BasicDataInterface_Terminate_FullMethodName            = "/toc_python_forwarder.BasicDataInterface/Terminate"
	BasicDataInterface_CheckUsage_FullMethodName           = "/toc_python_forwarder.BasicDataInterface/CheckUsage"
	BasicDataInterface_Login_FullMethodName                = "/toc_python_forwarder.BasicDataInterface/Login"
	BasicDataInterface_GetAllStockDetail_FullMethodName    = "/toc_python_forwarder.BasicDataInterface/GetAllStockDetail"
	BasicDataInterface_GetAllFutureDetail_FullMethodName   = "/toc_python_forwarder.BasicDataInterface/GetAllFutureDetail"
	BasicDataInterface_GetAllOptionDetail_FullMethodName   = "/toc_python_forwarder.BasicDataInterface/GetAllOptionDetail"
)

// BasicDataInterfaceClient is the client API for BasicDataInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BasicDataInterfaceClient interface {
	// CreateLongConnection is the function to create long connection
	CreateLongConnection(ctx context.Context, opts ...grpc.CallOption) (BasicDataInterface_CreateLongConnectionClient, error)
	// Terminate is the terminate function
	Terminate(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// CheckUsage get shioaji usage
	CheckUsage(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ShioajiUsage, error)
	// Login log in
	Login(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// GetAllStockDetail is the function to get stock detail
	GetAllStockDetail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StockDetailResponse, error)
	// GetAllFutureDetail is the function to get future detail
	GetAllFutureDetail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FutureDetailResponse, error)
	// GetAllOptionDetail is the function to get option detail
	GetAllOptionDetail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OptionDetailResponse, error)
}

type basicDataInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewBasicDataInterfaceClient(cc grpc.ClientConnInterface) BasicDataInterfaceClient {
	return &basicDataInterfaceClient{cc}
}

func (c *basicDataInterfaceClient) CreateLongConnection(ctx context.Context, opts ...grpc.CallOption) (BasicDataInterface_CreateLongConnectionClient, error) {
	stream, err := c.cc.NewStream(ctx, &BasicDataInterface_ServiceDesc.Streams[0], BasicDataInterface_CreateLongConnection_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &basicDataInterfaceCreateLongConnectionClient{stream}
	return x, nil
}

type BasicDataInterface_CreateLongConnectionClient interface {
	Send(*emptypb.Empty) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type basicDataInterfaceCreateLongConnectionClient struct {
	grpc.ClientStream
}

func (x *basicDataInterfaceCreateLongConnectionClient) Send(m *emptypb.Empty) error {
	return x.ClientStream.SendMsg(m)
}

func (x *basicDataInterfaceCreateLongConnectionClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *basicDataInterfaceClient) Terminate(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BasicDataInterface_Terminate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicDataInterfaceClient) CheckUsage(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ShioajiUsage, error) {
	out := new(ShioajiUsage)
	err := c.cc.Invoke(ctx, BasicDataInterface_CheckUsage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicDataInterfaceClient) Login(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BasicDataInterface_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicDataInterfaceClient) GetAllStockDetail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StockDetailResponse, error) {
	out := new(StockDetailResponse)
	err := c.cc.Invoke(ctx, BasicDataInterface_GetAllStockDetail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicDataInterfaceClient) GetAllFutureDetail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FutureDetailResponse, error) {
	out := new(FutureDetailResponse)
	err := c.cc.Invoke(ctx, BasicDataInterface_GetAllFutureDetail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicDataInterfaceClient) GetAllOptionDetail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OptionDetailResponse, error) {
	out := new(OptionDetailResponse)
	err := c.cc.Invoke(ctx, BasicDataInterface_GetAllOptionDetail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BasicDataInterfaceServer is the server API for BasicDataInterface service.
// All implementations must embed UnimplementedBasicDataInterfaceServer
// for forward compatibility
type BasicDataInterfaceServer interface {
	// CreateLongConnection is the function to create long connection
	CreateLongConnection(BasicDataInterface_CreateLongConnectionServer) error
	// Terminate is the terminate function
	Terminate(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// CheckUsage get shioaji usage
	CheckUsage(context.Context, *emptypb.Empty) (*ShioajiUsage, error)
	// Login log in
	Login(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// GetAllStockDetail is the function to get stock detail
	GetAllStockDetail(context.Context, *emptypb.Empty) (*StockDetailResponse, error)
	// GetAllFutureDetail is the function to get future detail
	GetAllFutureDetail(context.Context, *emptypb.Empty) (*FutureDetailResponse, error)
	// GetAllOptionDetail is the function to get option detail
	GetAllOptionDetail(context.Context, *emptypb.Empty) (*OptionDetailResponse, error)
	mustEmbedUnimplementedBasicDataInterfaceServer()
}

// UnimplementedBasicDataInterfaceServer must be embedded to have forward compatible implementations.
type UnimplementedBasicDataInterfaceServer struct {
}

func (UnimplementedBasicDataInterfaceServer) CreateLongConnection(BasicDataInterface_CreateLongConnectionServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateLongConnection not implemented")
}
func (UnimplementedBasicDataInterfaceServer) Terminate(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Terminate not implemented")
}
func (UnimplementedBasicDataInterfaceServer) CheckUsage(context.Context, *emptypb.Empty) (*ShioajiUsage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUsage not implemented")
}
func (UnimplementedBasicDataInterfaceServer) Login(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedBasicDataInterfaceServer) GetAllStockDetail(context.Context, *emptypb.Empty) (*StockDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllStockDetail not implemented")
}
func (UnimplementedBasicDataInterfaceServer) GetAllFutureDetail(context.Context, *emptypb.Empty) (*FutureDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllFutureDetail not implemented")
}
func (UnimplementedBasicDataInterfaceServer) GetAllOptionDetail(context.Context, *emptypb.Empty) (*OptionDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllOptionDetail not implemented")
}
func (UnimplementedBasicDataInterfaceServer) mustEmbedUnimplementedBasicDataInterfaceServer() {}

// UnsafeBasicDataInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BasicDataInterfaceServer will
// result in compilation errors.
type UnsafeBasicDataInterfaceServer interface {
	mustEmbedUnimplementedBasicDataInterfaceServer()
}

func RegisterBasicDataInterfaceServer(s grpc.ServiceRegistrar, srv BasicDataInterfaceServer) {
	s.RegisterService(&BasicDataInterface_ServiceDesc, srv)
}

func _BasicDataInterface_CreateLongConnection_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BasicDataInterfaceServer).CreateLongConnection(&basicDataInterfaceCreateLongConnectionServer{stream})
}

type BasicDataInterface_CreateLongConnectionServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*emptypb.Empty, error)
	grpc.ServerStream
}

type basicDataInterfaceCreateLongConnectionServer struct {
	grpc.ServerStream
}

func (x *basicDataInterfaceCreateLongConnectionServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *basicDataInterfaceCreateLongConnectionServer) Recv() (*emptypb.Empty, error) {
	m := new(emptypb.Empty)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _BasicDataInterface_Terminate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicDataInterfaceServer).Terminate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasicDataInterface_Terminate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicDataInterfaceServer).Terminate(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicDataInterface_CheckUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicDataInterfaceServer).CheckUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasicDataInterface_CheckUsage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicDataInterfaceServer).CheckUsage(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicDataInterface_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicDataInterfaceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasicDataInterface_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicDataInterfaceServer).Login(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicDataInterface_GetAllStockDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicDataInterfaceServer).GetAllStockDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasicDataInterface_GetAllStockDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicDataInterfaceServer).GetAllStockDetail(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicDataInterface_GetAllFutureDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicDataInterfaceServer).GetAllFutureDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasicDataInterface_GetAllFutureDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicDataInterfaceServer).GetAllFutureDetail(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicDataInterface_GetAllOptionDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicDataInterfaceServer).GetAllOptionDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasicDataInterface_GetAllOptionDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicDataInterfaceServer).GetAllOptionDetail(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BasicDataInterface_ServiceDesc is the grpc.ServiceDesc for BasicDataInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BasicDataInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "toc_python_forwarder.BasicDataInterface",
	HandlerType: (*BasicDataInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Terminate",
			Handler:    _BasicDataInterface_Terminate_Handler,
		},
		{
			MethodName: "CheckUsage",
			Handler:    _BasicDataInterface_CheckUsage_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _BasicDataInterface_Login_Handler,
		},
		{
			MethodName: "GetAllStockDetail",
			Handler:    _BasicDataInterface_GetAllStockDetail_Handler,
		},
		{
			MethodName: "GetAllFutureDetail",
			Handler:    _BasicDataInterface_GetAllFutureDetail_Handler,
		},
		{
			MethodName: "GetAllOptionDetail",
			Handler:    _BasicDataInterface_GetAllOptionDetail_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateLongConnection",
			Handler:       _BasicDataInterface_CreateLongConnection_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "basic.proto",
}
