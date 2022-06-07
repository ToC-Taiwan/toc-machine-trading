// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BasicgRPCAPI -.
type BasicgRPCAPI struct {
	conn *sinopac.Connection
}

// NewBasic -.
func NewBasic(client *sinopac.Connection) *BasicgRPCAPI {
	return &BasicgRPCAPI{client}
}

// GetServerToken GetServerToken
func (t *BasicgRPCAPI) GetServerToken() (string, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetServerToken(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return r.GetToken(), nil
}

// GetAllStockDetail GetAllStockDetail
func (t *BasicgRPCAPI) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// GetAllStockSnapshot GetAllStockSnapshot
func (t *BasicgRPCAPI) GetAllStockSnapshot() ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr GetStockSnapshotByNumArr
func (t *BasicgRPCAPI) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.StockSnapshotMessage, error) {
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

// GetStockSnapshotTSE GetStockSnapshotTSE
func (t *BasicgRPCAPI) GetStockSnapshotTSE() ([]*pb.StockSnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockSnapshotMessage{}, err
	}
	return r.GetData(), nil
}
