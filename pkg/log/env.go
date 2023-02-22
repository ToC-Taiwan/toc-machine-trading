package log

type env struct {
	Level      string `env:"LOG_LEVEL" env-required:"true"`
	Format     string `env:"LOG_FORMAT" env-required:"true"`
	NeedCaller bool   `env:"LOG_NEED_CALLER" env-required:"true"`

	Token     string `env:"SLACK_TOKEN"`
	ChannelID string `env:"SLACK_CHANNEL_ID"`
}
