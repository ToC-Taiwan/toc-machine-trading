// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"tmt/internal/usecase"
	"tmt/pb"
	"tmt/pkg/grpc"
)

// HistorygRPCAPI -.
type HistorygRPCAPI struct {
	conn *grpc.Connection
}

// NewHistory -.
func NewHistory(client *grpc.Connection) usecase.HistorygRPCAPI {
	return &HistorygRPCAPI{client}
}

// GetStockHistoryTick GetStockHistoryTick
func (t *HistorygRPCAPI) GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockHistoryTick(context.Background(), &pb.StockNumArrWithDate{
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
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockHistoryKbar(context.Background(), &pb.StockNumArrWithDate{
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
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockHistoryClose(context.Background(), &pb.StockNumArrWithDate{
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
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockHistoryCloseByDateArr(context.Background(), &pb.StockNumArrWithDateArr{
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
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockTSEHistoryTick(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.HistoryTickMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryKbar GetStockTSEHistoryKbar
func (t *HistorygRPCAPI) GetStockTSEHistoryKbar(date string) ([]*pb.HistoryKbarMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockTSEHistoryKbar(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.HistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryClose GetStockTSEHistoryClose
func (t *HistorygRPCAPI) GetStockTSEHistoryClose(date string) ([]*pb.HistoryCloseMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetStockTSEHistoryClose(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.HistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}

// GetFutureHistoryTick -.
func (t *HistorygRPCAPI) GetFutureHistoryTick(codeArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetFutureHistoryTick(context.Background(), &pb.FutureCodeArrWithDate{
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
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetFutureHistoryClose(context.Background(), &pb.FutureCodeArrWithDate{
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
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewHistoryDataInterfaceClient(conn)
	r, err := c.GetFutureHistoryKbar(context.Background(), &pb.FutureCodeArrWithDate{
		FutureCodeArr: codeArr,
		Date:          date,
	})
	if err != nil {
		return []*pb.HistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}
