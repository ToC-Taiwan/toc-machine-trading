package usecase

import (
	"context"
	"fmt"

	"tmt/internal/entity"
	"tmt/internal/usecase/modules/calendar"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FcmUseCase struct {
	app      *firebase.App
	logger   *log.Log
	bus      *eventbus.Bus
	tradeDay *calendar.Calendar
}

// NewFCM -.
func NewFCM() FCM {
	logger := log.Get()
	fb, err := newFCM()
	if err != nil {
		logger.Fatal(err)
	}

	uc := &FcmUseCase{
		app:      fb,
		logger:   logger,
		bus:      eventbus.Get(),
		tradeDay: calendar.Get(),
	}

	uc.bus.SubscribeAsync(topicFetchStockHistory, true, uc.sendTargets)
	return uc
}

func newFCM() (*firebase.App, error) {
	opt := option.WithCredentialsFile("configs/service_account.json")
	fb, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

func (uc *FcmUseCase) sendTargets(targetArr []*entity.StockTarget) error {
	ctx := context.Background()
	client, err := uc.app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Found New Targets",
			Body:  fmt.Sprintf("%s has %d targets", uc.tradeDay.GetStockTradeDay().TradeDay.Format(entity.ShortTimeLayout), len(targetArr)),
		},
		Data: map[string]string{
			"new_targets_count": fmt.Sprintf("%d", len(targetArr)),
		},
		Topic: "new_targets",
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (uc *FcmUseCase) AnnounceMessage(msg string) error {
	ctx := context.Background()
	client, err := uc.app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Announcement",
			Body:  msg,
		},
		Topic: "announcement",
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
