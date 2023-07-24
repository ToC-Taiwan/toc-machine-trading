// Package grpcapi package grpcapi
package grpcapi

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

// GetStockHistoryCloseByDateArr GetStockHistoryCloseByDateArr
func (t *HistorygRPCAPI) GetStockHistoryCloseByDateArr(stockNumArr []string, date []string) ([]*pb.HistoryCloseMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockHistoryCloseByDateArr(context.Background(), &pb.StockNumArrWithDateArr{
		StockNumArr: stockNumArr,
		DateArr:     date,
	})
	if err != nil {
		return []*pb.HistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryTick GetStockTSEHistoryTick
func (t *HistorygRPCAPI) GetStockTSEHistoryTick(date string) ([]*pb.HistoryTickMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockTSEHistoryTick(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.HistoryTickMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryKbar GetStockTSEHistoryKbar
func (t *HistorygRPCAPI) GetStockTSEHistoryKbar(date string) ([]*pb.HistoryKbarMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockTSEHistoryKbar(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.HistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryClose GetStockTSEHistoryClose
func (t *HistorygRPCAPI) GetStockTSEHistoryClose(date string) ([]*pb.HistoryCloseMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetStockTSEHistoryClose(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.HistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}

// GetFutureHistoryTick -.
func (t *HistorygRPCAPI) GetFutureHistoryTick(codeArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetFutureHistoryTick(context.Background(), &pb.FutureCodeArrWithDate{
		FutureCodeArr: codeArr,
		Date:          date,
	})
	if err != nil {
		return []*pb.HistoryTickMessage{}, err
	}
	return r.GetData(), nil
}

// GetFutureHistoryClose -.
func (t *HistorygRPCAPI) GetFutureHistoryClose(codeArr []string, date string) ([]*pb.HistoryCloseMessage, error) {
	conn := t.pool.Get()
	defer t.pool.Put(conn)

	r, err := pb.NewHistoryDataInterfaceClient(conn).GetFutureHistoryClose(context.Background(), &pb.FutureCodeArrWithDate{
		FutureCodeArr: codeArr,
		Date:          date,
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
