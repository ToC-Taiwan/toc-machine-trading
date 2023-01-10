package log

import "github.com/sirupsen/logrus"

type Level string

const (
	LevelPanic Level = "panic"
	LevelFatal Level = "fatal"
	LevelError Level = "error"
	LevelWarn  Level = "warn"
	LevelInfo  Level = "info"
	LevelDebug Level = "debug"
	LevelTrace Level = "trace"
)

func (l Level) String() string {
	return string(l)
}

func (l Level) Level() logrus.Level {
	switch l {
	case LevelPanic:
		return logrus.PanicLevel
	case LevelFatal:
		return logrus.FatalLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelDebug:
		return logrus.DebugLevel
	case LevelTrace:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
