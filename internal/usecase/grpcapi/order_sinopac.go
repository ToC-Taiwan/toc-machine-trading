package grpcapi

import (
	"context"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// OrdergRPCAPI -.
type OrdergRPCAPI struct {
	conn *grpc.Connection
}

// NewOrder -.
func NewOrder(client *grpc.Connection) *OrdergRPCAPI {
	return &OrdergRPCAPI{client}
}

// GetFuturePosition -.
func (t *OrdergRPCAPI) GetFuturePosition() (*pb.FuturePositionArr, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetFuturePosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyStock BuyStock
func (t *OrdergRPCAPI) BuyStock(order *entity.StockOrder, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
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
func (t *OrdergRPCAPI) SellStock(order *entity.StockOrder, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
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
func (t *OrdergRPCAPI) SellFirstStock(order *entity.StockOrder, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
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
	c := pb.NewTradeInterfaceClient(conn)
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
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.GetOrderStatusByID(context.Background(), &pb.OrderID{
		OrderId:  orderID,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetLocalOrderStatusArr -.
func (t *OrdergRPCAPI) GetLocalOrderStatusArr() error {
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
func (t *OrdergRPCAPI) GetSimulateOrderStatusArr() error {
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
func (t *OrdergRPCAPI) GetNonBlockOrderStatusArr() (*pb.ErrorMessage, error) {
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
func (t *OrdergRPCAPI) BuyFuture(order *entity.FutureOrder, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.BuyFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFuture -.
func (t *OrdergRPCAPI) SellFuture(order *entity.FutureOrder, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.SellFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFirstFuture -.
func (t *OrdergRPCAPI) SellFirstFuture(order *entity.FutureOrder, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.SellFirstFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Quantity,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// CancelFuture -.
func (t *OrdergRPCAPI) CancelFuture(orderID string, sim bool) (*pb.TradeResult, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewTradeInterfaceClient(conn)
	r, err := c.CancelFuture(context.Background(), &pb.FutureOrderID{
		OrderId:  orderID,
		Simulate: sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
