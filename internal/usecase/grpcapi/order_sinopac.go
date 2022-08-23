package grpcapi

import (
	"context"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// OrdergRPCAPI -.
type OrdergRPCAPI struct {
	conn *sinopac.Connection
}

// NewOrder -.
func NewOrder(client *sinopac.Connection) *OrdergRPCAPI {
	return &OrdergRPCAPI{client}
}

// BuyStock BuyStock
func (t *OrdergRPCAPI) BuyStock(order *entity.Order, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.BuyStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellStock SellStock
func (t *OrdergRPCAPI) SellStock(order *entity.Order, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.SellStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFirstStock SellFirstStock
func (t *OrdergRPCAPI) SellFirstStock(order *entity.Order, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.SellFirstStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// CancelStock CancelStock
func (t *OrdergRPCAPI) CancelStock(orderID string, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.CancelStock(context.Background(), &pb.OrderID{
		OrderId:  orderID,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetOrderStatusByID GetOrderStatusByID
func (t *OrdergRPCAPI) GetOrderStatusByID(orderID string, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.GetOrderStatusByID(context.Background(), &pb.OrderID{
		OrderId:  orderID,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetOrderStatusArr GetOrderStatusArr
func (t *OrdergRPCAPI) GetOrderStatusArr() ([]*pb.StockOrderStatus, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.GetOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockOrderStatus{}, err
	}
	return r.GetData(), nil
}

// GetNonBlockOrderStatusArr GetNonBlockOrderStatusArr
func (t *OrdergRPCAPI) GetNonBlockOrderStatusArr() (*pb.FunctionErr, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeServiceClient(conn)
	r, err := c.GetNonBlockOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}
