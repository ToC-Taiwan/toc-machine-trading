package log

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type slackHook struct {
	api       *slack.Client
	channelID string
}

func newSlackHook(token, channelID string) *slackHook {
	return &slackHook{
		api:       slack.New(token),
		channelID: channelID,
	}
}

func (s *slackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (s *slackHook) Fire(e *logrus.Entry) error {
	_, _, err := s.api.PostMessage(s.channelID, slack.MsgOptionText(e.Message, false))
	if err != nil {
		return err
	}
	return nil
}
