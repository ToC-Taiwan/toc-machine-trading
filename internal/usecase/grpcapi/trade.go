// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TradegRPCAPI struct {
	conn *grpc.Connection
	sim  bool
}

func NewTrade(client *grpc.Connection, sim bool) *TradegRPCAPI {
	return &TradegRPCAPI{
		conn: client,
		sim:  sim,
	}
}

// GetFuturePosition -.
func (t *TradegRPCAPI) GetFuturePosition() (*pb.FuturePositionArr, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetFuturePosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetAccountBalance -.
func (t *TradegRPCAPI) GetAccountBalance() (*pb.AccountBalance, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetAccountBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetMargin -.
func (t *TradegRPCAPI) GetMargin() (*pb.Margin, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetMargin(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetSettlement -.
func (t *TradegRPCAPI) GetSettlement() (*pb.SettlementList, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetSettlement(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyStock BuyStock
func (t *TradegRPCAPI) BuyStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.BuyStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellStock SellStock
func (t *TradegRPCAPI) SellStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.SellStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFirstStock SellFirstStock
func (t *TradegRPCAPI) SellFirstStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.SellFirstStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// CancelStock CancelStock
func (t *TradegRPCAPI) CancelStock(orderID string) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.CancelStock(context.Background(), &pb.OrderID{
		OrderId:  orderID,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetOrderStatusByID GetOrderStatusByID
func (t *TradegRPCAPI) GetOrderStatusByID(orderID string) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetOrderStatusByID(context.Background(), &pb.OrderID{
		OrderId:  orderID,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetLocalOrderStatusArr -.
func (t *TradegRPCAPI) GetLocalOrderStatusArr() error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	_, err := c.GetLocalOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetSimulateOrderStatusArr -.
func (t *TradegRPCAPI) GetSimulateOrderStatusArr() error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	_, err := c.GetSimulateOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetNonBlockOrderStatusArr GetNonBlockOrderStatusArr
func (t *TradegRPCAPI) GetNonBlockOrderStatusArr() (*pb.ErrorMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetNonBlockOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyFuture -.
func (t *TradegRPCAPI) BuyFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.BuyFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFuture -.
func (t *TradegRPCAPI) SellFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.SellFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFirstFuture -.
func (t *TradegRPCAPI) SellFirstFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.SellFirstFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// CancelFuture -.
func (t *TradegRPCAPI) CancelFuture(orderID string) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.CancelFuture(context.Background(), &pb.FutureOrderID{
		OrderId:  orderID,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
