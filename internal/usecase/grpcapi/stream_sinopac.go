package grpcapi

import (
	"context"
	"errors"

	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// StreamgRPCAPI -.
type StreamgRPCAPI struct {
	conn *grpc.Connection
}

// NewStream -.
func NewStream(client *grpc.Connection) *StreamgRPCAPI {
	return &StreamgRPCAPI{client}
}

// GetAllStockSnapshot -.
func (t *StreamgRPCAPI) GetAllStockSnapshot() ([]*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr -.
func (t *StreamgRPCAPI) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.GetStockSnapshotByNumArr(context.Background(), &pb.StockNumArr{
		StockNumArr: stockNumArr,
	})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotTSE -.
func (t *StreamgRPCAPI) GetStockSnapshotTSE() (*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.SnapshotMessage{}, err
	}
	return r, nil
}

// GetFutureSnapshotByCode -.
func (t *StreamgRPCAPI) GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.GetFutureSnapshotByCodeArr(context.Background(), &pb.FutureCodeArr{
		FutureCodeArr: []string{code},
	})
	if err != nil {
		return &pb.SnapshotMessage{}, err
	}

	if data := r.GetData(); len(data) > 0 {
		return data[0], nil
	}
	return nil, errors.New("no data")
}
