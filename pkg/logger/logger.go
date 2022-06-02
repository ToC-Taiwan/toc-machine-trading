// Package logger package logger
package logger

import (
	"os"
	"path/filepath"
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

func initLogger() {
	if globalLogger != nil {
		return
	}

	// Get current path
	basePath := global.GetBasePath()
	// if basePath == "" {
	// 	basePath = global.GetRuntimePath()
	// }

	// create new instance
	globalLogger = logrus.New()
	if global.GetIsDevelopment() {
		globalLogger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
			ForceQuote:       true,
		})
	} else {
		globalLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: global.LongTimeLayout,
			PrettyPrint:     false,
		})
	}

	// globalLogger.SetReportCaller(true)
	folderName := time.Now().Format(time.RFC3339)[:10]
	folderName = strings.ReplaceAll(folderName, ":", "")
	globalLogger.SetLevel(logrus.TraceLevel)
	globalLogger.SetOutput(os.Stdout)
	pathMap := lfshook.PathMap{
		logrus.PanicLevel: filepath.Join(basePath, "/logs/", folderName, "/panic.json"),
		logrus.FatalLevel: filepath.Join(basePath, "/logs/", folderName, "/fetal.json"),
		logrus.ErrorLevel: filepath.Join(basePath, "/logs/", folderName, "/error.json"),
		logrus.WarnLevel:  filepath.Join(basePath, "/logs/", folderName, "/warn.json"),
		logrus.InfoLevel:  filepath.Join(basePath, "/logs/", folderName, "/info.json"),
		logrus.DebugLevel: filepath.Join(basePath, "/logs/", folderName, "/debug.json"),
		logrus.TraceLevel: filepath.Join(basePath, "/logs/", folderName, "/error.json"),
	}
	globalLogger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
}

// Get Get
func Get() *logrus.Logger {
	if globalLogger != nil {
		return globalLogger
	}
	once.Do(initLogger)
	return globalLogger
}
