// Package log package log
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

const (
	_defaultNeedCaller = false
	_defaultLogLevel   = logrus.InfoLevel
	_defaultLogFormat  = formatText
	_defaultTimeFormat = time.RFC3339
	_defaultFileName   = "app"
)

const (
	formatJSON string = "json"
	formatText string = "text"
)

type Log struct {
	*logrus.Logger
	*config
}

type config struct {
	NeedCaller bool   `env:"LOG_NEED_CALLER"`
	Format     string `env:"LOG_FORMAT"`
	TimeFormat string `env:"LOG_TIME_FORMAT"`
	Level      string `env:"LOG_LEVEL"`
	FileName   string `env:"LOG_FILE_NAME"`
	LinkSlack  bool   `env:"LOG_LINK_SLACK"`
}

var (
	global *Log
	once   sync.Once
)

func Get() *Log {
	if global != nil {
		return global
	}
	once.Do(func() {
		global = newLogger()
	})
	return global
}

func newLogger() *Log {
	l := new(Log)
	l.Logger = logrus.New()

	l.readConfig()
	l.setFormatter()
	l.setFileHook()
	l.setSlackHook()
	return l
}

func (l *Log) readConfig() {
	cfg := config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}
	l.config = &cfg

	l.Logger.SetReportCaller(cfg.NeedCaller)

	if cfg.Format != formatJSON && cfg.Format != formatText {
		cfg.Format = _defaultLogFormat
	}

	_, err := time.Parse(cfg.TimeFormat, "2006-01-02")
	if err != nil {
		cfg.TimeFormat = _defaultTimeFormat
	}

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		cfg.Level = _defaultLogLevel.String()
	}
	l.Logger.SetLevel(level)

	if cfg.FileName == "" {
		cfg.FileName = _defaultFileName
	}
}

func (l *Log) setFormatter() {
	var formatter logrus.Formatter
	if l.config.Format == formatJSON {
		formatter = &logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   l.config.TimeFormat,
			PrettyPrint:       false,
			CallerPrettyfier:  l.callerPrettyfier(),
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat:  l.config.TimeFormat,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
			CallerPrettyfier: l.callerPrettyfier(),
		}
	}
	l.Logger.SetFormatter(formatter)
}

func (l *Log) setFileHook() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Clean(filepath.Dir(ex))
	fileHook := NewFilekHook(
		l.Logger.Level,
		filepath.Join(
			basePath,
			"logs",
		),
		l.config.FileName,
	)
	if l.config.NeedCaller {
		fileHook.SetReportCaller(true, l.callerPrettyfier())
	}
	l.Logger.AddHook(fileHook)
}

func (l *Log) setSlackHook() {
	if !l.config.LinkSlack {
		return
	}
	l.Logger.AddHook(NewSlackHook())
}

func (l *Log) callerPrettyfier() func(*runtime.Frame) (string, string) {
	return func(frame *runtime.Frame) (string, string) {
		path := frame.File
		caller := []string{}
		for filepath.Base(filepath.Dir(path)) != "toc-machine-trading" {
			caller = append(caller, filepath.Base(path))
			path = filepath.Dir(path)
		}
		for i, j := 0, len(caller)-1; i < j; i, j = i+1, j-1 {
			caller[i], caller[j] = caller[j], caller[i]
		}
		if l.config.Format == formatJSON {
			return fmt.Sprintf("%s:%d", filepath.Join(caller...), frame.Line), ""
		}
		return fmt.Sprintf("[%s:%d]", filepath.Join(caller...), frame.Line), ""
	}
}
