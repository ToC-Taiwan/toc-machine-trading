package grpcapi

import (
	"context"

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

// GetFutureSnapshotByCodeArr -.
func (t *StreamgRPCAPI) GetFutureSnapshotByCodeArr(codeArr []string) (*pb.SnapshotResponse, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewStreamDataInterfaceClient(conn)
	r, err := c.GetFutureSnapshotByCodeArr(context.Background(), &pb.FutureCodeArr{
		FutureCodeArr: codeArr,
	})
	if err != nil {
		return &pb.SnapshotResponse{}, err
	}
	return r, nil
}
