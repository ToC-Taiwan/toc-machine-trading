package log

type env struct {
	Level      string `env:"LOG_LEVEL"`
	Format     string `env:"LOG_FORMAT"`
	NeedCaller bool   `env:"LOG_NEED_CALLER"`
}
