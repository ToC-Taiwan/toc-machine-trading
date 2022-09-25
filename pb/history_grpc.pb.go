// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: history.proto

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

// HistoryDataInterfaceClient is the client API for HistoryDataInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HistoryDataInterfaceClient interface {
	GetStockHistoryTick(ctx context.Context, in *StockNumArrWithDate, opts ...grpc.CallOption) (*HistoryTickResponse, error)
	GetStockHistoryKbar(ctx context.Context, in *StockNumArrWithDate, opts ...grpc.CallOption) (*HistoryKbarResponse, error)
	GetStockHistoryClose(ctx context.Context, in *StockNumArrWithDate, opts ...grpc.CallOption) (*HistoryCloseResponse, error)
	GetStockHistoryCloseByDateArr(ctx context.Context, in *StockNumArrWithDateArr, opts ...grpc.CallOption) (*HistoryCloseResponse, error)
	GetStockTSEHistoryTick(ctx context.Context, in *Date, opts ...grpc.CallOption) (*HistoryTickResponse, error)
	GetStockTSEHistoryKbar(ctx context.Context, in *Date, opts ...grpc.CallOption) (*HistoryKbarResponse, error)
	GetStockTSEHistoryClose(ctx context.Context, in *Date, opts ...grpc.CallOption) (*HistoryCloseResponse, error)
	GetFutureHistoryTick(ctx context.Context, in *FutureCodeArrWithDate, opts ...grpc.CallOption) (*HistoryTickResponse, error)
	GetFutureHistoryClose(ctx context.Context, in *FutureCodeArrWithDate, opts ...grpc.CallOption) (*HistoryCloseResponse, error)
	GetFutureHistoryKbar(ctx context.Context, in *FutureCodeArrWithDate, opts ...grpc.CallOption) (*HistoryKbarResponse, error)
}

type historyDataInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewHistoryDataInterfaceClient(cc grpc.ClientConnInterface) HistoryDataInterfaceClient {
	return &historyDataInterfaceClient{cc}
}

