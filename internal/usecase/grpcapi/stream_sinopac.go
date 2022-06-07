package grpcapi

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pb"
	"toc-machine-trading/pkg/global"
	"toc-machine-trading/pkg/logger"
	"toc-machine-trading/pkg/sinopac"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ActionListMap ActionListMap
var ActionListMap = map[string]int64{
	"Buy":  1,
	"Sell": 2,
}

// StatusListMap StatusListMap
var StatusListMap = map[string]int64{
	"PendingSubmit": 1, // 傳送中
	"PreSubmitted":  2, // 預約單
	"Submitted":     3, // 傳送成功
	"Failed":        4, // 失敗
	"Cancelled":     5, // 已刪除
	"Filled":        6, // 完全成交
	"Filling":       7, // 部分成交
}

// StreamgRPCAPI -.
type StreamgRPCAPI struct {
	conn *sinopac.Connection
}

// NewStream -.
func NewStream(client *sinopac.Connection) *StreamgRPCAPI {
	return &StreamgRPCAPI{client}
}

// EventChannel EventChannel
func (t *StreamgRPCAPI) EventChannel(eventChan chan *entity.SinopacEvent) error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewLongConeectionServiceClient(conn)

	var stream pb.LongConeectionService_EventChannelClient
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if stream == nil {
				if stream, err = c.EventChannel(context.Background(), &emptypb.Empty{}); err != nil {
					logger.Get().Error(err)
					time.Sleep(time.Second * 5)
					continue
				}
			} else {
				break
			}
		}
	}()
	wg.Wait()
	for {
		response, err := stream.Recv()
		if err != nil {
			if !errors.Is(io.EOF, err) {
				return err
			}
			continue
		}
		eventChan <- &entity.SinopacEvent{
			Event:     response.GetEvent(),
			EventCode: response.GetEventCode(),
			Info:      response.GetInfo(),
			Response:  response.GetRespCode(),
		}
	}
}

// TickChannel TickChannel
func (t *StreamgRPCAPI) TickChannel(tickChan chan *entity.RealTimeTick) error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewLongConeectionServiceClient(conn)

	var stream pb.LongConeectionService_TickChannelClient
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if stream == nil {
				if stream, err = c.TickChannel(context.Background(), &emptypb.Empty{}); err != nil {
					logger.Get().Error(err)
					time.Sleep(time.Second * 5)
					continue
				}
			} else {
				break
			}
		}
	}()
	wg.Wait()
	for {
		response, err := stream.Recv()
		if err != nil {
			if !errors.Is(io.EOF, err) {
				return err
			}
			continue
		}
		dataTime, err := time.ParseInLocation(global.LongTimeLayout, response.GetDateTime(), time.Local)
		if err != nil {
			return err
		}
		tickChan <- &entity.RealTimeTick{
			StockNum:        response.GetCode(),
			TickTime:        dataTime,
			Open:            response.GetOpen(),
			AvgPrice:        response.GetAvgPrice(),
			Close:           response.GetClose(),
			High:            response.GetHigh(),
			Low:             response.GetLow(),
			Amount:          response.GetAmount(),
			AmountSum:       response.GetTotalAmount(),
			Volume:          response.GetVolume(),
			VolumeSum:       response.GetTotalVolume(),
			TickType:        response.GetTickType(),
			ChgType:         response.GetChgType(),
			PriceChg:        response.GetPriceChg(),
			PctChg:          response.GetPctChg(),
			BidSideTotalVol: response.GetBidSideTotalVol(),
			AskSideTotalVol: response.GetAskSideTotalVol(),
			BidSideTotalCnt: response.GetBidSideTotalCnt(),
			AskSideTotalCnt: response.GetAskSideTotalCnt(),
			Suspend:         response.GetSuspend(),
			Simtrade:        response.GetSimtrade(),
		}
	}
}

