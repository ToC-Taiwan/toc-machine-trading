// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.11
// source: trade.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// code (str): contract id.
// direction (Action): action. {Buy, Sell}
// quantity (int): quantity.
// price (float): the average price.
// last_price (float): last price.
// pnl (float): unrealized profit.
type FuturePosition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code      string  `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Direction string  `protobuf:"bytes,2,opt,name=direction,proto3" json:"direction,omitempty"`
	Quantity  int32   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price     float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	LastPrice float64 `protobuf:"fixed64,5,opt,name=last_price,json=lastPrice,proto3" json:"last_price,omitempty"`
	Pnl       float64 `protobuf:"fixed64,6,opt,name=pnl,proto3" json:"pnl,omitempty"`
}

func (x *FuturePosition) Reset() {
	*x = FuturePosition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FuturePosition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FuturePosition) ProtoMessage() {}

func (x *FuturePosition) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FuturePosition.ProtoReflect.Descriptor instead.
func (*FuturePosition) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{0}
}

func (x *FuturePosition) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *FuturePosition) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

func (x *FuturePosition) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *FuturePosition) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *FuturePosition) GetLastPrice() float64 {
	if x != nil {
		return x.LastPrice
	}
	return 0
}

func (x *FuturePosition) GetPnl() float64 {
	if x != nil {
		return x.Pnl
	}
	return 0
}

type FuturePositionArr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PositionArr []*FuturePosition `protobuf:"bytes,1,rep,name=position_arr,json=positionArr,proto3" json:"position_arr,omitempty"`
}

func (x *FuturePositionArr) Reset() {
	*x = FuturePositionArr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FuturePositionArr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FuturePositionArr) ProtoMessage() {}

func (x *FuturePositionArr) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FuturePositionArr.ProtoReflect.Descriptor instead.
func (*FuturePositionArr) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{1}
}

func (x *FuturePositionArr) GetPositionArr() []*FuturePosition {
	if x != nil {
		return x.PositionArr
	}
	return nil
}

type StockOrderDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StockNum string  `protobuf:"bytes,1,opt,name=stock_num,json=stockNum,proto3" json:"stock_num,omitempty"`
	Price    float64 `protobuf:"fixed64,2,opt,name=price,proto3" json:"price,omitempty"`
	Quantity int64   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Simulate bool    `protobuf:"varint,4,opt,name=simulate,proto3" json:"simulate,omitempty"`
}

func (x *StockOrderDetail) Reset() {
	*x = StockOrderDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StockOrderDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockOrderDetail) ProtoMessage() {}

func (x *StockOrderDetail) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockOrderDetail.ProtoReflect.Descriptor instead.
func (*StockOrderDetail) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{2}
}

func (x *StockOrderDetail) GetStockNum() string {
	if x != nil {
		return x.StockNum
	}
	return ""
}

func (x *StockOrderDetail) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *StockOrderDetail) GetQuantity() int64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *StockOrderDetail) GetSimulate() bool {
	if x != nil {
		return x.Simulate
	}
	return false
}

type FutureOrderDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     string  `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Price    float64 `protobuf:"fixed64,2,opt,name=price,proto3" json:"price,omitempty"`
	Quantity int64   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Simulate bool    `protobuf:"varint,4,opt,name=simulate,proto3" json:"simulate,omitempty"`
}

func (x *FutureOrderDetail) Reset() {
	*x = FutureOrderDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FutureOrderDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FutureOrderDetail) ProtoMessage() {}

func (x *FutureOrderDetail) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FutureOrderDetail.ProtoReflect.Descriptor instead.
func (*FutureOrderDetail) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{3}
}

func (x *FutureOrderDetail) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *FutureOrderDetail) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *FutureOrderDetail) GetQuantity() int64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *FutureOrderDetail) GetSimulate() bool {
	if x != nil {
		return x.Simulate
	}
	return false
}

type TradeResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Status  string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Error   string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *TradeResult) Reset() {
	*x = TradeResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradeResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradeResult) ProtoMessage() {}

