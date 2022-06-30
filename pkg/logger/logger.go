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
	globalLogger.SetReportCaller(true)

	var dev bool
	mode, ok := os.LookupEnv("DEPLOYMENT")
	if !ok || mode != "prod" {
		dev = true
	}

	if dev {
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
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
				return fmt.Sprintf("%s:%d", fileName, frame.Line), ""
			},
		})
	}

	globalLogger.SetLevel(logrus.TraceLevel)
	globalLogger.SetOutput(os.Stdout)
	globalLogger.Hooks.Add(fileHook(basePath))
}

func fileHook(basePath string) *lfshook.LfsHook {
	folderName := time.Now().Format("20060102")
	folderName = strings.ReplaceAll(folderName, ":", "")

	pathMap := lfshook.PathMap{
		logrus.PanicLevel: filepath.Join(basePath, "/logs/", folderName, "/panic.json"),
		logrus.FatalLevel: filepath.Join(basePath, "/logs/", folderName, "/fetal.json"),
		logrus.ErrorLevel: filepath.Join(basePath, "/logs/", folderName, "/error.json"),
		logrus.WarnLevel:  filepath.Join(basePath, "/logs/", folderName, "/warn.json"),
		logrus.InfoLevel:  filepath.Join(basePath, "/logs/", folderName, "/info.json"),
		logrus.DebugLevel: filepath.Join(basePath, "/logs/", folderName, "/debug.json"),
		logrus.TraceLevel: filepath.Join(basePath, "/logs/", folderName, "/error.json"),
	}

	return lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{
			DisableHTMLEscape: true,
			TimestampFormat:   global.LongTimeLayout,
			PrettyPrint:       false,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := strings.ReplaceAll(frame.File, fmt.Sprintf("%s/", basePath), "")
				return fmt.Sprintf("%s:%d", fileName, frame.Line), ""
			},
		},
	)
}