// BidAskChannel BidAskChannel
func (t *StreamgRPCAPI) BidAskChannel(bidAskChan chan *entity.RealTimeBidAsk) error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewLongConeectionServiceClient(conn)

	var stream pb.LongConeectionService_BidAskChannelClient
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if stream == nil {
				if stream, err = c.BidAskChannel(context.Background(), &emptypb.Empty{}); err != nil {
					logger.Get().Error(err)
					time.Sleep(time.Second * 5)
					continue
				}
			} else {
				break
			}
		}
	}()
	wg.Wait()
	for {
		response, err := stream.Recv()
		if err != nil {
			if !errors.Is(io.EOF, err) {
				return err
			}
			continue
		}
		dataTime, err := time.ParseInLocation(global.LongTimeLayout, response.GetDateTime(), time.Local)
		if err != nil {
			return nil
		}
		bidAskChan <- &entity.RealTimeBidAsk{
			StockNum:  response.GetCode(),
			TickTime:  dataTime,
			BidPrice1: response.GetBidPrice()[0], BidVolume1: response.GetBidVolume()[0], DiffBidVol1: response.GetDiffBidVol()[0],
			BidPrice2: response.GetBidPrice()[1], BidVolume2: response.GetBidVolume()[1], DiffBidVol2: response.GetDiffBidVol()[1],
			BidPrice3: response.GetBidPrice()[2], BidVolume3: response.GetBidVolume()[2], DiffBidVol3: response.GetDiffBidVol()[2],
			BidPrice4: response.GetBidPrice()[3], BidVolume4: response.GetBidVolume()[3], DiffBidVol4: response.GetDiffBidVol()[3],
			BidPrice5: response.GetBidPrice()[4], BidVolume5: response.GetBidVolume()[4], DiffBidVol5: response.GetDiffBidVol()[4],
			AskPrice1: response.GetAskPrice()[0], AskVolume1: response.GetAskVolume()[0], DiffAskVol1: response.GetDiffAskVol()[0],
			AskPrice2: response.GetAskPrice()[1], AskVolume2: response.GetAskVolume()[1], DiffAskVol2: response.GetDiffAskVol()[1],
			AskPrice3: response.GetAskPrice()[2], AskVolume3: response.GetAskVolume()[2], DiffAskVol3: response.GetDiffAskVol()[2],
			AskPrice4: response.GetAskPrice()[3], AskVolume4: response.GetAskVolume()[3], DiffAskVol4: response.GetDiffAskVol()[3],
			AskPrice5: response.GetAskPrice()[4], AskVolume5: response.GetAskVolume()[4], DiffAskVol5: response.GetDiffAskVol()[4],
			Suspend:  response.GetSuspend(),
			Simtrade: response.GetSimtrade(),
		}
	}
}

// OrderStatusChannel OrderStatusChannel
func (t *StreamgRPCAPI) OrderStatusChannel(orderStatusChan chan *entity.OrderStatus) error {
	conn := t.conn.GetReadyConn()
	defer t.conn.PutReadyConn(conn)
	c := pb.NewLongConeectionServiceClient(conn)

	var stream pb.LongConeectionService_OrderStatusChannelClient
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if stream == nil {
				if stream, err = c.OrderStatusChannel(context.Background(), &emptypb.Empty{}); err != nil {
					logger.Get().Error(err)
					time.Sleep(time.Second * 5)
					continue
				}
			} else {
				break
			}
		}
	}()
	wg.Wait()
	for {
		response, err := stream.Recv()
		if err != nil {
			if !errors.Is(io.EOF, err) {
				return err
			}
			continue
		}
		actionMap := ActionListMap
		statusMap := StatusListMap
		orderTime, err := time.ParseInLocation(global.LongTimeLayout, response.GetOrderTime(), time.Local)
		if err != nil {
			return err
		}
		orderStatusChan <- &entity.OrderStatus{
			StockNum:  response.GetCode(),
			Action:    actionMap[response.GetAction()],
			Price:     response.GetPrice(),
			Quantity:  response.GetQuantity(),
			Status:    statusMap[response.GetStatus()],
			OrderID:   response.GetOrderId(),
			OrderTime: orderTime,
		}
	}
}
