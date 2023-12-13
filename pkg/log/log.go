// Package log package log
package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"tmt/pkg/log/hook/file"
	"tmt/pkg/log/hook/slack"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

const (
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
	NeedCaller     bool   `env:"LOG_NEED_CALLER"`
	Format         string `env:"LOG_FORMAT"`
	TimeFormat     string `env:"LOG_TIME_FORMAT"`
	Level          string `env:"LOG_LEVEL"`
	LinkSlack      bool   `env:"LOG_LINK_SLACK"`
	DisableConsole bool   `env:"LOG_DISABLE_CONSOLE"`
	DisableFile    bool   `env:"LOG_DISABLE_FILE"`
	FileName       string `env:"LOG_FILE_NAME"`
}

func Get() *Log {
	l := &Log{
		Logger: logrus.New(),
	}
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

	if cfg.FileName == "" {
		cfg.FileName = _defaultFileName
	}
	if cfg.Format != formatJSON && cfg.Format != formatText {
		cfg.Format = _defaultLogFormat
	}
	_, err := time.Parse(cfg.TimeFormat, time.DateOnly)
	if err != nil {
		cfg.TimeFormat = _defaultTimeFormat
	}
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = _defaultLogLevel
	}
	l.Logger.SetLevel(level)
	if cfg.DisableConsole {
		l.Logger.SetOutput(io.Discard)
	}
	l.Logger.SetReportCaller(cfg.NeedCaller)
	l.config = &cfg
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
			PadLevelText:     false,
			ForceColors:      true,
			CallerPrettyfier: l.callerPrettyfier(),
		}
	}
	l.Logger.SetFormatter(formatter)
}

func (l *Log) setFileHook() {
	if l.config.DisableFile {
		return
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fileHook := file.Get(
		l.Logger.Level,
		filepath.Join(filepath.Clean(filepath.Dir(ex)), "logs"),
		l.config.FileName,
	)
	fileHook.SetReportCaller(l.config.NeedCaller, l.callerPrettyfier())
	l.Logger.AddHook(fileHook)
}

func (l *Log) setSlackHook() {
	if !l.config.LinkSlack {
		return
	}
	l.Logger.AddHook(slack.Get())
}

func (l *Log) callerPrettyfier() func(*runtime.Frame) (string, string) {
	if !l.config.NeedCaller {
		return nil
	}
	return func(frame *runtime.Frame) (string, string) {
		path := frame.File
		if path == "" {
			return "", ""
		}
		split := strings.Split(path, "/")
		if len(split) < 2 {
			return "", ""
		}
		if l.config.Format == formatJSON {
			return fmt.Sprintf("%s:%d", split[len(split)-1], frame.Line), ""
		}
		return fmt.Sprintf("[%s:%d]", split[len(split)-1], frame.Line), ""
	}
}