func (x *TradeResult) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradeResult.ProtoReflect.Descriptor instead.
func (*TradeResult) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{4}
}

func (x *TradeResult) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *TradeResult) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *TradeResult) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type OrderID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId  string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Simulate bool   `protobuf:"varint,2,opt,name=simulate,proto3" json:"simulate,omitempty"`
}

func (x *OrderID) Reset() {
	*x = OrderID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderID) ProtoMessage() {}

func (x *OrderID) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderID.ProtoReflect.Descriptor instead.
func (*OrderID) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{5}
}

func (x *OrderID) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *OrderID) GetSimulate() bool {
	if x != nil {
		return x.Simulate
	}
	return false
}

type FutureOrderID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId  string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Simulate bool   `protobuf:"varint,2,opt,name=simulate,proto3" json:"simulate,omitempty"`
}

func (x *FutureOrderID) Reset() {
	*x = FutureOrderID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FutureOrderID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FutureOrderID) ProtoMessage() {}

func (x *FutureOrderID) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FutureOrderID.ProtoReflect.Descriptor instead.
func (*FutureOrderID) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{6}
}

func (x *FutureOrderID) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *FutureOrderID) GetSimulate() bool {
	if x != nil {
		return x.Simulate
	}
	return false
}

type OrderStatusArr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*OrderStatus `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *OrderStatusArr) Reset() {
	*x = OrderStatusArr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderStatusArr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderStatusArr) ProtoMessage() {}

func (x *OrderStatusArr) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderStatusArr.ProtoReflect.Descriptor instead.
func (*OrderStatusArr) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{7}
}

func (x *OrderStatusArr) GetData() []*OrderStatus {
	if x != nil {
		return x.Data
	}
	return nil
}

type OrderStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status    string  `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Code      string  `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Action    string  `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`
	Price     float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	Quantity  int64   `protobuf:"varint,5,opt,name=quantity,proto3" json:"quantity,omitempty"`
	OrderId   string  `protobuf:"bytes,6,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	OrderTime string  `protobuf:"bytes,7,opt,name=order_time,json=orderTime,proto3" json:"order_time,omitempty"`
}

func (x *OrderStatus) Reset() {
	*x = OrderStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderStatus) ProtoMessage() {}

func (x *OrderStatus) ProtoReflect() protoreflect.Message {
	mi := &file_trade_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderStatus.ProtoReflect.Descriptor instead.
func (*OrderStatus) Descriptor() ([]byte, []int) {
	return file_trade_proto_rawDescGZIP(), []int{8}
}

func (x *OrderStatus) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *OrderStatus) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *OrderStatus) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *OrderStatus) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *OrderStatus) GetQuantity() int64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *OrderStatus) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *OrderStatus) GetOrderTime() string {
	if x != nil {
		return x.OrderTime
	}
	return ""
}

var File_trade_proto protoreflect.FileDescriptor

