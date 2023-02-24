// Package slack package slack
package slack

import (
	"crypto/tls"
	"net/http"

	"tmt/pkg/log"

	"github.com/slack-go/slack"
)

var logger = log.Get()

type Slack struct {
	api       *slack.Client
	channelID string
}

func NewSlack(token, channelID string) *Slack {
	return &Slack{
		api: slack.New(token, slack.OptionHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		})),
		channelID: channelID,
	}
}

func (s *Slack) PostMessage(message string) {
	_, _, e := s.api.PostMessage(s.channelID, slack.MsgOptionText(message, false))
	if e != nil {
		logger.Errorf("PostMessage to slack error: %v", e)
	}
}
