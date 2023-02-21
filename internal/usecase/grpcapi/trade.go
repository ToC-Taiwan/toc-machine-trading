// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/grpc"

	"github.com/slack-go/slack"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TradegRPCAPI struct {
	conn *grpc.Connection
	sim  bool

	slack          *slack.Client
	slackChannelID string
}

func NewTrade(client *grpc.Connection, slackCfg config.Slack, sim bool) *TradegRPCAPI {
	return &TradegRPCAPI{
		conn:           client,
		sim:            sim,
		slack:          slack.New(slackCfg.Token),
		slackChannelID: slackCfg.ChannelID,
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

func (t *TradegRPCAPI) NotifyToSlack(message string) {
	_, _, e := t.slack.PostMessage(t.slackChannelID, slack.MsgOptionText(message, false))
	if e != nil {
		logger.Errorf("NotifyToSlack error: %v", e)
	}
}
