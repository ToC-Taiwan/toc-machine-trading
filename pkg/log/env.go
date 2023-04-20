package log

type env struct {
	Level          string `env:"LOG_LEVEL"`
	Format         string `env:"LOG_FORMAT"`
	NeedCaller     bool   `env:"LOG_NEED_CALLER"`
	SlackToken     string `env:"SLACK_TOKEN"`
	SlackChannelID string `env:"SLACK_CHANNEL_ID"`
}
