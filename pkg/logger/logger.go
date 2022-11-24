// Package logger package logger
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Format string

const (
	// LogFormatJSON -.
	LogFormatJSON Format = "json"
	// LogFormatText -.
	LogFormatText Format = "text"
)

type Level string

const (
	LogLevelPanic Level = "panic"
	LogLevelFatal Level = "fatal"
	LogLevelError Level = "error"
	LogLevelWarn  Level = "warn"
	LogLevelInfo  Level = "info"
	LogLevelDebug Level = "debug"
	LogLevelTrace Level = "trace"
)

func (l Level) Level() logrus.Level {
	switch l {
	case LogLevelPanic:
		return logrus.PanicLevel
	case LogLevelFatal:
		return logrus.FatalLevel
	case LogLevelError:
		return logrus.ErrorLevel
	case LogLevelWarn:
		return logrus.WarnLevel
	case LogLevelInfo:
		return logrus.InfoLevel
	case LogLevelDebug:
		return logrus.DebugLevel
	case LogLevelTrace:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

const (
	_defaultLogFormat  = LogFormatText
	_defaultLogLevel   = LogLevelInfo
	_defaultNeedCaller = false
	_defaultTimeFormat = "2006-01-02 15:04:05"
)

// Logger -.
type Logger struct {
	*logrus.Logger
}

type loggerConfig struct {
	timeFormat string
	format     Format
	level      Level
	needCaller bool
}

func newLoggerConfig(opts ...Option) *loggerConfig {
	cfg := &loggerConfig{
		format:     _defaultLogFormat,
		level:      _defaultLogLevel,
		needCaller: _defaultNeedCaller,
		timeFormat: _defaultTimeFormat,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// NewLogger -.
func NewLogger(opts ...Option) *Logger {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Clean(filepath.Dir(ex))
	cfg := newLoggerConfig(opts...)

	l := logrus.New()
	jsonFormatter := &logrus.JSONFormatter{
		DisableHTMLEscape: true,
		TimestampFormat:   cfg.timeFormat,
		PrettyPrint:       false,
		CallerPrettyfier:  newCallerPrettyfier(basePath, true),
	}

	textFormatter := &logrus.TextFormatter{
		TimestampFormat:  cfg.timeFormat,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		PadLevelText:     false,
		ForceColors:      true,
		CallerPrettyfier: newCallerPrettyfier(basePath, false),
	}

	var formatter logrus.Formatter
	if cfg.format == LogFormatJSON {
		formatter = jsonFormatter
	} else {
		formatter = textFormatter
	}

	if cfg.needCaller {
		l.SetReportCaller(true)
	}

	l.Hooks.Add(fileHook(basePath, jsonFormatter, cfg.timeFormat))
	l.SetFormatter(formatter)
	l.SetLevel(cfg.level.Level())
	l.SetOutput(os.Stdout)

	return &Logger{
		Logger: l,
	}
}

func fileHook(basePath string, formatter logrus.Formatter, timeFormat string) *lfshook.LfsHook {
	date := strings.ReplaceAll(time.Now().Format(timeFormat), ":", "")
	return lfshook.NewHook(
		lfshook.PathMap{
			logrus.PanicLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
			logrus.FatalLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
			logrus.ErrorLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
			logrus.WarnLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
			logrus.InfoLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		},
		formatter,
	)
}

func newCallerPrettyfier(basePath string, isJSON bool) func(*runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
		fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
		if isJSON {
			return fmt.Sprintf("%s:%d", fileName, frame.Line), ""
		}
		return fmt.Sprintf("[%s:%d]", fileName, frame.Line), ""
	}
}
