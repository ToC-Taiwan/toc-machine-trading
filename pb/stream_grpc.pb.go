// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: stream.proto

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

// StreamDataInterfaceClient is the client API for StreamDataInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamDataInterfaceClient interface {
	GetStockSnapshotByNumArr(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SnapshotResponse, error)
	GetAllStockSnapshot(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error)
	GetStockSnapshotTSE(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotMessage, error)
	GetStockVolumeRank(ctx context.Context, in *VolumeRankRequest, opts ...grpc.CallOption) (*StockVolumeRankResponse, error)
	SubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	UnSubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	SubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	UnSubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	SubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	UnSubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error)
	UnSubscribeStockAllTick(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error)
	UnSubscribeStockAllBidAsk(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error)
	GetFutureSnapshotByCodeArr(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SnapshotResponse, error)
}

type streamDataInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamDataInterfaceClient(cc grpc.ClientConnInterface) StreamDataInterfaceClient {
	return &streamDataInterfaceClient{cc}
}

func (c *streamDataInterfaceClient) GetStockSnapshotByNumArr(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/GetStockSnapshotByNumArr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) GetAllStockSnapshot(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/GetAllStockSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) GetStockSnapshotTSE(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotMessage, error) {
	out := new(SnapshotMessage)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/GetStockSnapshotTSE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) GetStockVolumeRank(ctx context.Context, in *VolumeRankRequest, opts ...grpc.CallOption) (*StockVolumeRankResponse, error) {
	out := new(StockVolumeRankResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/GetStockVolumeRank", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) SubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/SubscribeStockTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) UnSubscribeStockTick(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) SubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/SubscribeStockBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) UnSubscribeStockBidAsk(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) SubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/SubscribeFutureTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) UnSubscribeFutureTick(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/UnSubscribeFutureTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) UnSubscribeStockAllTick(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error) {
	out := new(ErrorMessage)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockAllTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) UnSubscribeStockAllBidAsk(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ErrorMessage, error) {
	out := new(ErrorMessage)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockAllBidAsk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamDataInterfaceClient) GetFutureSnapshotByCodeArr(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.StreamDataInterface/GetFutureSnapshotByCodeArr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreamDataInterfaceServer is the server API for StreamDataInterface service.
// All implementations must embed UnimplementedStreamDataInterfaceServer
// for forward compatibility
type StreamDataInterfaceServer interface {
	GetStockSnapshotByNumArr(context.Context, *StockNumArr) (*SnapshotResponse, error)
	GetAllStockSnapshot(context.Context, *emptypb.Empty) (*SnapshotResponse, error)
	GetStockSnapshotTSE(context.Context, *emptypb.Empty) (*SnapshotMessage, error)
	GetStockVolumeRank(context.Context, *VolumeRankRequest) (*StockVolumeRankResponse, error)
	SubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error)
	UnSubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error)
	SubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error)
	UnSubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error)
	SubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error)
	UnSubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error)
	UnSubscribeStockAllTick(context.Context, *emptypb.Empty) (*ErrorMessage, error)
	UnSubscribeStockAllBidAsk(context.Context, *emptypb.Empty) (*ErrorMessage, error)
	GetFutureSnapshotByCodeArr(context.Context, *FutureCodeArr) (*SnapshotResponse, error)
	mustEmbedUnimplementedStreamDataInterfaceServer()
}

// UnimplementedStreamDataInterfaceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamDataInterfaceServer struct {
}

func (UnimplementedStreamDataInterfaceServer) GetStockSnapshotByNumArr(context.Context, *StockNumArr) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockSnapshotByNumArr not implemented")
}
func (UnimplementedStreamDataInterfaceServer) GetAllStockSnapshot(context.Context, *emptypb.Empty) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllStockSnapshot not implemented")
}
func (UnimplementedStreamDataInterfaceServer) GetStockSnapshotTSE(context.Context, *emptypb.Empty) (*SnapshotMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockSnapshotTSE not implemented")
}
func (UnimplementedStreamDataInterfaceServer) GetStockVolumeRank(context.Context, *VolumeRankRequest) (*StockVolumeRankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockVolumeRank not implemented")
}
func (UnimplementedStreamDataInterfaceServer) SubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeStockTick not implemented")
}
func (UnimplementedStreamDataInterfaceServer) UnSubscribeStockTick(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeStockTick not implemented")
}
func (UnimplementedStreamDataInterfaceServer) SubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeStockBidAsk not implemented")
}
func (UnimplementedStreamDataInterfaceServer) UnSubscribeStockBidAsk(context.Context, *StockNumArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeStockBidAsk not implemented")
}
func (UnimplementedStreamDataInterfaceServer) SubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeFutureTick not implemented")
}
func (UnimplementedStreamDataInterfaceServer) UnSubscribeFutureTick(context.Context, *FutureCodeArr) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeFutureTick not implemented")
}
func (UnimplementedStreamDataInterfaceServer) UnSubscribeStockAllTick(context.Context, *emptypb.Empty) (*ErrorMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeStockAllTick not implemented")
}
func (UnimplementedStreamDataInterfaceServer) UnSubscribeStockAllBidAsk(context.Context, *emptypb.Empty) (*ErrorMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnSubscribeStockAllBidAsk not implemented")
}
func (UnimplementedStreamDataInterfaceServer) GetFutureSnapshotByCodeArr(context.Context, *FutureCodeArr) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFutureSnapshotByCodeArr not implemented")
}
func (UnimplementedStreamDataInterfaceServer) mustEmbedUnimplementedStreamDataInterfaceServer() {}

// UnsafeStreamDataInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamDataInterfaceServer will
// result in compilation errors.
type UnsafeStreamDataInterfaceServer interface {
	mustEmbedUnimplementedStreamDataInterfaceServer()
}

func RegisterStreamDataInterfaceServer(s grpc.ServiceRegistrar, srv StreamDataInterfaceServer) {
	s.RegisterService(&StreamDataInterface_ServiceDesc, srv)
}

func _StreamDataInterface_GetStockSnapshotByNumArr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).GetStockSnapshotByNumArr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/GetStockSnapshotByNumArr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).GetStockSnapshotByNumArr(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_GetAllStockSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).GetAllStockSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/GetAllStockSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).GetAllStockSnapshot(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_GetStockSnapshotTSE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).GetStockSnapshotTSE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/GetStockSnapshotTSE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).GetStockSnapshotTSE(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_GetStockVolumeRank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VolumeRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).GetStockVolumeRank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/GetStockVolumeRank",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).GetStockVolumeRank(ctx, req.(*VolumeRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_SubscribeStockTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).SubscribeStockTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/SubscribeStockTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).SubscribeStockTick(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_UnSubscribeStockTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockTick(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_SubscribeStockBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).SubscribeStockBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/SubscribeStockBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).SubscribeStockBidAsk(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_UnSubscribeStockBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockBidAsk(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_SubscribeFutureTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).SubscribeFutureTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/SubscribeFutureTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).SubscribeFutureTick(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_UnSubscribeFutureTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).UnSubscribeFutureTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/UnSubscribeFutureTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).UnSubscribeFutureTick(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_UnSubscribeStockAllTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockAllTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockAllTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockAllTick(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_UnSubscribeStockAllBidAsk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockAllBidAsk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/UnSubscribeStockAllBidAsk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).UnSubscribeStockAllBidAsk(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamDataInterface_GetFutureSnapshotByCodeArr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamDataInterfaceServer).GetFutureSnapshotByCodeArr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.StreamDataInterface/GetFutureSnapshotByCodeArr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamDataInterfaceServer).GetFutureSnapshotByCodeArr(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

// StreamDataInterface_ServiceDesc is the grpc.ServiceDesc for StreamDataInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamDataInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sinopac_forwarder.StreamDataInterface",
	HandlerType: (*StreamDataInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStockSnapshotByNumArr",
			Handler:    _StreamDataInterface_GetStockSnapshotByNumArr_Handler,
		},
		{
			MethodName: "GetAllStockSnapshot",
			Handler:    _StreamDataInterface_GetAllStockSnapshot_Handler,
		},
		{
			MethodName: "GetStockSnapshotTSE",
			Handler:    _StreamDataInterface_GetStockSnapshotTSE_Handler,
		},
		{
			MethodName: "GetStockVolumeRank",
			Handler:    _StreamDataInterface_GetStockVolumeRank_Handler,
		},
		{
			MethodName: "SubscribeStockTick",
			Handler:    _StreamDataInterface_SubscribeStockTick_Handler,
		},
		{
			MethodName: "UnSubscribeStockTick",
			Handler:    _StreamDataInterface_UnSubscribeStockTick_Handler,
		},
		{
			MethodName: "SubscribeStockBidAsk",
			Handler:    _StreamDataInterface_SubscribeStockBidAsk_Handler,
		},
		{
			MethodName: "UnSubscribeStockBidAsk",
			Handler:    _StreamDataInterface_UnSubscribeStockBidAsk_Handler,
		},
		{
			MethodName: "SubscribeFutureTick",
			Handler:    _StreamDataInterface_SubscribeFutureTick_Handler,
		},
		{
			MethodName: "UnSubscribeFutureTick",
			Handler:    _StreamDataInterface_UnSubscribeFutureTick_Handler,
		},
		{
			MethodName: "UnSubscribeStockAllTick",
			Handler:    _StreamDataInterface_UnSubscribeStockAllTick_Handler,
		},
		{
			MethodName: "UnSubscribeStockAllBidAsk",
			Handler:    _StreamDataInterface_UnSubscribeStockAllBidAsk_Handler,
		},
		{
			MethodName: "GetFutureSnapshotByCodeArr",
			Handler:    _StreamDataInterface_GetFutureSnapshotByCodeArr_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stream.proto",
}
