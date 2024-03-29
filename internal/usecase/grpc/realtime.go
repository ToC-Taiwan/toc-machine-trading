// Package grpc package grpc
package grpc

import (
	"context"
	"errors"

	"tmt/pb"
	"tmt/pkg/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// realtime -.
type realtime struct {
	pool *grpc.ConnPool
}

func NewRealTime(client *grpc.ConnPool) RealTimegRPCAPI {
	return &realtime{client}
}

// GetAllStockSnapshot -.
func (t *realtime) GetAllStockSnapshot() ([]*pb.SnapshotMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr -.
func (t *realtime) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.SnapshotMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetStockSnapshotByNumArr(context.Background(), &pb.StockNumArr{
		StockNumArr: stockNumArr,
	})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotTSE -.
func (t *realtime) GetStockSnapshotTSE() (*pb.SnapshotMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.SnapshotMessage{}, err
	}
	if len(r.GetData()) > 0 {
		return r.GetData()[0], nil
	}
	return &pb.SnapshotMessage{}, nil
}

// GetStockSnapshotOTC -.
func (t *realtime) GetStockSnapshotOTC() (*pb.SnapshotMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetStockSnapshotOTC(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.SnapshotMessage{}, err
	}
	if len(r.GetData()) > 0 {
		return r.GetData()[0], nil
	}
	return &pb.SnapshotMessage{}, nil
}

// GetFutureSnapshotByCode -.
func (t *realtime) GetFutureSnapshotByCode(code string) (*pb.SnapshotMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetFutureSnapshotByCodeArr(context.Background(), &pb.FutureCodeArr{
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

func (t *realtime) GetNasdaq() (*pb.YahooFinancePrice, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetNasdaq(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.YahooFinancePrice{}, err
	}
	return r, nil
}

func (t *realtime) GetNasdaqFuture() (*pb.YahooFinancePrice, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetNasdaqFuture(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.YahooFinancePrice{}, err
	}
	return r, nil
}

// GetStockVolumeRank -.
func (t *realtime) GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetStockVolumeRank(context.Background(), &pb.VolumeRankRequest{
		Count: 200,
		Date:  date,
	})
	if err != nil {
		return []*pb.StockVolumeRankMessage{}, err
	}
	return r.GetData(), nil
}

func (t *realtime) GetStockVolumeRankPB(date string) (*pb.StockVolumeRankResponse, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewRealTimeDataInterfaceClient(conn).GetStockVolumeRank(context.Background(), &pb.VolumeRankRequest{
		Count: 200,
		Date:  date,
	})
	if err != nil {
		return &pb.StockVolumeRankResponse{}, err
	}
	return r, nil
}
