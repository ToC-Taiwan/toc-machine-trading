// Package logger package logger
package logger

import (
	"tmt/cmd/config"
	"tmt/pkg/logger"
)

func Get() *logger.Logger {
	cfg := config.GetConfig().Log
	return logger.NewLogger(
		logger.LogFormat(logger.Format(cfg.LogFormat)),
		logger.LogLevel(logger.Level(cfg.Level)),
		logger.NeedCaller(cfg.NeedCaller),
	)
}
