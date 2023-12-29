package config

type EnvConfig struct {
	Database Database
	Server   Server
	Sinopac  Sinopac
	Fugle    Fugle
	RabbitMQ RabbitMQ
	SMTP     SMTP
}

type Database struct {
	DBName  string
	URL     string
	PoolMax int
}

type Server struct {
	HTTP                      string
	DisableSwaggerHTTPHandler string
}

// Sinopac -.
type Sinopac struct {
	PoolMax int
	URL     string
}

// Fugle -.
type Fugle struct {
	PoolMax int
	URL     string
}

// RabbitMQ -.
type RabbitMQ struct {
	URL      string
	Exchange string
	WaitTime int64
	Attempts int
}

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}