func (c *historyDataInterfaceClient) GetStockHistoryTick(ctx context.Context, in *StockNumArrWithDate, opts ...grpc.CallOption) (*HistoryTickResponse, error) {
	out := new(HistoryTickResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetStockHistoryKbar(ctx context.Context, in *StockNumArrWithDate, opts ...grpc.CallOption) (*HistoryKbarResponse, error) {
	out := new(HistoryKbarResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryKbar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetStockHistoryClose(ctx context.Context, in *StockNumArrWithDate, opts ...grpc.CallOption) (*HistoryCloseResponse, error) {
	out := new(HistoryCloseResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryClose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetStockHistoryCloseByDateArr(ctx context.Context, in *StockNumArrWithDateArr, opts ...grpc.CallOption) (*HistoryCloseResponse, error) {
	out := new(HistoryCloseResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryCloseByDateArr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetStockTSEHistoryTick(ctx context.Context, in *Date, opts ...grpc.CallOption) (*HistoryTickResponse, error) {
	out := new(HistoryTickResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockTSEHistoryTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetStockTSEHistoryKbar(ctx context.Context, in *Date, opts ...grpc.CallOption) (*HistoryKbarResponse, error) {
	out := new(HistoryKbarResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockTSEHistoryKbar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetStockTSEHistoryClose(ctx context.Context, in *Date, opts ...grpc.CallOption) (*HistoryCloseResponse, error) {
	out := new(HistoryCloseResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetStockTSEHistoryClose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetFutureHistoryTick(ctx context.Context, in *FutureCodeArrWithDate, opts ...grpc.CallOption) (*HistoryTickResponse, error) {
	out := new(HistoryTickResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetFutureHistoryTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetFutureHistoryClose(ctx context.Context, in *FutureCodeArrWithDate, opts ...grpc.CallOption) (*HistoryCloseResponse, error) {
	out := new(HistoryCloseResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetFutureHistoryClose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *historyDataInterfaceClient) GetFutureHistoryKbar(ctx context.Context, in *FutureCodeArrWithDate, opts ...grpc.CallOption) (*HistoryKbarResponse, error) {
	out := new(HistoryKbarResponse)
	err := c.cc.Invoke(ctx, "/sinopac_forwarder.HistoryDataInterface/GetFutureHistoryKbar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HistoryDataInterfaceServer is the server API for HistoryDataInterface service.
// All implementations must embed UnimplementedHistoryDataInterfaceServer
// for forward compatibility
type HistoryDataInterfaceServer interface {
	GetStockHistoryTick(context.Context, *StockNumArrWithDate) (*HistoryTickResponse, error)
	GetStockHistoryKbar(context.Context, *StockNumArrWithDate) (*HistoryKbarResponse, error)
	GetStockHistoryClose(context.Context, *StockNumArrWithDate) (*HistoryCloseResponse, error)
	GetStockHistoryCloseByDateArr(context.Context, *StockNumArrWithDateArr) (*HistoryCloseResponse, error)
	GetStockTSEHistoryTick(context.Context, *Date) (*HistoryTickResponse, error)
	GetStockTSEHistoryKbar(context.Context, *Date) (*HistoryKbarResponse, error)
	GetStockTSEHistoryClose(context.Context, *Date) (*HistoryCloseResponse, error)
	GetFutureHistoryTick(context.Context, *FutureCodeArrWithDate) (*HistoryTickResponse, error)
	GetFutureHistoryClose(context.Context, *FutureCodeArrWithDate) (*HistoryCloseResponse, error)
	GetFutureHistoryKbar(context.Context, *FutureCodeArrWithDate) (*HistoryKbarResponse, error)
	mustEmbedUnimplementedHistoryDataInterfaceServer()
}

// UnimplementedHistoryDataInterfaceServer must be embedded to have forward compatible implementations.
type UnimplementedHistoryDataInterfaceServer struct {
}

func (UnimplementedHistoryDataInterfaceServer) GetStockHistoryTick(context.Context, *StockNumArrWithDate) (*HistoryTickResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockHistoryTick not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetStockHistoryKbar(context.Context, *StockNumArrWithDate) (*HistoryKbarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockHistoryKbar not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetStockHistoryClose(context.Context, *StockNumArrWithDate) (*HistoryCloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockHistoryClose not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetStockHistoryCloseByDateArr(context.Context, *StockNumArrWithDateArr) (*HistoryCloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockHistoryCloseByDateArr not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetStockTSEHistoryTick(context.Context, *Date) (*HistoryTickResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockTSEHistoryTick not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetStockTSEHistoryKbar(context.Context, *Date) (*HistoryKbarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockTSEHistoryKbar not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetStockTSEHistoryClose(context.Context, *Date) (*HistoryCloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockTSEHistoryClose not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetFutureHistoryTick(context.Context, *FutureCodeArrWithDate) (*HistoryTickResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFutureHistoryTick not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetFutureHistoryClose(context.Context, *FutureCodeArrWithDate) (*HistoryCloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFutureHistoryClose not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) GetFutureHistoryKbar(context.Context, *FutureCodeArrWithDate) (*HistoryKbarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFutureHistoryKbar not implemented")
}
func (UnimplementedHistoryDataInterfaceServer) mustEmbedUnimplementedHistoryDataInterfaceServer() {}

// UnsafeHistoryDataInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HistoryDataInterfaceServer will
// result in compilation errors.
type UnsafeHistoryDataInterfaceServer interface {
	mustEmbedUnimplementedHistoryDataInterfaceServer()
}

func RegisterHistoryDataInterfaceServer(s grpc.ServiceRegistrar, srv HistoryDataInterfaceServer) {
	s.RegisterService(&HistoryDataInterface_ServiceDesc, srv)
}

func _HistoryDataInterface_GetStockHistoryTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArrWithDate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryTick(ctx, req.(*StockNumArrWithDate))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetStockHistoryKbar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArrWithDate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryKbar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryKbar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryKbar(ctx, req.(*StockNumArrWithDate))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetStockHistoryClose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArrWithDate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryClose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryClose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryClose(ctx, req.(*StockNumArrWithDate))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetStockHistoryCloseByDateArr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockNumArrWithDateArr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryCloseByDateArr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockHistoryCloseByDateArr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockHistoryCloseByDateArr(ctx, req.(*StockNumArrWithDateArr))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetStockTSEHistoryTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Date)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockTSEHistoryTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockTSEHistoryTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockTSEHistoryTick(ctx, req.(*Date))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetStockTSEHistoryKbar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Date)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockTSEHistoryKbar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockTSEHistoryKbar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockTSEHistoryKbar(ctx, req.(*Date))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetStockTSEHistoryClose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Date)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetStockTSEHistoryClose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetStockTSEHistoryClose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetStockTSEHistoryClose(ctx, req.(*Date))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetFutureHistoryTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArrWithDate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetFutureHistoryTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetFutureHistoryTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetFutureHistoryTick(ctx, req.(*FutureCodeArrWithDate))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetFutureHistoryClose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArrWithDate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetFutureHistoryClose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetFutureHistoryClose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetFutureHistoryClose(ctx, req.(*FutureCodeArrWithDate))
	}
	return interceptor(ctx, in, info, handler)
}

func _HistoryDataInterface_GetFutureHistoryKbar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FutureCodeArrWithDate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HistoryDataInterfaceServer).GetFutureHistoryKbar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sinopac_forwarder.HistoryDataInterface/GetFutureHistoryKbar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HistoryDataInterfaceServer).GetFutureHistoryKbar(ctx, req.(*FutureCodeArrWithDate))
	}
	return interceptor(ctx, in, info, handler)
}

// HistoryDataInterface_ServiceDesc is the grpc.ServiceDesc for HistoryDataInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HistoryDataInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sinopac_forwarder.HistoryDataInterface",
	HandlerType: (*HistoryDataInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStockHistoryTick",
			Handler:    _HistoryDataInterface_GetStockHistoryTick_Handler,
		},
		{
			MethodName: "GetStockHistoryKbar",
			Handler:    _HistoryDataInterface_GetStockHistoryKbar_Handler,
		},
		{
			MethodName: "GetStockHistoryClose",
			Handler:    _HistoryDataInterface_GetStockHistoryClose_Handler,
		},
		{
			MethodName: "GetStockHistoryCloseByDateArr",
			Handler:    _HistoryDataInterface_GetStockHistoryCloseByDateArr_Handler,
		},
		{
			MethodName: "GetStockTSEHistoryTick",
			Handler:    _HistoryDataInterface_GetStockTSEHistoryTick_Handler,
		},
		{
			MethodName: "GetStockTSEHistoryKbar",
			Handler:    _HistoryDataInterface_GetStockTSEHistoryKbar_Handler,
		},
		{
			MethodName: "GetStockTSEHistoryClose",
			Handler:    _HistoryDataInterface_GetStockTSEHistoryClose_Handler,
		},
		{
			MethodName: "GetFutureHistoryTick",
			Handler:    _HistoryDataInterface_GetFutureHistoryTick_Handler,
		},
		{
			MethodName: "GetFutureHistoryClose",
			Handler:    _HistoryDataInterface_GetFutureHistoryClose_Handler,
		},
		{
			MethodName: "GetFutureHistoryKbar",
			Handler:    _HistoryDataInterface_GetFutureHistoryKbar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "history.proto",
}
