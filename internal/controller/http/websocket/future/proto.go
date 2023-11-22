package future

import (
	"tmt/internal/entity"
	"tmt/pb"
)

func newFutureTickProto(r *entity.RealTimeFutureTick) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_FUTURE_TICK,
		Data: &pb.WSMessage_FutureTick{
			FutureTick: &pb.WSFutureTick{
				Code:            r.Code,
				TickTime:        r.TickTime.Format(entity.LongTimeLayout),
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
					OrderTime: f.OrderTime.Format(entity.LongTimeLayout),
					Action:    int64(f.Action),
					Price:     f.Price,
					Quantity:  f.Quantity,
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
				Tse:    newIndexStatusProto(t.TSE),
				Otc:    newIndexStatusProto(t.OTC),
				Nasdaq: newIndexStatusProto(t.Nasdaq),
				Nf:     newIndexStatusProto(t.NF),
			},
		},
	}
}

func newIndexStatusProto(status *entity.IndexStatus) *pb.WSIndexStatus {
	return &pb.WSIndexStatus{
		BreakCount: status.BreakCount,
		PriceChg:   status.PriceChg,
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
			KbarTime: v.KbarTime.Format(entity.LongTimeLayout),
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

func newFutureDetailProto(r *entity.Future) *pb.WSMessage {
	return &pb.WSMessage{
		Type: pb.WSType_TYPE_FUTURE_DETAIL,
		Data: &pb.WSMessage_FutureDetail{
			FutureDetail: &pb.WSFutureDetail{
				Code:           r.Code,
				Symbol:         r.Symbol,
				Name:           r.Name,
				Category:       r.Category,
				DeliveryMonth:  r.DeliveryMonth,
				DeliveryDate:   r.DeliveryDate.Format(entity.ShortTimeLayout),
				UnderlyingKind: r.UnderlyingKind,
				Unit:           r.Unit,
				LimitUp:        r.LimitUp,
				LimitDown:      r.LimitDown,
				Reference:      r.Reference,
				UpdateDate:     r.UpdateDate.Format(entity.ShortTimeLayout),
			},
		},
	}
}
