// Package grpc package grpc
package grpc

import (
	"context"

	"tmt/pb"
	"tmt/pkg/grpc"
)

// HistorygRPCAPI -.
type HistorygRPCAPI struct {
	pool *grpc.ConnPool
}

// NewHistory -.
func NewHistory(client *grpc.ConnPool) *HistorygRPCAPI {
	return &HistorygRPCAPI{client}
}

// GetStockHistoryTick GetStockHistoryTick
func (t *HistorygRPCAPI) GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error) {
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
func (t *HistorygRPCAPI) GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error) {
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
func (t *HistorygRPCAPI) GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error) {
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
func (t *HistorygRPCAPI) GetFutureHistoryKbar(codeArr []string, date string) ([]*pb.HistoryKbarMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetFutureHistoryKbar(context.Background(), &pb.FutureCodeArrWithDate{
		FutureCodeArr: codeArr,
		Date:          date,
	})
	if err != nil {
		return []*pb.HistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}
