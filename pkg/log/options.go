// Package log package log
package log

type config struct {
	timeFormat     string
	format         Format
	level          Level
	needCaller     bool
	fileName       string
	slackToken     string
	slackChannelID string
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
		c.format = Format(format)
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

func SlackToken(token string) Option {
	return func(c *config) {
		c.slackToken = token
	}
}

func SlackChannelID(id string) Option {
	return func(c *config) {
		c.slackChannelID = id
	}
}
