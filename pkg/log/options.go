// Package log package log
package log

const (
	FormatJSON string = "json"
	FormatText string = "text"
)

type config struct {
	level      Level
	needCaller bool
	timeFormat string
	format     string
	fileName   string
}

// Option -.
type Option func(*config)

func LogLevel(level string) Option {
	return func(c *config) {
		switch level {
		case LevelPanic.String(), LevelFatal.String(), LevelError.String(), LevelWarn.String(), LevelInfo.String(), LevelDebug.String(), LevelTrace.String():
			c.level = Level(level)
		default:
			c.level = _defaultLogLevel
		}
	}
}

func TimeFormat(format string) Option {
	return func(c *config) {
		c.timeFormat = format
	}
}

func LogFormat(format string) Option {
	return func(c *config) {
		if format != "json" && format != "text" {
			c.format = _defaultLogFormat
		}
		c.format = format
	}
}

func NeedCaller(need bool) Option {
	return func(c *config) {
		c.needCaller = need
	}
}

func FileName(name string) Option {
	return func(c *config) {
		c.fileName = name
	}
}
