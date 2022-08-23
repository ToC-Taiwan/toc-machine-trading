package grpcapi

import (
	"context"

	"tmt/pb"
	"tmt/pkg/sinopac"
)

// HistorygRPCAPI -.
type HistorygRPCAPI struct {
	conn *sinopac.Connection
}

// NewHistory -.
func NewHistory(client *sinopac.Connection) *HistorygRPCAPI {
	return &HistorygRPCAPI{client}
}

// GetStockHistoryTick GetStockHistoryTick
func (t *HistorygRPCAPI) GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.StockHistoryTickMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockHistoryTick(context.Background(), &pb.StockNumArrWithDate{
		StockNumArr: stockNumArr,
		Date:        date,
	})
	if err != nil {
		return []*pb.StockHistoryTickMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockHistoryKbar GetStockHistoryKbar
func (t *HistorygRPCAPI) GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.StockHistoryKbarMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockHistoryKbar(context.Background(), &pb.StockNumArrWithDate{
		StockNumArr: stockNumArr,
		Date:        date,
	})
	if err != nil {
		return []*pb.StockHistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockHistoryClose GetStockHistoryClose
func (t *HistorygRPCAPI) GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.StockHistoryCloseMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockHistoryClose(context.Background(), &pb.StockNumArrWithDate{
		StockNumArr: stockNumArr,
		Date:        date,
	})
	if err != nil {
		return []*pb.StockHistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockHistoryCloseByDateArr GetStockHistoryCloseByDateArr
func (t *HistorygRPCAPI) GetStockHistoryCloseByDateArr(stockNumArr []string, date []string) ([]*pb.StockHistoryCloseMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockHistoryCloseByDateArr(context.Background(), &pb.StockNumArrWithDateArr{
		StockNumArr: stockNumArr,
		DateArr:     date,
	})
	if err != nil {
		return []*pb.StockHistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryTick GetStockTSEHistoryTick
func (t *HistorygRPCAPI) GetStockTSEHistoryTick(date string) ([]*pb.StockHistoryTickMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockTSEHistoryTick(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.StockHistoryTickMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryKbar GetStockTSEHistoryKbar
func (t *HistorygRPCAPI) GetStockTSEHistoryKbar(date string) ([]*pb.StockHistoryKbarMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockTSEHistoryKbar(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.StockHistoryKbarMessage{}, err
	}
	return r.GetData(), nil
}

// GetStockTSEHistoryClose GetStockTSEHistoryClose
func (t *HistorygRPCAPI) GetStockTSEHistoryClose(date string) ([]*pb.StockHistoryCloseMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetStockTSEHistoryClose(context.Background(), &pb.Date{
		Date: date,
	})
	if err != nil {
		return []*pb.StockHistoryCloseMessage{}, err
	}
	return r.GetData(), nil
}
