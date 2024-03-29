// Package slack package slack
package slack

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	slackEmojiPanic string = ":rotating_light:"
	slackEmojiFatal string = ":skull_and_crossbones:"
	slackEmojiError string = ":no_entry:"
	slackEmojiWarn  string = ":warning:"
	slackEmojiInfo  string = ":information_source:"
	slackEmojiDebug string = ":mag_right:"
	slackEmojiTrace string = ":microscope:"
)

var (
	singlton *slackHook
	once     sync.Once
)

type slackHook struct {
	api       *slack.Client
	channelID string

	levels []logrus.Level

	msgChan chan string
	msgLock sync.Mutex
}

type slackSetting struct {
	SlackToken     string `env:"SLACK_TOKEN"`
	SlackChannelID string `env:"SLACK_CHANNEL_ID"`
	SlackLogLevel  string `env:"SLACK_LOG_LEVEL"`
}

func Get() *slackHook {
	if singlton == nil {
		once.Do(func() {
			setting := slackSetting{}
			if err := cleanenv.ReadEnv(&setting); err != nil {
				panic(err)
			}
			hook := &slackHook{
				api: slack.New(
					setting.SlackToken,
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
				channelID: setting.SlackChannelID,
				msgChan:   make(chan string),
			}
			level, err := logrus.ParseLevel(setting.SlackLogLevel)
			if err != nil {
				level = logrus.WarnLevel
			}
			for _, l := range logrus.AllLevels {
				if l <= level {
					hook.levels = append(hook.levels, l)
				}
			}
			go hook.postMessage()
			singlton = hook
		})
		return Get()
	}
	return singlton
}

func (s *slackHook) postMessage() {
	for {
		message := <-s.msgChan
		go func() {
			defer s.msgLock.Unlock()
			s.msgLock.Lock()
			_, _, e := s.api.PostMessage(s.channelID, slack.MsgOptionText(message, false))
			if e != nil {
				fmt.Printf("SlackHook error: %s", e.Error())
				return
			}
		}()
	}
}

func (s *slackHook) Levels() []logrus.Level {
	return s.levels
}

func (s *slackHook) Fire(entry *logrus.Entry) error {
	msg, err := s.Format(entry)
	if err != nil {
		return err
	}

	s.msgChan <- string(msg)
	return nil
}

func (s *slackHook) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var slackEmoji string
	switch entry.Level {
	case logrus.PanicLevel:
		slackEmoji = slackEmojiPanic
	case logrus.FatalLevel:
		slackEmoji = slackEmojiFatal
	case logrus.ErrorLevel:
		slackEmoji = slackEmojiError
	case logrus.WarnLevel:
		slackEmoji = slackEmojiWarn
	case logrus.InfoLevel:
		slackEmoji = slackEmojiInfo
	case logrus.DebugLevel:
		slackEmoji = slackEmojiDebug
	case logrus.TraceLevel:
		slackEmoji = slackEmojiTrace
	}

	// levelText := strings.ToUpper(entry.Level.String())[0:4]
	_, e := b.WriteString(fmt.Sprintf("%s %s", slackEmoji, entry.Message))
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}
