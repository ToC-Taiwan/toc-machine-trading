// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: forwarder/realtime.proto

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
	RealTimeDataInterface_GetAllStockSnapshot_FullMethodName        = "/forwarder.RealTimeDataInterface/GetAllStockSnapshot"
	RealTimeDataInterface_GetStockSnapshotByNumArr_FullMethodName   = "/forwarder.RealTimeDataInterface/GetStockSnapshotByNumArr"
	RealTimeDataInterface_GetStockSnapshotTSE_FullMethodName        = "/forwarder.RealTimeDataInterface/GetStockSnapshotTSE"
	RealTimeDataInterface_GetStockSnapshotOTC_FullMethodName        = "/forwarder.RealTimeDataInterface/GetStockSnapshotOTC"
	RealTimeDataInterface_GetNasdaq_FullMethodName                  = "/forwarder.RealTimeDataInterface/GetNasdaq"
	RealTimeDataInterface_GetNasdaqFuture_FullMethodName            = "/forwarder.RealTimeDataInterface/GetNasdaqFuture"
	RealTimeDataInterface_GetStockVolumeRank_FullMethodName         = "/forwarder.RealTimeDataInterface/GetStockVolumeRank"
	RealTimeDataInterface_GetFutureSnapshotByCodeArr_FullMethodName = "/forwarder.RealTimeDataInterface/GetFutureSnapshotByCodeArr"
)

// RealTimeDataInterfaceClient is the client API for RealTimeDataInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RealTimeDataInterfaceClient interface {
	GetAllStockSnapshot(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error)
	GetStockSnapshotByNumArr(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SnapshotResponse, error)
	GetStockSnapshotTSE(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error)
	GetStockSnapshotOTC(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error)
	GetNasdaq(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*YahooFinancePrice, error)
	GetNasdaqFuture(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*YahooFinancePrice, error)
	GetStockVolumeRank(ctx context.Context, in *VolumeRankRequest, opts ...grpc.CallOption) (*StockVolumeRankResponse, error)
	GetFutureSnapshotByCodeArr(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SnapshotResponse, error)
}

type realTimeDataInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewRealTimeDataInterfaceClient(cc grpc.ClientConnInterface) RealTimeDataInterfaceClient {
	return &realTimeDataInterfaceClient{cc}
}

