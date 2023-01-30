// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: subscribe.proto

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

// SubscribeDataInterfaceClient is the client API for SubscribeDataInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubscribeDataInterfaceClient interface {
	// SubscribeStockTick is the interface for subscribe stock tick
	SubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// UnSubscribeStockTick is the interface for unsubscribe stock tick
	UnSubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// SubscribeStockBidAsk is the interface for subscribe stock bid ask
	SubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// UnSubscribeStockBidAsk is the interface for unsubscribe stock bid ask
	UnSubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// SubscribeFutureTick is the interface for subscribe stock all tick
	SubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// UnSubscribeFutureTick is the interface for unsubscribe stock all tick
	UnSubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// SubscribeFutureBidAsk is the interface for subscribe stock all bid ask
	SubscribeFutureBidAsk(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// UnSubscribeFutureBidAsk is the interface for unsubscribe stock all bid ask
	UnSubscribeFutureBidAsk(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	// UnSubscribeAllTick is the interface for unsubscribe stock all tick
	UnSubscribeAllTick(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error)
	// UnSubscribeStockAllBidAsk is the interface for unsubscribe stock all bid ask
	UnSubscribeAllBidAsk(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error)
}

type subscribeDataInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewSubscribeDataInterfaceClient(cc grpc.ClientConnInterface) SubscribeDataInterfaceClient {
	return &subscribeDataInterfaceClient{cc}
}

func (c *subscribeDataInterfaceClient) SubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/SubscribeStockTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) UnSubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeStockTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) SubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/SubscribeStockBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) UnSubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeStockBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) SubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/SubscribeFutureTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) UnSubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeFutureTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) SubscribeFutureBidAsk(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/SubscribeFutureBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) UnSubscribeFutureBidAsk(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeFutureBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) UnSubscribeAllTick(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error) {
	out := new(ErrorMessage)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeAllTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscribeDataInterfaceClient) UnSubscribeAllBidAsk(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error) {
	out := new(ErrorMessage)
	err := c.cc.Invoke(ctx, "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeAllBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubscribeDataInterfaceServer is the server API for SubscribeDataInterface service.
// All implementations must embed UnimplementedSubscribeDataInterfaceServer
// for forward compatibility
type SubscribeDataInterfaceServer interface {
	// SubscribeStockTick is the interface for subscribe stock tick
	SubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error)
	// UnSubscribeStockTick is the interface for unsubscribe stock tick
	UnSubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error)
	// SubscribeStockBidAsk is the interface for subscribe stock bid ask
	SubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error)
	// UnSubscribeStockBidAsk is the interface for unsubscribe stock bid ask
	UnSubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error)
	// SubscribeFutureTick is the interface for subscribe stock all tick
	SubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error)
	// UnSubscribeFutureTick is the interface for unsubscribe stock all tick
	UnSubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error)
	// SubscribeFutureBidAsk is the interface for subscribe stock all bid ask
	SubscribeFutureBidAsk(context.Context, *FutureCodeArr) (*SubscribeResponse, error)
	// UnSubscribeFutureBidAsk is the interface for unsubscribe stock all bid ask
	UnSubscribeFutureBidAsk(context.Context, *FutureCodeArr) (*SubscribeResponse, error)
	// UnSubscribeAllTick is the interface for unsubscribe stock all tick
	UnSubscribeAllTick(context.Context, *emptypb.Empty) (*ErrorMessage, error)
	// UnSubscribeStockAllBidAsk is the interface for unsubscribe stock all bid ask
	UnSubscribeAllBidAsk(context.Context, *emptypb.Empty) (*ErrorMessage, error)
	mustEmbedUnimplementedSubscribeDataInterfaceServer()
}

// UnimplementedSubscribeDataInterfaceServer must be embedded to have forward compatible implementations.
type UnimplementedSubscribeDataInterfaceServer struct {
}

