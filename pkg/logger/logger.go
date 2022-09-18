// Package logger package logger
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	globalLogger *logrus.Logger
	once         sync.Once
)

// Get Get
func Get() *logrus.Logger {
	if globalLogger != nil {
		return globalLogger
	}

	once.Do(initLogger)

	for {
		if globalLogger != nil {
			break
		}
	}
	return globalLogger
}

type logConfig struct {
	jsonFormat bool
	isDev      bool
	logLevel   int
}

func initLogger() {
	// Get current path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Clean(filepath.Dir(ex))

	logCfg := parseEnvToLogConfig()
	newLogger := logrus.New()

	var formatter logrus.Formatter
	if logCfg.jsonFormat {
		formatter = &logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   time.RFC3339,
			PrettyPrint:       false,
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
			ForceQuote:       true,
			CallerPrettyfier: newCallerPrettyfier(basePath),
		}
	}

	if logCfg.isDev {
		newLogger.SetReportCaller(true)
	}

	newLogger.SetFormatter(formatter)
	newLogger.SetLevel(logrus.Level(logCfg.logLevel))
	newLogger.SetOutput(os.Stdout)
	newLogger.Hooks.Add(fileHook(basePath))

	globalLogger = newLogger
}

func fileHook(basePath string) *lfshook.LfsHook {
	date := time.Now().Format(time.RFC3339)
	date = strings.ReplaceAll(date, ":", "")

	pathMap := lfshook.PathMap{
		logrus.PanicLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.FatalLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.ErrorLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.WarnLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.InfoLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
	}

	return lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   time.RFC3339,
			PrettyPrint:       true,
		},
	)
}

func newCallerPrettyfier(basePath string) func(*runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
		fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
		return fmt.Sprintf("[%s:%d]", fileName, frame.Line), ""
	}
}

func parseEnvToLogConfig() logConfig {
	var cfg logConfig
	if mode := os.Getenv("LOG_FORMAT"); mode == "json" {
		cfg.jsonFormat = true
	}

	if deployment := os.Getenv("DEPLOYMENT"); deployment == "dev" {
		cfg.isDev = true
	}

	logLevelString := os.Getenv("LOG_LEVEL")
	switch logLevelString {
	case "panic":
		cfg.logLevel = PanicLevel
	case "fatal":
		cfg.logLevel = FatalLevel
	case "error":
		cfg.logLevel = ErrorLevel
	case "warn":
		cfg.logLevel = WarnLevel
	case "info":
		cfg.logLevel = InfoLevel
	case "debug":
		cfg.logLevel = DebugLevel
	case "trace":
		cfg.logLevel = TraceLevel
	default:
		cfg.logLevel = InfoLevel
	}

	return cfg
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	PanicLevel int = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)
