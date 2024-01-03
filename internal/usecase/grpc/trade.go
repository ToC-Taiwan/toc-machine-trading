// Package grpc package grpc
package grpc

import (
	"context"

	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TradegRPCAPI struct {
	pool *grpc.ConnPool
	sim  bool
}

func NewTrade(client *grpc.ConnPool, sim bool) *TradegRPCAPI {
	return &TradegRPCAPI{
		pool: client,
		sim:  sim,
	}
}

// GetFuturePosition -.
func (t *TradegRPCAPI) GetFuturePosition() (*pb.FuturePositionArr, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetFuturePosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetStockPosition -.
func (t *TradegRPCAPI) GetStockPosition() (*pb.StockPositionArr, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetStockPosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetAccountBalance -.
func (t *TradegRPCAPI) GetAccountBalance() (*pb.AccountBalance, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetAccountBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetMargin -.
func (t *TradegRPCAPI) GetMargin() (*pb.Margin, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetMargin(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetSettlement -.
func (t *TradegRPCAPI) GetSettlement() (*pb.SettlementList, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetSettlement(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyStock BuyStock
func (t *TradegRPCAPI) BuyStock(order *entity.StockOrder) (*pb.TradeResult, error) {
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
func (t *TradegRPCAPI) SellStock(order *entity.StockOrder) (*pb.TradeResult, error) {
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

func (t *TradegRPCAPI) BuyOddStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).BuyOddStock(context.Background(), &pb.OddStockOrderDetail{
		StockNum: order.StockNum,
		Price:    order.Price,
		Share:    order.Lot,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (t *TradegRPCAPI) SellOddStock(order *entity.StockOrder) (*pb.TradeResult, error) {
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
func (t *TradegRPCAPI) SellFirstStock(order *entity.StockOrder) (*pb.TradeResult, error) {
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

// CancelStock CancelStock
func (t *TradegRPCAPI) CancelStock(orderID string) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).CancelStock(context.Background(), &pb.OrderID{
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
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetOrderStatusByID(context.Background(), &pb.OrderID{
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
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewTradeInterfaceClient(conn).GetLocalOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetSimulateOrderStatusArr -.
func (t *TradegRPCAPI) GetSimulateOrderStatusArr() error {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	_, err := pb.NewTradeInterfaceClient(conn).GetSimulateOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetNonBlockOrderStatusArr GetNonBlockOrderStatusArr
func (t *TradegRPCAPI) GetNonBlockOrderStatusArr() (*pb.ErrorMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).GetNonBlockOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyFuture -.
func (t *TradegRPCAPI) BuyFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
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
func (t *TradegRPCAPI) SellFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
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
func (t *TradegRPCAPI) SellFirstFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
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

// CancelFuture -.
func (t *TradegRPCAPI) CancelFuture(orderID string) (*pb.TradeResult, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewTradeInterfaceClient(conn).CancelFuture(context.Background(), &pb.FutureOrderID{
		OrderId:  orderID,
		Simulate: t.sim,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
