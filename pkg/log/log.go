// Package log package log
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Log struct {
	*logrus.Logger
	*config
}

type env struct {
	Level      string `env:"LOG_LEVEL"`
	Format     string `env:"LOG_FORMAT"`
	NeedCaller bool   `env:"LOG_NEED_CALLER"`
	TimeFormat string `env:"LOG_TIME_FORMAT"`

	SlackToken     string `env:"SLACK_TOKEN"`
	SlackChannelID string `env:"SLACK_CHANNEL_ID"`
	SlackLogLevel  string `env:"SLACK_LOG_LEVEL"`
}

var (
	singleton *Log
	initOnce  sync.Once
)

// Get -.
func Get() *Log {
	if singleton != nil {
		return singleton
	}

	initOnce.Do(func() {
		l := &Log{
			Logger: logrus.New(),
			config: &config{
				timeFormat: _defaultTimeFormat,
				format:     _defaultLogFormat,
				level:      _defaultLogLevel,
				needCaller: _defaultNeedCaller,
				fileName:   _defaultFileName,
			},
		}
		l.readEnv()

		var formatter logrus.Formatter
		if l.format == FormatJSON {
			formatter = &logrus.JSONFormatter{
				DisableHTMLEscape: true,
				TimestampFormat:   l.timeFormat,
				PrettyPrint:       false,
				CallerPrettyfier:  l.callerPrettyfier(true),
			}
		} else {
			formatter = &logrus.TextFormatter{
				TimestampFormat:  l.timeFormat,
				FullTimestamp:    true,
				QuoteEmptyFields: true,
				PadLevelText:     false,
				ForceColors:      true,
				CallerPrettyfier: l.callerPrettyfier(false),
			}
		}

		if l.needCaller {
			l.SetReportCaller(true)
		}

		l.Hooks.Add(
			NewFilekHook(
				l.level.Level(),
				filepath.Join(
					getExecPath(),
					_defaultFilePath,
				),
				l.fileName,
			),
		)
		l.SetFormatter(formatter)
		l.SetLevel(l.level.Level())
		l.SetOutput(os.Stdout)
		l.Info("Logger initialized successfully")

		singleton = l
	})

	return singleton
}

func (l *Log) readEnv() {
	cfg := env{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	TimeFormat(cfg.TimeFormat)(l.config)
	LogLevel(cfg.Level)(l.config)
	LogFormat(cfg.Format)(l.config)
	NeedCaller(cfg.NeedCaller)(l.config)

	var slackLevel logrus.Level
	switch cfg.SlackLogLevel {
	case LevelPanic.String(), LevelFatal.String(), LevelError.String(), LevelWarn.String(), LevelInfo.String(), LevelDebug.String(), LevelTrace.String():
		slackLevel = Level(cfg.SlackLogLevel).Level()
	default:
		slackLevel = logrus.WarnLevel
	}

	if cfg.SlackToken != "" && cfg.SlackChannelID != "" {
		l.Hooks.Add(
			NewSlackHook(
				cfg.SlackToken,
				cfg.SlackChannelID,
				slackLevel,
			),
		)
	}
}

func (l *Log) callerPrettyfier(isJSON bool) func(*runtime.Frame) (function string, file string) {
	basePath := getExecPath()
	return func(frame *runtime.Frame) (function string, file string) {
		fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
		if isJSON {
			return fmt.Sprintf("%s:%d", fileName, frame.Line), ""
		}
		return fmt.Sprintf("[%s:%d]", fileName, frame.Line), ""
	}
}
