// Package grpcapi package grpcapi
package grpcapi

import (
	"context"
	"errors"

	"tmt/internal/usecase"
	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// RealTimegRPCAPI -.
type RealTimegRPCAPI struct {
	conn *grpc.Connection
}

func NewRealTime(client *grpc.Connection) usecase.RealTimegRPCAPI {
	return &RealTimegRPCAPI{client}
}

// GetAllStockSnapshot -.
func (t *RealTimegRPCAPI) GetAllStockSnapshot() ([]*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr -.
func (t *RealTimegRPCAPI) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetStockSnapshotByNumArr(context.Background(), &pb.StockNumArr{
		StockNumArr: stockNumArr,
	})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotTSE -.
func (t *RealTimegRPCAPI) GetStockSnapshotTSE() (*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.SnapshotMessage{}, err
	}
	return r, nil
}

// GetStockSnapshotOTC -.
func (t *RealTimegRPCAPI) GetStockSnapshotOTC() (*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetStockSnapshotOTC(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.SnapshotMessage{}, err
	}
	return r, nil
}

// GetFutureSnapshotByCode -.
func (t *RealTimegRPCAPI) GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
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

func (t *RealTimegRPCAPI) GetNasdaq() (*pb.YahooFinancePrice, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetNasdaq(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.YahooFinancePrice{}, err
	}
	return r, nil
}

func (t *RealTimegRPCAPI) GetNasdaqFuture() (*pb.YahooFinancePrice, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetNasdaqFuture(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.YahooFinancePrice{}, err
	}
	return r, nil
}

// GetStockVolumeRank -.
func (t *RealTimegRPCAPI) GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewRealTimeDataInterfaceClient(conn)
	r, err := c.GetStockVolumeRank(context.Background(), &pb.VolumeRankRequest{
		Count: 200,
		Date:  date,
	})
	if err != nil {
		return []*pb.StockVolumeRankMessage{}, err
	}
	return r.GetData(), nil
}
