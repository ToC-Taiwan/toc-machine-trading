package grpcapi

import (
	"context"

	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// StreamgRPCAPI -.
type StreamgRPCAPI struct {
	conn *sinopac.Connection
}

// NewStream -.
func NewStream(client *sinopac.Connection) *StreamgRPCAPI {
	return &StreamgRPCAPI{client}
}

// GetAllStockSnapshot -.
func (t *StreamgRPCAPI) GetAllStockSnapshot() ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr -.
func (t *StreamgRPCAPI) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockSnapshotByNumArr(context.Background(), &pb.StockNumArr{
		StockNumArr: stockNumArr,
	})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotTSE -.
func (t *StreamgRPCAPI) GetStockSnapshotTSE() (*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.StockSnapshotMessage{}, err
	}
	return r, nil
}

// GetFutureSnapshotFIMTX -.
func (t *StreamgRPCAPI) GetFutureSnapshotFIMTX() (*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewFutureForwarderClient(conn)
	r, err := c.GetFIMTXSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.StockSnapshotMessage{}, err
	}
	return r, nil
}
