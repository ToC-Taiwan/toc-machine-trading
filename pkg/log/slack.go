package log

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackHook struct {
	api       *slack.Client
	channelID string

	levels []logrus.Level

	msgChan  chan string
	hookLock sync.Mutex
}

func NewSlackHook(token, channelID string, level logrus.Level) *SlackHook {
	hook := &SlackHook{
		api: slack.New(
			token,
			slack.OptionHTTPClient(
				&http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				},
			),
		),
		channelID: channelID,
		msgChan:   make(chan string),
	}

	for _, l := range logrus.AllLevels {
		if l <= level {
			hook.levels = append(hook.levels, l)
		}
	}

	go hook.PostMessage()
	return hook
}

func (s *SlackHook) PostMessage() {
	for {
		message := <-s.msgChan
		go func() {
			defer s.hookLock.Unlock()
			s.hookLock.Lock()
			_, _, e := s.api.PostMessage(s.channelID, slack.MsgOptionText(message, false))
			if e != nil {
				fmt.Printf("SlackHook error: %s", e.Error())
				return
			}
		}()
	}
}

func (s *SlackHook) Levels() []logrus.Level {
	return s.levels
}

func (s *SlackHook) Fire(entry *logrus.Entry) error {
	msg, err := s.Format(entry)
	if err != nil {
		return err
	}

	s.msgChan <- string(msg)
	return nil
}

func (s *SlackHook) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	_, e := b.WriteString(fmt.Sprintf("[%s][%s]  %-30s", levelText, entry.Time.Format(_defaultTimeFormat), strings.TrimSuffix(entry.Message, "\n")))
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}
