// Package grpc package grpc
package grpc

import (
	"context"

	"tmt/pkg/grpc"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"
)

// history -.
type history struct {
	pool *grpc.ConnPool
}

// NewHistory -.
func NewHistory(client *grpc.ConnPool) HistorygRPCAPI {
	return &history{client}
}

// GetStockHistoryTick GetStockHistoryTick
func (t *history) GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockHistoryTick(context.Background(), &pb.StockNumArrWithDate{
		StockNumArr: stockNumArr,
		Date:        date,
	})
	if err != nil {
		return []*pb.HistoryTickMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockHistoryKbar GetStockHistoryKbar
func (t *history) GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockHistoryKbar(context.Background(), &pb.StockNumArrWithDate{
		StockNumArr: stockNumArr,
		Date:        date,
	})
	if err != nil {
		return []*pb.HistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockHistoryClose GetStockHistoryClose
func (t *history) GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockHistoryClose(context.Background(), &pb.StockNumArrWithDate{
		StockNumArr: stockNumArr,
		Date:        date,
	})
	if err != nil {
		return []*pb.HistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}

// GetFutureHistoryKbar -.
func (t *history) GetFutureHistoryKbar(codeArr []string, date string) (*pb.HistoryKbarResponse, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetFutureHistoryKbar(context.Background(), &pb.FutureCodeArrWithDate{
		FutureCodeArr: codeArr,
		Date:          date,
	})
	if err != nil {
		return &pb.HistoryKbarResponse{}, err
	}
	return r, nil
}
