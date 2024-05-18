// Package grpc package grpc
package grpc

import (
	"context"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"
	"google.golang.org/grpc"
)

// history -.
type history struct {
	conn *grpc.ClientConn
}

// NewHistory -.
func NewHistory(client *grpc.ClientConn) HistorygRPCAPI {
	return &history{client}
}

// GetStockHistoryTick GetStockHistoryTick
func (t *history) GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	r, err := pb.NewHistoryDataInterfaceClient(t.conn).GetStockHistoryTick(context.Background(), &pb.StockNumArrWithDate{
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
	r, err := pb.NewHistoryDataInterfaceClient(t.conn).GetStockHistoryKbar(context.Background(), &pb.StockNumArrWithDate{
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
	r, err := pb.NewHistoryDataInterfaceClient(t.conn).GetStockHistoryClose(context.Background(), &pb.StockNumArrWithDate{
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
	r, err := pb.NewHistoryDataInterfaceClient(t.conn).GetFutureHistoryKbar(context.Background(), &pb.FutureCodeArrWithDate{
		FutureCodeArr: codeArr,
		Date:          date,
	})
	if err != nil {
		return &pb.HistoryKbarResponse{}, err
	}
	return r, nil
}
