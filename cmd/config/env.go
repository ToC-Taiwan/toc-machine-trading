package config

type Database struct {
	DBName  string `env:"DB_NAME" env-required:"true"`
	URL     string `env:"DB_URL" env-required:"true"`
	PoolMax int    `env:"DB_POOL_MAX" env-required:"true"`
}

type Server struct {
	HTTP                      string `env:"HTTP" env-required:"true"`
	RouterDebugMode           string `env:"GIN_MODE" env-required:"true"`
	DisableSwaggerHTTPHandler string `env:"DISABLE_SWAGGER_HTTP_HANDLER" env-required:"true"`
}

// Sinopac -.
type Sinopac struct {
	PoolMax int    `env:"SINOPAC_POOL_MAX" env-required:"true"`
	URL     string `env:"SINOPAC_URL" env-required:"true"`
}

// Fugle -.
type Fugle struct {
	PoolMax int    `env:"FUGLE_POOL_MAX" env-required:"true"`
	URL     string `env:"FUGLE_URL" env-required:"true"`
}

// RabbitMQ -.
type RabbitMQ struct {
	URL      string `env:"RABBITMQ_URL" env-required:"true"`
	Exchange string `env:"RABBITMQ_EXCHANGE" env-required:"true"`
	WaitTime int64  `env:"RABBITMQ_WAIT_TIME" env-required:"true"`
	Attempts int    `env:"RABBITMQ_ATTEMPTS" env-required:"true"`
}
