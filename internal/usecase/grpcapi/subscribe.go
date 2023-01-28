// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"tmt/internal/usecase"
	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// SubscribegRPCAPI -.
type SubscribegRPCAPI struct {
	conn *grpc.Connection
}

func NewSubscribe(client *grpc.Connection) usecase.SubscribegRPCAPI {
	return &SubscribegRPCAPI{client}
}

// SubscribeStockTick return arry means fail to subscribe
func (t *SubscribegRPCAPI) SubscribeStockTick(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.SubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockTick return arry means fail to subscribe
func (t *SubscribegRPCAPI) UnSubscribeStockTick(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockAllTick -.
func (t *SubscribegRPCAPI) UnSubscribeStockAllTick() (*pb.ErrorMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockAllTick(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SubscribeStockBidAsk return arry means fail to subscribe
func (t *SubscribegRPCAPI) SubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.SubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockBidAsk return arry means fail to subscribe
func (t *SubscribegRPCAPI) UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockAllBidAsk -.
func (t *SubscribegRPCAPI) UnSubscribeStockAllBidAsk() (*pb.ErrorMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockAllBidAsk(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SubscribeFutureTick return arry means fail to subscribe
func (t *SubscribegRPCAPI) SubscribeFutureTick(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.SubscribeFutureTick(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeFutureTick return arry means fail to subscribe
func (t *SubscribegRPCAPI) UnSubscribeFutureTick(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.UnSubscribeFutureTick(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// SubscribeFutureBidAsk return arry means fail to subscribe
func (t *SubscribegRPCAPI) SubscribeFutureBidAsk(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.SubscribeFutureBidAsk(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeFutureBidAsk return arry means fail to subscribe
func (t *SubscribegRPCAPI) UnSubscribeFutureBidAsk(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSubscribeDataInterfaceClient(conn)
	r, err := c.UnSubscribeFutureBidAsk(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}
