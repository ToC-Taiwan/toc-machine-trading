package future

import (
	"tmt/internal/entity"
	"tmt/pb"
	"tmt/pkg/common"
)

func newFutureTickProto(r *entity.RealTimeFutureTick) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_FUTURE_TICK,
		Data: &pb.WSMessage_FutureTick{
			FutureTick: &pb.WSFutureTick{
				Code:            r.Code,
				TickTime:        r.TickTime.Format(common.LongTimeLayout),
				Open:            r.Open,
				UnderlyingPrice: r.UnderlyingPrice,
				BidSideTotalVol: r.BidSideTotalVol,
				AskSideTotalVol: r.AskSideTotalVol,
				AvgPrice:        r.AvgPrice,
				Close:           r.Close,
				High:            r.High,
				Low:             r.Low,
				Amount:          r.Amount,
				TotalAmount:     r.TotalAmount,
				Volume:          r.Volume,
				TotalVolume:     r.TotalVolume,
				TickType:        r.TickType,
				ChgType:         r.ChgType,
				PriceChg:        r.PriceChg,
				PctChg:          r.PctChg,
			},
		},
	}
}

func newFutureOrderProto(f *entity.FutureOrder) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_FUTURE_ORDER,
		Data: &pb.WSMessage_FutureOrder{
			FutureOrder: &pb.WSFutureOrder{
				Code: f.Code,
				BaseOrder: &pb.WSOrder{
					OrderId:   f.OrderID,
					Status:    int64(f.Status),
					OrderTime: f.OrderTime.Format(common.LongTimeLayout),
					Action:    int64(f.Action),
					Price:     f.Price,
					Quantity:  f.Quantity,
					TradeTime: f.TradeTime.Format(common.LongTimeLayout),
					TickTime:  f.TickTime.Format(common.LongTimeLayout),
					GroupId:   f.GroupID,
				},
			},
		},
	}
}

func newTradeVolumeProto(firstPeriod, secondPeriod, thirdPeriod, fourthPeriod entity.RealTimeFutureTickArr) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_PERIOD_TRADE_VOLUME,
		Data: &pb.WSMessage_PeriodTradeVolume{
			PeriodTradeVolume: &pb.WSPeriodTradeVolume{
				FirstPeriod: &pb.OutInVolume{
					OutVolume: firstPeriod.GetOutInVolume().OutVolume,
					InVolume:  firstPeriod.GetOutInVolume().InVolume,
				},
				SecondPeriod: &pb.OutInVolume{
					OutVolume: secondPeriod.GetOutInVolume().OutVolume,
					InVolume:  secondPeriod.GetOutInVolume().InVolume,
				},
				ThirdPeriod: &pb.OutInVolume{
					OutVolume: thirdPeriod.GetOutInVolume().OutVolume,
					InVolume:  thirdPeriod.GetOutInVolume().InVolume,
				},
				FourthPeriod: &pb.OutInVolume{
					OutVolume: fourthPeriod.GetOutInVolume().OutVolume,
					InVolume:  fourthPeriod.GetOutInVolume().InVolume,
				},
			},
		},
	}
}

func newTradeIndexProto(t *entity.TradeIndex) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_TRADE_INDEX,
		Data: &pb.WSMessage_TradeIndex{
			TradeIndex: &pb.WSTradeIndex{
				Tse:    newStockSnapshotProto(t.TSE),
				Otc:    newStockSnapshotProto(t.OTC),
				Nasdaq: newYahooPriceProto(t.Nasdaq),
				Nf:     newYahooPriceProto(t.NF),
			},
		},
	}
}

func newYahooPriceProto(y *entity.YahooPrice) *pb.WSYahooPrice {
	return &pb.WSYahooPrice{
		Last:      y.Last,
		Price:     y.Price,
		UpdatedAt: y.UpdatedAt.Format(common.LongTimeLayout),
	}
}

func newStockSnapshotProto(s *entity.StockSnapShot) *pb.WSStockSnapShot {
	return &pb.WSStockSnapShot{
		StockNum:        s.StockNum,
		StockName:       s.StockName,
		SnapTime:        s.SnapTime.Format(common.LongTimeLayout),
		Open:            s.Open,
		High:            s.High,
		Low:             s.Low,
		Close:           s.Close,
		TickType:        s.TickType,
		PriceChg:        s.PriceChg,
		PctChg:          s.PctChg,
		ChgType:         s.ChgType,
		Volume:          s.Volume,
		VolumeSum:       s.VolumeSum,
		Amount:          s.Amount,
		AmountSum:       s.AmountSum,
		YesterdayVolume: s.YesterdayVolume,
		VolumeRatio:     s.VolumeRatio,
	}
}

func newFuturePositionProto(position entity.FuturePositionArr) *pb.WSMessage {
	var ret []*pb.Position
	for _, v := range position {
		ret = append(ret, &pb.Position{
			Code:      v.Code,
			Direction: v.Direction,
			Quantity:  v.Quantity,
			Price:     v.Price,
			LastPrice: v.LastPrice,
			Pnl:       v.Pnl,
		})
	}

	return &pb.WSMessage{
		Type: pb.WSType_TYPE_FUTURE_POSITION,
		Data: &pb.WSMessage_FuturePosition{
			FuturePosition: &pb.WSFuturePosition{
				Position: ret,
			},
		},
	}
}

func newAssistStatusProto(status bool) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_ASSIST_STATUS,
		Data: &pb.WSMessage_AssitStatus{
			AssitStatus: &pb.WSAssitStatus{Running: status},
		},
	}
}

func newErrMessageProto(err *futureTradeError) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_ERR_MESSAGE,
		Data: &pb.WSMessage_ErrMessage{
			ErrMessage: &pb.WSErrMessage{
				ErrCode:  int64(err.ErrCode),
				Response: err.Response,
			},
		},
	}
}

func newKbarArrProto(r []*entity.FutureHistoryKbar) *pb.WSMessage {
	var ret []*pb.Kbar
	for _, v := range r {
		ret = append(ret, &pb.Kbar{
			KbarTime: v.KbarTime.Format(common.LongTimeLayout),
			Close:    v.Close,
			Open:     v.Open,
			High:     v.High,
			Low:      v.Low,
			Volume:   v.Volume,
		})
	}

	return &pb.WSMessage{
		Type: pb.WSType_TYPE_KBAR_ARR,
		Data: &pb.WSMessage_HistoryKbar{
			HistoryKbar: &pb.WSHistoryKbarMessage{
				Arr: ret,
			},
		},
	}
}
