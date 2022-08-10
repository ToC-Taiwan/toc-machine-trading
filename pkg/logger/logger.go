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

	"toc-machine-trading/pkg/global"

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

	var dev bool
	mode, ok := os.LookupEnv("DEPLOYMENT")
	if !ok || mode != "prod" {
		dev = true
	}

	if dev {
		globalLogger.SetReportCaller(true)
		globalLogger.SetFormatter(&logrus.TextFormatter{
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
		})
	} else {
		globalLogger.SetFormatter(&logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   global.LongTimeLayout,
			PrettyPrint:       false,
		})
	}

	globalLogger.SetLevel(logrus.InfoLevel)
	if dev {
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
			// CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// 	fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
			// 	return fmt.Sprintf("%s:%d", fileName, frame.Line), ""
			// },
		},
	)
}
