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

	"tmt/global"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	globalLogger *logrus.Logger
	once         sync.Once
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	global.SetBasePath(filepath.Clean(filepath.Dir(ex)))
}

// Get Get
func Get() *logrus.Logger {
	if globalLogger != nil {
		return globalLogger
	}

	once.Do(initLogger)
	return globalLogger
}

type logConfig struct {
	jsonFormat bool
	isDev      bool
	logLevel   int
}

func initLogger() {
	// Get current path
	logCfg := parseEnvToLogConfig()
	basePath := global.GetBasePath()
	globalLogger = logrus.New()

	var formatter logrus.Formatter
	if logCfg.jsonFormat {
		formatter = &logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   global.LongTimeLayout,
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
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
				return fmt.Sprintf("[%s:%d]", fileName, frame.Line), ""
			},
		}
	}

	if logCfg.isDev {
		globalLogger.SetReportCaller(true)
	}

	globalLogger.SetFormatter(formatter)
	globalLogger.SetLevel(logrus.Level(logCfg.logLevel))
	globalLogger.SetOutput(os.Stdout)
	globalLogger.Hooks.Add(fileHook(basePath))
}

func fileHook(basePath string) *lfshook.LfsHook {
	date := time.Now().Format(global.ShortTimeLayoutNoDash)
	date = strings.ReplaceAll(date, ":", "")

	pathMap := lfshook.PathMap{
		logrus.PanicLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.FatalLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.ErrorLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.WarnLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.InfoLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		// logrus.DebugLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		// logrus.TraceLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
	}

	return lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   global.LongTimeLayout,
			PrettyPrint:       true,
		},
	)
}

// NewCallerPrettyfier -.
func NewCallerPrettyfier(basePath string) func(*runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
		fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
		return fmt.Sprintf("[%s:%d]", fileName, frame.Line), ""
	}
}

func parseEnvToLogConfig() logConfig {
	var jsonFormat, isDev bool
	var logLevel int

	if mode := os.Getenv("LOG_FORMAT"); mode == "json" {
		jsonFormat = true
	}

	if deployment := os.Getenv("DEPLOYMENT"); deployment == "dev" {
		isDev = true
	}

	logLevelString := os.Getenv("LOG_LEVEL")
	switch logLevelString {
	case "panic":
		logLevel = PanicLevel
	case "fatal":
		logLevel = FatalLevel
	case "error":
		logLevel = ErrorLevel
	case "warn":
		logLevel = WarnLevel
	case "info":
		logLevel = InfoLevel
	case "debug":
		logLevel = DebugLevel
	case "trace":
		logLevel = TraceLevel
	default:
		logLevel = InfoLevel
	}

	return logConfig{
		jsonFormat: jsonFormat,
		isDev:      isDev,
		logLevel:   logLevel,
	}
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
