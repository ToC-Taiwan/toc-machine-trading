package websocket

import (
	"context"
	"time"

	"tmt/internal/entity"
)

func (w *WSRouter) sendFuture(ctx context.Context) {
	timestamp := time.Now().UnixNano()
	tickChan := make(chan *entity.RealTimeFutureTick)

	go func() {
		for {
			tick, ok := <-tickChan
			if !ok {
				close(w.msgChan)
				return
			}
			w.msgChan <- tick
		}
	}()

	defer w.s.DeleteFutureRealTimeConnection(timestamp)
	w.s.NewFutureRealTimeConnection(timestamp, tickChan)

	<-ctx.Done()
}
