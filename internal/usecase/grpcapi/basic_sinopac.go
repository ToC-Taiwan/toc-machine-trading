// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BasicgRPCAPI -.
type BasicgRPCAPI struct {
	conn *sinopac.Connection
}

// NewBasic -.
func NewBasic(client *sinopac.Connection) *BasicgRPCAPI {
	return &BasicgRPCAPI{client}
}

// GetAllStockDetail GetAllStockDetail
func (t *BasicgRPCAPI) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// // HealthCheck HealthCheck
// func (s *Connection) HealthCheck() (string, error) {
// 	conn := s.getReadyConn()
// 	defer s.putReadyConn(conn)
// 	c := pb.NewSinopacForwarderClient(conn)
// 	r, err := c.GetServerToken(context.Background(), &emptypb.Empty{})
// 	if err != nil {
// 		return "", err
// 	}
// 	return r.GetToken(), nil
// }

// // GetAllSnapshot GetAllSnapshot
// func GetAllSnapshot() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetAllStockSnapshot(ctx, &emptypb.Empty{})
// 	if err != nil {
// 		fmt.Printf("GetAllSnapshot: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// // GetAllSnapshotTSE GetAllSnapshotTSE
// func GetAllSnapshotTSE() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockSnapshotTSE(ctx, &emptypb.Empty{})
// 	if err != nil {
// 		fmt.Printf("GetAllSnapshotTSE: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// // GetStockSnapshotByStockNumArr GetStockSnapshotByStockNumArr
// func GetStockSnapshotByStockNumArr() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockSnapshotByNumArr(ctx, &pb.StockNumArr{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		// StockNumArr: []string{"2890", "2884", "1605"},
// 	})
// 	if err != nil {
// 		fmt.Printf("GetStockSnapshotByStockNumArr: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// // GetStockTSEHistoryTick GetStockTSEHistoryTick
// func GetStockTSEHistoryTick() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockTSEHistoryTick(ctx, &pb.Date{
// 		Date: "2022-05-26",
// 	})
// 	if err != nil {
// 		fmt.Printf("GetStockTSEHistoryTick: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// 	fmt.Println(len(r.GetData()))
// }

// // GetStockHistoryTick GetStockHistoryTick
// func GetStockHistoryTick() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	requestTime := time.Now().UnixNano()
// 	r, err := c.GetStockHistoryTick(ctx, &pb.StockNumArrWithDate{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		Date:        "2022-05-23",
// 	})
// 	if err != nil {
// 		fmt.Printf("GetStockHistoryTick: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// 	fmt.Println(len(r.GetData()))
// 	fmt.Println((time.Now().UnixNano() - requestTime) / 1000 / 1000 / 1000)
// }

// // GetStockHistoryKbar GetStockHistoryKbar
// func GetStockHistoryKbar() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockHistoryKbar(ctx, &pb.StockNumArrWithDate{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		Date:        "2022-05-23",
// 	})
// 	if err != nil {
// 		fmt.Printf("GetStockHistoryKbar: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// // GetStockHistoryClose GetStockHistoryClose
// func GetStockHistoryClose() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockHistoryClose(ctx, &pb.StockNumArrWithDate{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		Date:        "2022-05-23",
// 	})
// 	if err != nil {
// 		fmt.Printf("GetStockHistoryClose: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// // GetStockVolumeRank GetStockVolumeRank
// func GetStockVolumeRank() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockVolumeRank(ctx, &pb.VolumeRankRequest{
// 		Count: 200,
// 		Date:  "2022-05-25",
// 	})
// 	if err != nil {
// 		fmt.Printf("GetStockVolumeRank: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// // SubscribeStockTick SubscribeStockTick
// func SubscribeStockTick() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.SubscribeStockTick(ctx, &pb.StockNumArr{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		// StockNumArr: []string{"2890", "2884", "1605"},
// 	})
// 	if err != nil {
// 		fmt.Printf("SubscribeStockTick: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetFailArr() {
// 		fmt.Println(v)
// 	}
// }

// // UnSubscribeStockTick UnSubscribeStockTick
// func UnSubscribeStockTick() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.UnSubscribeStockTick(ctx, &pb.StockNumArr{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		// StockNumArr: []string{"2890", "2884", "1605"},
// 	})
// 	if err != nil {
// 		fmt.Printf("UnSubscribeStockTick: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetFailArr() {
// 		fmt.Println(v)
// 	}
// }

// // SubscribeStockBidAsk SubscribeStockBidAsk
// func SubscribeStockBidAsk() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.SubscribeStockBidAsk(ctx, &pb.StockNumArr{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		// StockNumArr: []string{"2890", "2884", "1605"},
// 	})
// 	if err != nil {
// 		fmt.Printf("SubscribeStockBidAsk: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetFailArr() {
// 		fmt.Println(v)
// 	}
// }

// // UnSubscribeStockBidAsk UnSubscribeStockBidAsk
// func UnSubscribeStockBidAsk() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.UnSubscribeStockBidAsk(ctx, &pb.StockNumArr{
// 		StockNumArr: []string{"2330", "2610", "2303", "1708", "2603"},
// 		// StockNumArr: []string{"2890", "2884", "1605"},
// 	})
// 	if err != nil {
// 		fmt.Printf("UnSubscribeStockBidAsk: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetFailArr() {
// 		fmt.Println(v)
// 	}
// }

// // GetStockTSEHistoryClose GetStockTSEHistoryClose
// func GetStockTSEHistoryClose() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewSinopacForwarderClient(conn)
// 	ctx := context.Background()

// 	r, err := c.GetStockTSEHistoryClose(ctx, &pb.Date{
// 		Date: "2022-05-23",
// 	})
// 	if err != nil {
// 		fmt.Printf("UnSubscribeStockBidAsk: %v", err)
// 		return
// 	}
// 	for _, v := range r.GetData() {
// 		fmt.Println(v)
// 	}
// }

// stuck := make(chan struct{})
// sinopac.HealthCheck()
// logger.New("INOF").Warn("message string")
// go sinopac.EventChannel()
// go sinopac.TickChannel()
// go sinopac.BidAskChannel()
// // time.Sleep(time.Second * 3)
// sinopac.SubscribeStockTick()
// sinopac.UnSubscribeStockTick()
// sinopac.SubscribeStockBidAsk()
// sinopac.UnSubscribeStockBidAsk()
// sinopac.GetAllStockDetail()
// sinopac.GetAllSnapshot()
// sinopac.GetStockSnapshotByStockNumArr()
// sinopac.GetAllSnapshotTSE()
// sinopac.GetStockHistoryTick()
// sinopac.GetStockHistoryKbar()
// sinopac.GetStockHistoryClose()
// sinopac.GetStockVolumeRank()
// sinopac.GetStockTSEHistoryTick()
// sinopac.GetStockTSEHistoryClose()
// <-stuck
