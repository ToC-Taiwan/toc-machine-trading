// Package grpc package grpc
package grpc

import (
	"context"

	"github.com/toc-taiwan/toc-trade-protobuf/golang/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// subscribe -.
type subscribe struct {
	conn *grpc.ClientConn
}

func NewSubscribe(client *grpc.ClientConn) SubscribegRPCAPI {
	return &subscribe{client}
}

// SubscribeStockTick return arry means fail to subscribe
func (t *subscribe) SubscribeStockTick(stockNumArr []string, odd bool) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).SubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr, Odd: odd})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockTick return arry means fail to subscribe
func (t *subscribe) UnSubscribeStockTick(stockNumArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).UnSubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeAllTick -.
func (t *subscribe) UnSubscribeAllTick() (*pb.ErrorMessage, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).UnSubscribeAllTick(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SubscribeStockBidAsk return arry means fail to subscribe
func (t *subscribe) SubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).SubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockBidAsk return arry means fail to subscribe
func (t *subscribe) UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).UnSubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeAllBidAsk -.
func (t *subscribe) UnSubscribeAllBidAsk() (*pb.ErrorMessage, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).UnSubscribeAllBidAsk(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SubscribeFutureTick return arry means fail to subscribe
func (t *subscribe) SubscribeFutureTick(codeArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).SubscribeFutureTick(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeFutureTick return arry means fail to subscribe
func (t *subscribe) UnSubscribeFutureTick(codeArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).UnSubscribeFutureTick(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// SubscribeFutureBidAsk return arry means fail to subscribe
func (t *subscribe) SubscribeFutureBidAsk(codeArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).SubscribeFutureBidAsk(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeFutureBidAsk return arry means fail to subscribe
func (t *subscribe) UnSubscribeFutureBidAsk(codeArr []string) ([]string, error) {
	r, err := pb.NewSubscribeDataInterfaceClient(t.conn).UnSubscribeFutureBidAsk(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}
