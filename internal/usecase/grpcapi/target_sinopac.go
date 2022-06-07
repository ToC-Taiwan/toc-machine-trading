package grpcapi

import (
	"context"

	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/sinopac"
)

// TargetgRPCAPI -.
type TargetgRPCAPI struct {
	conn *sinopac.Connection
}

// NewTarget -.
func NewTarget(client *sinopac.Connection) *TargetgRPCAPI {
	return &TargetgRPCAPI{client}
}

// GetStockVolumeRank -.
func (t *TargetgRPCAPI) GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
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
	c := pb.NewSinopacForwarderClient(conn)
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
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.UnSubscribeStockTick(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}

// SubscribeStockBidAsk return arry means fail to subscribe
func (t *TargetgRPCAPI) SubscribeStockBidAsk(stockNumArr []string) ([]string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
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
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.UnSubscribeStockBidAsk(context.Background(), &pb.StockNumArr{StockNumArr: stockNumArr})
	if err != nil {
		return []string{}, err
	}
	return r.GetFailArr(), nil
}