func (UnimplementedSubscribeDataInterfaceServer) SubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeStockTick not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) UnSubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeStockTick not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) SubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeStockBidAsk not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) UnSubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeStockBidAsk not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) SubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeFutureTick not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) UnSubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeFutureTick not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) SubscribeFutureBidAsk(context.Context, *FutureCodeArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeFutureBidAsk not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) UnSubscribeFutureBidAsk(context.Context, *FutureCodeArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeFutureBidAsk not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) UnSubscribeAllTick(context.Context, *emptypb.Empty) (*ErrorMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeAllTick not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) UnSubscribeAllBidAsk(context.Context, *emptypb.Empty) (*ErrorMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeAllBidAsk not implemented")
}
func (UnimplementedSubscribeDataInterfaceServer) mustEmbedUnimplementedSubscribeDataInterfaceServer() {
}

// UnsafeSubscribeDataInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubscribeDataInterfaceServer will
// result in compilation errors.
type UnsafeSubscribeDataInterfaceServer interface {
	mustEmbedUnimplementedSubscribeDataInterfaceServer()
}

func RegisterSubscribeDataInterfaceServer(s grpc.ServiceRegistrar, srv SubscribeDataInterfaceServer) {
	s.RegisterService(&SubscribeDataInterface_ServiceDesc, srv)
}

func _SubscribeDataInterface_SubscribeStockTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).SubscribeStockTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/SubscribeStockTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).SubscribeStockTick(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_UnSubscribeStockTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeStockTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeStockTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeStockTick(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_SubscribeStockBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).SubscribeStockBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/SubscribeStockBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).SubscribeStockBidAsk(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_UnSubscribeStockBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeStockBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeStockBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeStockBidAsk(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_SubscribeFutureTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).SubscribeFutureTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/SubscribeFutureTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).SubscribeFutureTick(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_UnSubscribeFutureTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeFutureTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeFutureTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeFutureTick(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_SubscribeFutureBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).SubscribeFutureBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/SubscribeFutureBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).SubscribeFutureBidAsk(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_UnSubscribeFutureBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeFutureBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeFutureBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeFutureBidAsk(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_UnSubscribeAllTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeAllTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeAllTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeAllTick(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SubscribeDataInterface_UnSubscribeAllBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeAllBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/toc_python_forwarder.SubscribeDataInterface/UnSubscribeAllBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscribeDataInterfaceServer).UnSubscribeAllBidAsk(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SubscribeDataInterface_ServiceDesc is the grpc.ServiceDesc for SubscribeDataInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SubscribeDataInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "toc_python_forwarder.SubscribeDataInterface",
	HandlerType: (*SubscribeDataInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubscribeStockTick",
			Handler:    _SubscribeDataInterface_SubscribeStockTick_Handler,
		},
		{
			MethodName: "UnSubscribeStockTick",
			Handler:    _SubscribeDataInterface_UnSubscribeStockTick_Handler,
		},
		{
			MethodName: "SubscribeStockBidAsk",
			Handler:    _SubscribeDataInterface_SubscribeStockBidAsk_Handler,
		},
		{
			MethodName: "UnSubscribeStockBidAsk",
			Handler:    _SubscribeDataInterface_UnSubscribeStockBidAsk_Handler,
		},
		{
			MethodName: "SubscribeFutureTick",
			Handler:    _SubscribeDataInterface_SubscribeFutureTick_Handler,
		},
		{
			MethodName: "UnSubscribeFutureTick",
			Handler:    _SubscribeDataInterface_UnSubscribeFutureTick_Handler,
		},
		{
			MethodName: "SubscribeFutureBidAsk",
			Handler:    _SubscribeDataInterface_SubscribeFutureBidAsk_Handler,
		},
		{
			MethodName: "UnSubscribeFutureBidAsk",
			Handler:    _SubscribeDataInterface_UnSubscribeFutureBidAsk_Handler,
		},
		{
			MethodName: "UnSubscribeAllTick",
			Handler:    _SubscribeDataInterface_UnSubscribeAllTick_Handler,
		},
		{
			MethodName: "UnSubscribeAllBidAsk",
			Handler:    _SubscribeDataInterface_UnSubscribeAllBidAsk_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "subscribe.proto",
}