func (c *realTimeDataInterfaceClient) GetAllStockSnapshot(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetAllStockSnapshot_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetStockSnapshotByNumArr(ctx context.Context, in *StockNumArr, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetStockSnapshotByNumArr_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetStockSnapshotTSE(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetStockSnapshotTSE_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetStockSnapshotOTC(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetStockSnapshotOTC_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetNasdaq(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*YahooFinancePrice, error) {
	out := new(YahooFinancePrice)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetNasdaq_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetNasdaqFuture(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*YahooFinancePrice, error) {
	out := new(YahooFinancePrice)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetNasdaqFuture_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetStockVolumeRank(ctx context.Context, in *VolumeRankRequest, opts ...grpc.CallOption) (*StockVolumeRankResponse, error) {
	out := new(StockVolumeRankResponse)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetStockVolumeRank_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realTimeDataInterfaceClient) GetFutureSnapshotByCodeArr(ctx context.Context, in *FutureCodeArr, opts ...grpc.CallOption) (*SnapshotResponse, error) {
	out := new(SnapshotResponse)
	err := c.cc.Invoke(ctx, RealTimeDataInterface_GetFutureSnapshotByCodeArr_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RealTimeDataInterfaceServer is the server API for RealTimeDataInterface service.
// All implementations must embed UnimplementedRealTimeDataInterfaceServer
// for forward compatibility
type RealTimeDataInterfaceServer interface {
	GetAllStockSnapshot(context.Context, *emptypb.Empty) (*SnapshotResponse, error)
	GetStockSnapshotByNumArr(context.Context, *StockNumArr) (*SnapshotResponse, error)
	GetStockSnapshotTSE(context.Context, *emptypb.Empty) (*SnapshotResponse, error)
	GetStockSnapshotOTC(context.Context, *emptypb.Empty) (*SnapshotResponse, error)
	GetNasdaq(context.Context, *emptypb.Empty) (*YahooFinancePrice, error)
	GetNasdaqFuture(context.Context, *emptypb.Empty) (*YahooFinancePrice, error)
	GetStockVolumeRank(context.Context, *VolumeRankRequest) (*StockVolumeRankResponse, error)
	GetFutureSnapshotByCodeArr(context.Context, *FutureCodeArr) (*SnapshotResponse, error)
	mustEmbedUnimplementedRealTimeDataInterfaceServer()
}

// UnimplementedRealTimeDataInterfaceServer must be embedded to have forward compatible implementations.
type UnimplementedRealTimeDataInterfaceServer struct {
}

func (UnimplementedRealTimeDataInterfaceServer) GetAllStockSnapshot(context.Context, *emptypb.Empty) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllStockSnapshot not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetStockSnapshotByNumArr(context.Context, *StockNumArr) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockSnapshotByNumArr not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetStockSnapshotTSE(context.Context, *emptypb.Empty) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockSnapshotTSE not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetStockSnapshotOTC(context.Context, *emptypb.Empty) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockSnapshotOTC not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetNasdaq(context.Context, *emptypb.Empty) (*YahooFinancePrice, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNasdaq not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetNasdaqFuture(context.Context, *emptypb.Empty) (*YahooFinancePrice, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNasdaqFuture not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetStockVolumeRank(context.Context, *VolumeRankRequest) (*StockVolumeRankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockVolumeRank not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) GetFutureSnapshotByCodeArr(context.Context, *FutureCodeArr) (*SnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFutureSnapshotByCodeArr not implemented")
}
func (UnimplementedRealTimeDataInterfaceServer) mustEmbedUnimplementedRealTimeDataInterfaceServer() {}

// UnsafeRealTimeDataInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RealTimeDataInterfaceServer will
// result in compilation errors.
type UnsafeRealTimeDataInterfaceServer interface {
	mustEmbedUnimplementedRealTimeDataInterfaceServer()
}

func RegisterRealTimeDataInterfaceServer(s grpc.ServiceRegistrar, srv RealTimeDataInterfaceServer) {
	s.RegisterService(&RealTimeDataInterface_ServiceDesc, srv)
}

func _RealTimeDataInterface_GetAllStockSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetAllStockSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetAllStockSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetAllStockSnapshot(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetStockSnapshotByNumArr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetStockSnapshotByNumArr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetStockSnapshotByNumArr_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetStockSnapshotByNumArr(ctx, req.(*StockNumArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetStockSnapshotTSE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetStockSnapshotTSE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetStockSnapshotTSE_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetStockSnapshotTSE(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetStockSnapshotOTC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetStockSnapshotOTC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetStockSnapshotOTC_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetStockSnapshotOTC(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetNasdaq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetNasdaq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetNasdaq_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetNasdaq(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetNasdaqFuture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetNasdaqFuture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetNasdaqFuture_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetNasdaqFuture(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetStockVolumeRank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VolumeRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetStockVolumeRank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetStockVolumeRank_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetStockVolumeRank(ctx, req.(*VolumeRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealTimeDataInterface_GetFutureSnapshotByCodeArr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealTimeDataInterfaceServer).GetFutureSnapshotByCodeArr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealTimeDataInterface_GetFutureSnapshotByCodeArr_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealTimeDataInterfaceServer).GetFutureSnapshotByCodeArr(ctx, req.(*FutureCodeArr))
	}
	return interceptor(ctx, in, info, handler)
}

// RealTimeDataInterface_ServiceDesc is the grpc.ServiceDesc for RealTimeDataInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RealTimeDataInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "forwarder.RealTimeDataInterface",
	HandlerType: (*RealTimeDataInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllStockSnapshot",
			Handler:    _RealTimeDataInterface_GetAllStockSnapshot_Handler,
		},
		{
			MethodName: "GetStockSnapshotByNumArr",
			Handler:    _RealTimeDataInterface_GetStockSnapshotByNumArr_Handler,
		},
		{
			MethodName: "GetStockSnapshotTSE",
			Handler:    _RealTimeDataInterface_GetStockSnapshotTSE_Handler,
		},
		{
			MethodName: "GetStockSnapshotOTC",
			Handler:    _RealTimeDataInterface_GetStockSnapshotOTC_Handler,
		},
		{
			MethodName: "GetNasdaq",
			Handler:    _RealTimeDataInterface_GetNasdaq_Handler,
		},
		{
			MethodName: "GetNasdaqFuture",
			Handler:    _RealTimeDataInterface_GetNasdaqFuture_Handler,
		},
		{
			MethodName: "GetStockVolumeRank",
			Handler:    _RealTimeDataInterface_GetStockVolumeRank_Handler,
		},
		{
			MethodName: "GetFutureSnapshotByCodeArr",
			Handler:    _RealTimeDataInterface_GetFutureSnapshotByCodeArr_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "forwarder/realtime.proto",
}