var file_trade_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x73,
	0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa5, 0x01, 0x0a, 0x0e,
	0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x6e, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03,
	0x70, 0x6e, 0x6c, 0x22, 0x59, 0x0a, 0x11, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x72, 0x72, 0x12, 0x44, 0x0a, 0x0c, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x72, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0b, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x72, 0x72, 0x22, 0x7d,
	0x0a, 0x10, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x22, 0x75, 0x0a,
	0x11, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x73, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x65, 0x22, 0x56, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x40, 0x0a, 0x07,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x22, 0x46,
	0x0a, 0x0d, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12,
	0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69,
	0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x73, 0x69,
	0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x22, 0x44, 0x0a, 0x0e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x41, 0x72, 0x72, 0x12, 0x32, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63,
	0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xbd, 0x01, 0x0a,
	0x0b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x32, 0xd3, 0x08, 0x0a,
	0x0e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12,
	0x51, 0x0a, 0x08, 0x42, 0x75, 0x79, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x23, 0x2e, 0x73, 0x69,
	0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x00, 0x12, 0x52, 0x0a, 0x09, 0x53, 0x65, 0x6c, 0x6c, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12,
	0x23, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72,
	0x64, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66,
	0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x0e, 0x53, 0x65, 0x6c, 0x6c, 0x46, 0x69,
	0x72, 0x73, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x23, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70,
	0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x1a, 0x1e, 0x2e,
	0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12,
	0x4b, 0x0a, 0x0b, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x1a,
	0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6e,
	0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x54,
	0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x16,
	0x47, 0x65, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x41, 0x72, 0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x53,
	0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x41, 0x72, 0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x49, 0x44, 0x12, 0x1a, 0x2e,
	0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6e, 0x6f,
	0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72,
	0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x19, 0x47,
	0x65, 0x74, 0x4e, 0x6f, 0x6e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x41, 0x72, 0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x1f, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x65, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x09, 0x42, 0x75, 0x79, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65,
	0x12, 0x24, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x65, 0x72, 0x2e, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63,
	0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0a, 0x53, 0x65, 0x6c, 0x6c,
	0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x12, 0x24, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63,
	0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x46, 0x75, 0x74, 0x75, 0x72,
	0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x1a, 0x1e, 0x2e, 0x73,
	0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72,
	0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x59,
	0x0a, 0x0f, 0x53, 0x65, 0x6c, 0x6c, 0x46, 0x69, 0x72, 0x73, 0x74, 0x46, 0x75, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x24, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x1a, 0x1e, 0x2e, 0x73, 0x69, 0x6e, 0x6f, 0x70, 0x61,
	0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x64,
	0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x0c, 0x43, 0x61, 0x6e,
	0x63, 0x65, 0x6c, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x12, 0x20, 0x2e, 0x73, 0x69, 0x6e, 0x6f,
	0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x46, 0x75,
	0x74, 0x75, 0x72, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x1e, 0x2e, 0x73, 0x69,
	0x6e, 0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e,
	0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x53, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x24, 0x2e, 0x73, 0x69, 0x6e,
	0x6f, 0x70, 0x61, 0x63, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x46,
	0x75, 0x74, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x72, 0x72,
	0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_trade_proto_rawDescOnce sync.Once
	file_trade_proto_rawDescData = file_trade_proto_rawDesc
)

func file_trade_proto_rawDescGZIP() []byte {
	file_trade_proto_rawDescOnce.Do(func() {
		file_trade_proto_rawDescData = protoimpl.X.CompressGZIP(file_trade_proto_rawDescData)
	})
	return file_trade_proto_rawDescData
}

