package logger

// Option -.
type Option func(*loggerConfig)

func TimeFormat(format string) Option {
	return func(c *loggerConfig) {
		c.timeFormat = format
	}
}

func LogFormat(format Format) Option {
	return func(c *loggerConfig) {
		if format == "" {
			return
		}

		if format != "json" && format != "text" {
			return
		}

		c.format = format
	}
}

func NeedCaller(need bool) Option {
	return func(c *loggerConfig) {
		c.needCaller = need
	}
}

func LogLevel(level Level) Option {
	return func(c *loggerConfig) {
		if level == "" {
			return
		}

		if level != "panic" && level != "fatal" && level != "error" && level != "warn" && level != "info" && level != "debug" && level != "trace" {
			return
		}

		c.level = level
	}
}
