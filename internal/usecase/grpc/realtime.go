// Package grpc package grpc
package grpc

import (
	"context"
	"errors"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// realtime -.
type realtime struct {
	conn *grpc.ClientConn
}

func NewRealTime(client *grpc.ClientConn) RealTimegRPCAPI {
	return &realtime{client}
}

// GetAllStockSnapshot -.
func (t *realtime) GetAllStockSnapshot() ([]*pb.SnapshotMessage, error) {
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetAllStockSnapshot(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotByNumArr -.
func (t *realtime) GetStockSnapshotByNumArr(stockNumArr []string) ([]*pb.SnapshotMessage, error) {
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetStockSnapshotByNumArr(context.Background(), &pb.StockNumArr{
		StockNumArr: stockNumArr,
	})
	if err != nil {
		return []*pb.SnapshotMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockSnapshotTSE -.
func (t *realtime) GetStockSnapshotTSE() (*pb.SnapshotMessage, error) {
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetStockSnapshotTSE(context.Background(), &emptypb.Empty{})
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
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetStockSnapshotOTC(context.Background(), &emptypb.Empty{})
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
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetFutureSnapshotByCodeArr(context.Background(), &pb.FutureCodeArr{
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
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetNasdaq(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.YahooFinancePrice{}, err
	}
	return r, nil
}

func (t *realtime) GetNasdaqFuture() (*pb.YahooFinancePrice, error) {
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetNasdaqFuture(context.Background(), &emptypb.Empty{})
	if err != nil {
		return &pb.YahooFinancePrice{}, err
	}
	return r, nil
}

// GetStockVolumeRank -.
func (t *realtime) GetStockVolumeRank(date string) ([]*pb.StockVolumeRankMessage, error) {
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetStockVolumeRank(context.Background(), &pb.VolumeRankRequest{
		Count: 200,
		Date:  date,
	})
	if err != nil {
		return []*pb.StockVolumeRankMessage{}, err
	}
	return r.GetData(), nil
}

func (t *realtime) GetStockVolumeRankPB(date string) (*pb.StockVolumeRankResponse, error) {
	r, err := pb.NewRealTimeDataInterfaceClient(t.conn).GetStockVolumeRank(context.Background(), &pb.VolumeRankRequest{
		Count: 200,
		Date:  date,
	})
	if err != nil {
		return &pb.StockVolumeRankResponse{}, err
	}
	return r, nil
}