var file_trade_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_trade_proto_goTypes = []interface{}{
	(*FuturePosition)(nil),    // 0: sinopac_forwarder.FuturePosition
	(*FuturePositionArr)(nil), // 1: sinopac_forwarder.FuturePositionArr
	(*StockOrderDetail)(nil),  // 2: sinopac_forwarder.StockOrderDetail
	(*FutureOrderDetail)(nil), // 3: sinopac_forwarder.FutureOrderDetail
	(*TradeResult)(nil),       // 4: sinopac_forwarder.TradeResult
	(*OrderID)(nil),           // 5: sinopac_forwarder.OrderID
	(*FutureOrderID)(nil),     // 6: sinopac_forwarder.FutureOrderID
	(*OrderStatusArr)(nil),    // 7: sinopac_forwarder.OrderStatusArr
	(*OrderStatus)(nil),       // 8: sinopac_forwarder.OrderStatus
	(*emptypb.Empty)(nil),     // 9: google.protobuf.Empty
	(*ErrorMessage)(nil),      // 10: sinopac_forwarder.ErrorMessage
}
var file_trade_proto_depIdxs = []int32{
	0,  // 0: sinopac_forwarder.FuturePositionArr.position_arr:type_name -> sinopac_forwarder.FuturePosition
	8,  // 1: sinopac_forwarder.OrderStatusArr.data:type_name -> sinopac_forwarder.OrderStatus
	2,  // 2: sinopac_forwarder.TradeInterface.BuyStock:input_type -> sinopac_forwarder.StockOrderDetail
	2,  // 3: sinopac_forwarder.TradeInterface.SellStock:input_type -> sinopac_forwarder.StockOrderDetail
	2,  // 4: sinopac_forwarder.TradeInterface.SellFirstStock:input_type -> sinopac_forwarder.StockOrderDetail
	5,  // 5: sinopac_forwarder.TradeInterface.CancelStock:input_type -> sinopac_forwarder.OrderID
	9,  // 6: sinopac_forwarder.TradeInterface.GetLocalOrderStatusArr:input_type -> google.protobuf.Empty
	9,  // 7: sinopac_forwarder.TradeInterface.GetSimulateOrderStatusArr:input_type -> google.protobuf.Empty
	5,  // 8: sinopac_forwarder.TradeInterface.GetOrderStatusByID:input_type -> sinopac_forwarder.OrderID
	9,  // 9: sinopac_forwarder.TradeInterface.GetNonBlockOrderStatusArr:input_type -> google.protobuf.Empty
	3,  // 10: sinopac_forwarder.TradeInterface.BuyFuture:input_type -> sinopac_forwarder.FutureOrderDetail
	3,  // 11: sinopac_forwarder.TradeInterface.SellFuture:input_type -> sinopac_forwarder.FutureOrderDetail
	3,  // 12: sinopac_forwarder.TradeInterface.SellFirstFuture:input_type -> sinopac_forwarder.FutureOrderDetail
	6,  // 13: sinopac_forwarder.TradeInterface.CancelFuture:input_type -> sinopac_forwarder.FutureOrderID
	9,  // 14: sinopac_forwarder.TradeInterface.GetFuturePosition:input_type -> google.protobuf.Empty
	4,  // 15: sinopac_forwarder.TradeInterface.BuyStock:output_type -> sinopac_forwarder.TradeResult
	4,  // 16: sinopac_forwarder.TradeInterface.SellStock:output_type -> sinopac_forwarder.TradeResult
	4,  // 17: sinopac_forwarder.TradeInterface.SellFirstStock:output_type -> sinopac_forwarder.TradeResult
	4,  // 18: sinopac_forwarder.TradeInterface.CancelStock:output_type -> sinopac_forwarder.TradeResult
	9,  // 19: sinopac_forwarder.TradeInterface.GetLocalOrderStatusArr:output_type -> google.protobuf.Empty
	9,  // 20: sinopac_forwarder.TradeInterface.GetSimulateOrderStatusArr:output_type -> google.protobuf.Empty
	4,  // 21: sinopac_forwarder.TradeInterface.GetOrderStatusByID:output_type -> sinopac_forwarder.TradeResult
	10, // 22: sinopac_forwarder.TradeInterface.GetNonBlockOrderStatusArr:output_type -> sinopac_forwarder.ErrorMessage
	4,  // 23: sinopac_forwarder.TradeInterface.BuyFuture:output_type -> sinopac_forwarder.TradeResult
	4,  // 24: sinopac_forwarder.TradeInterface.SellFuture:output_type -> sinopac_forwarder.TradeResult
	4,  // 25: sinopac_forwarder.TradeInterface.SellFirstFuture:output_type -> sinopac_forwarder.TradeResult
	4,  // 26: sinopac_forwarder.TradeInterface.CancelFuture:output_type -> sinopac_forwarder.TradeResult
	1,  // 27: sinopac_forwarder.TradeInterface.GetFuturePosition:output_type -> sinopac_forwarder.FuturePositionArr
	15, // [15:28] is the sub-list for method output_type
	2,  // [2:15] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_trade_proto_init() }
func file_trade_proto_init() {
	if File_trade_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_trade_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FuturePosition); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FuturePositionArr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StockOrderDetail); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FutureOrderDetail); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradeResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FutureOrderID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderStatusArr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trade_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_trade_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trade_proto_goTypes,
		DependencyIndexes: file_trade_proto_depIdxs,
		MessageInfos:      file_trade_proto_msgTypes,
	}.Build()
	File_trade_proto = out.File
	file_trade_proto_rawDesc = nil
	file_trade_proto_goTypes = nil
	file_trade_proto_depIdxs = nil
}
