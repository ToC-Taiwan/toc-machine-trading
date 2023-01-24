// Package sinopac package sinopac
package sinopac

import (
	"context"

	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// TargetgRPCAPI -.
type TargetgRPCAPI struct {
	conn *grpc.Connection
}

// NewTarget -.
func NewTarget(client *grpc.Connection) *TargetgRPCAPI {
	return &TargetgRPCAPI{client}
}

// GetStockVolumeRank -.
func (t *TargetgRPCAPI) GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.GetStockVolumeRank(context.Background(), &pb.VolumeRankRequest{
		Count: 200,
		Date:  date,
	})
	if err != nil {
		return []*pb.StockVolumeRankMessage{}, err
	}
	return r.GetData(), nil
}

// SubscribeStockTick return arry means fail to subscribe
func (t *TargetgRPCAPI) SubscribeStockTick(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.SubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockTick return arry means fail to subscribe
func (t *TargetgRPCAPI) UnSubscribeStockTick(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockAllTick -.
func (t *TargetgRPCAPI) UnSubscribeStockAllTick() (*pb.ErrorMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockAllTick(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SubscribeStockBidAsk return arry means fail to subscribe
func (t *TargetgRPCAPI) SubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.SubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockBidAsk return arry means fail to subscribe
func (t *TargetgRPCAPI) UnSubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeStockAllBidAsk -.
func (t *TargetgRPCAPI) UnSubscribeStockAllBidAsk() (*pb.ErrorMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.UnSubscribeStockAllBidAsk(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SubscribeFutureTick return arry means fail to subscribe
func (t *TargetgRPCAPI) SubscribeFutureTick(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.SubscribeFutureTick(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeFutureTick return arry means fail to subscribe
func (t *TargetgRPCAPI) UnSubscribeFutureTick(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.UnSubscribeFutureTick(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// SubscribeFutureBidAsk return arry means fail to subscribe
func (t *TargetgRPCAPI) SubscribeFutureBidAsk(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.SubscribeFutureBidAsk(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// UnSubscribeFutureBidAsk return arry means fail to subscribe
func (t *TargetgRPCAPI) UnSubscribeFutureBidAsk(codeArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.UnSubscribeFutureBidAsk(context.Background(), &pb.FutureCodeArr{FutureCodeArr: codeArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}
