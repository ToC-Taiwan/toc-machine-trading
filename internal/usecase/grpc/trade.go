// Package grpc package grpc
package grpc

import (
	"context"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

type trade struct {
	pool *grpc.ConnPool
	sim  bool
}

func NewTrade(client *grpc.ConnPool, sim bool) TradegRPCAPI {
	return &trade{
		pool: client,
		sim:  sim,
	}
}

// GetFuturePosition -.
func (t *trade) GetFuturePosition() (*pb.FuturePositionArr, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetFuturePosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetStockPosition -.
func (t *trade) GetStockPosition() (*pb.StockPositionArr, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetStockPosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetAccountBalance -.
func (t *trade) GetAccountBalance() (*pb.AccountBalance, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetAccountBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetMargin -.
func (t *trade) GetMargin() (*pb.Margin, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetMargin(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetSettlement -.
func (t *trade) GetSettlement() (*pb.SettlementList, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetSettlement(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyStock BuyStock
func (t *trade) BuyStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).BuyStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Lot,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellStock SellStock
func (t *trade) SellStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).SellStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Lot,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (t *trade) BuyOddStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).BuyOddStock(context.Background(), &pb.OddStockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Share:    order.Share,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (t *trade) SellOddStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).SellOddStock(context.Background(), &pb.OddStockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Share:    order.Share,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFirstStock SellFirstStock
func (t *trade) SellFirstStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).SellFirstStock(context.Background(), &pb.StockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Quantity: order.Lot,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// CancelOrder -.
func (t *trade) CancelOrder(orderID string) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).CancelOrder(context.Background(), &pb.OrderID{
		OrderId:  orderID,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetLocalOrderStatusArr -.
func (t *trade) GetLocalOrderStatusArr() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewTradeInterfaceClient(conn).GetLocalOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetSimulateOrderStatusArr -.
func (t *trade) GetSimulateOrderStatusArr() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewTradeInterfaceClient(conn).GetSimulateOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// BuyFuture -.
func (t *trade) BuyFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).BuyFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Position,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFuture -.
func (t *trade) SellFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).SellFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Position,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SellFirstFuture -.
func (t *trade) SellFirstFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).SellFirstFuture(context.Background(), &pb.FutureOrderDetail{
		Code:     order.Code,
		Price:    order.Price,
		Quantity: order.Position,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
