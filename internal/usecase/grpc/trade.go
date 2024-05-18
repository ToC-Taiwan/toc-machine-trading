// Package grpc package grpc
package grpc

import (
	"context"

	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type trade struct {
	conn *grpc.ClientConn
	sim  bool
}

func NewTrade(client *grpc.ClientConn, sim bool) TradegRPCAPI {
	return &trade{
		conn: client,
		sim:  sim,
	}
}

// GetFuturePosition -.
func (t *trade) GetFuturePosition() (*pb.FuturePositionArr, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).GetFuturePosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetStockPosition -.
func (t *trade) GetStockPosition() (*pb.StockPositionArr, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).GetStockPosition(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetAccountBalance -.
func (t *trade) GetAccountBalance() (*pb.AccountBalance, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).GetAccountBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetMargin -.
func (t *trade) GetMargin() (*pb.Margin, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).GetMargin(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetSettlement -.
func (t *trade) GetSettlement() (*pb.SettlementList, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).GetSettlement(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// BuyStock BuyStock
func (t *trade) BuyStock(order *entity.StockOrder) (*pb.TradeResult, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).BuyStock(context.Background(), &pb.StockOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).SellStock(context.Background(), &pb.StockOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).BuyOddStock(context.Background(), &pb.OddStockOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).SellOddStock(context.Background(), &pb.OddStockOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).SellFirstStock(context.Background(), &pb.StockOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).CancelOrder(context.Background(), &pb.OrderID{
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
	_, err := pb.NewTradeInterfaceClient(t.conn).GetLocalOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// GetSimulateOrderStatusArr -.
func (t *trade) GetSimulateOrderStatusArr() error {
	_, err := pb.NewTradeInterfaceClient(t.conn).GetSimulateOrderStatusArr(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// BuyFuture -.
func (t *trade) BuyFuture(order *entity.FutureOrder) (*pb.TradeResult, error) {
	r, err := pb.NewTradeInterfaceClient(t.conn).BuyFuture(context.Background(), &pb.FutureOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).SellFuture(context.Background(), &pb.FutureOrderDetail{
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
	r, err := pb.NewTradeInterfaceClient(t.conn).SellFirstFuture(context.Background(), &pb.FutureOrderDetail{
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
