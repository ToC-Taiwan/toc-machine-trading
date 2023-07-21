// Package log package log
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	_defaultNeedCaller = false
	_defaultLogLevel   = LevelInfo
	_defaultLogFormat  = FormatText
	_defaultTimeFormat = "2006-01-02 15:04:05"
	_defaultFilePath   = "logs"
	_defaultFileName   = "log"
)

var singleton *Log

// Log -.
type Log struct {
	*logrus.Logger
	*config

	basePath string
}

type env struct {
	Level          string `env:"LOG_LEVEL"`
	Format         string `env:"LOG_FORMAT"`
	NeedCaller     bool   `env:"LOG_NEED_CALLER"`
	SlackToken     string `env:"SLACK_TOKEN"`
	SlackChannelID string `env:"SLACK_CHANNEL_ID"`
	SlackLogLevel  string `env:"SLACK_LOG_LEVEL"`
}

// Get -.
func Get() *Log {
	if singleton != nil {
		return singleton
	}

	l := &Log{
		Logger: logrus.New(),
		config: &config{
			timeFormat: _defaultTimeFormat,
			format:     _defaultLogFormat,
			level:      _defaultLogLevel,
			needCaller: _defaultNeedCaller,
			fileName:   _defaultFileName,
		},
		basePath: getExecPath(),
	}
	l.readEnv()

	jsonFormatter := &logrus.JSONFormatter{
		DisableHTMLEscape: true,
		TimestampFormat:   l.timeFormat,
		PrettyPrint:       false,
		CallerPrettyfier:  l.setCallerPrettyfier(true),
	}

	textFormatter := &logrus.TextFormatter{
		TimestampFormat:  l.timeFormat,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		PadLevelText:     false,
		ForceColors:      true,
		CallerPrettyfier: l.setCallerPrettyfier(false),
	}

	var formatter logrus.Formatter
	if l.format == FormatJSON {
		formatter = jsonFormatter
	} else {
		formatter = textFormatter
	}

	if l.needCaller {
		l.SetReportCaller(true)
	}

	l.Hooks.Add(NewFilekHook(l.level.Level(), filepath.Join(l.basePath, _defaultFilePath), l.fileName))
	l.SetFormatter(formatter)
	l.SetLevel(l.level.Level())
	l.SetOutput(os.Stdout)

	singleton = l
	return singleton
}

func (l *Log) readEnv() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(filepath.Join(filepath.Dir(ex), ".env"))
	if err != nil {
		panic(err)
	}

	cfg := env{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

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

func (l *Log) setCallerPrettyfier(isJSON bool) func(*runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
		fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", l.basePath), "")
		if isJSON {
			return fmt.Sprintf("%s:%d", fileName, frame.Line), ""
		}
		return fmt.Sprintf("[%s:%d]", fileName, frame.Line), ""
	}
}
