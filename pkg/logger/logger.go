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
	"tmt/pkg/global"

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

func initLogger() {
	// Get current path
	basePath := global.GetBasePath()
	globalLogger = logrus.New()

	var jsonFormat, prodMode bool
	mode, ok := os.LookupEnv("LOG_FORMAT")
	if !ok || mode == "json" {
		jsonFormat = true
	}

	deployment, ok := os.LookupEnv("DEPLOYMENT")
	if !ok || deployment == "prod" {
		prodMode = true
	}

	var formatter logrus.Formatter
	if jsonFormat {
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

	globalLogger.SetFormatter(formatter)
	globalLogger.SetLevel(logrus.InfoLevel)

	if !prodMode {
		globalLogger.SetReportCaller(true)
		globalLogger.SetLevel(logrus.TraceLevel)
	}

	globalLogger.SetOutput(os.Stdout)
	globalLogger.Hooks.Add(fileHook(basePath))
}

func fileHook(basePath string) *lfshook.LfsHook {
	date := time.Now().Format("20060102")
	date = strings.ReplaceAll(date, ":", "")

	pathMap := lfshook.PathMap{
		logrus.PanicLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.FatalLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.ErrorLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.WarnLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.InfoLevel:  filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.DebugLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
		logrus.TraceLevel: filepath.Join(basePath, fmt.Sprintf("/logs/tmt-%s.log", date)),
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
