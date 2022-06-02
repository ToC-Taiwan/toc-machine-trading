// Package grpcapi package grpcapi
package grpcapi

import (
	"context"

	"toc-machine-trading/pkg/logger"
	"toc-machine-trading/pkg/pb"
	"toc-machine-trading/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// StockgRPCAPI -.
type StockgRPCAPI struct {
	conn *sinopac.Connection
}

// New -.
func New(url string, poolSize int) *StockgRPCAPI {
	client, err := sinopac.New(url, sinopac.MaxPoolSize(poolSize))
	if err != nil {
		logger.Get().Panic(err)
	}

	return &StockgRPCAPI{
		conn: client,
	}
}

// GetAllStockDetail GetAllStockDetail
func (t *StockgRPCAPI) GetAllStockDetail() ([]*pb.StockDetailMessage, error) {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewSinopacForwarderClient(conn)
	r, err := c.GetAllStockDetail(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []*pb.StockDetailMessage{}, err
	}
	return r.GetStock(), nil
}

// func connectSinopac() {
// 	// gRPC
// 	client, err := sinopac.New("127.0.0.1:56666")
// 	if err != nil {
// 		logger.Get().Panic(err)
// 	}
// 	// go client.EventChannel()
// 	// logger.Get().Warn(client.HealthCheck())
// 	// client.GetAllStockDetail()
// }

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

// // EventChannel EventChannel
// func (s *Connection) EventChannel() {
// 	conn := s.getReadyConn()
// 	defer s.putReadyConn(conn)
// 	c := pb.NewLongConeectionServiceClient(conn)
// 	var stream pb.LongConeectionService_EventChannelClient
// 	var err error
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for {
// 			if stream == nil {
// 				if stream, err = c.EventChannel(context.Background(), &emptypb.Empty{}); err != nil {
// 					logger.Get().Error(err)
// 					time.Sleep(time.Second * 5)
// 					continue
// 				}
// 			} else {
// 				break
// 			}
// 		}
// 	}()
// 	wg.Wait()
// 	for {
// 		response, err := stream.Recv()
// 		if err != nil {
// 			if !errors.Is(io.EOF, err) {
// 				panic(err)
// 			}
// 			continue
// 		}
// 		logger.Get().Infof("%d %d %s %s\n", response.GetRespCode(), response.GetEventCode(), response.GetInfo(), response.GetEvent())
// 	}
// }

// var wg sync.WaitGroup

// // BidAskChannel BidAskChannel
// func BidAskChannel() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewLongConeectionServiceClient(conn)
// 	var stream pb.LongConeectionService_BidAskChannelClient
// 	ctx := context.Background()
// 	var err error
// 	for {
// 		if stream == nil {
// 			if stream, err = c.BidAskChannel(ctx, &emptypb.Empty{}); err != nil {
// 				fmt.Println(err)
// 				time.Sleep(time.Second * 5)
// 				continue
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	wg.Done()
// 	for {
// 		response, err := stream.Recv()
// 		if err != nil {
// 			if !errors.Is(io.EOF, err) {
// 				panic(err)
// 			}
// 			continue
// 		}
// 		fmt.Println(response)
// 	}
// }

// // TickChannel TickChannel
// func TickChannel() {
// 	conn := newConnection()
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	c := pb.NewLongConeectionServiceClient(conn)
// 	var stream pb.LongConeectionService_TickChannelClient
// 	ctx := context.Background()
// 	var err error
// 	for {
// 		if stream == nil {
// 			if stream, err = c.TickChannel(ctx, &emptypb.Empty{}); err != nil {
// 				fmt.Println(err)
// 				time.Sleep(time.Second * 5)
// 				continue
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	wg.Done()
// 	for {
// 		response, err := stream.Recv()
// 		if err != nil {
// 			if !errors.Is(io.EOF, err) {
// 				panic(err)
// 			}
// 			continue
// 		}
// 		fmt.Println(response)
// 	}
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
