package config

type EnvConfig struct {
	Database Database `json:"Database" yaml:"Database"`
	Server   Server   `json:"Server" yaml:"Server"`
	Sinopac  Sinopac  `json:"Sinopac" yaml:"Sinopac"`
	Fugle    Fugle    `json:"Fugle" yaml:"Fugle"`
	RabbitMQ RabbitMQ `json:"RabbitMQ" yaml:"RabbitMQ"`
	SMTP     SMTP     `json:"SMTP" yaml:"SMTP"`
}

type Database struct {
	DBName  string `json:"DBName" yaml:"DBName"`
	URL     string `json:"URL" yaml:"URL"`
	PoolMax int    `json:"PoolMax" yaml:"PoolMax"`
}

type Server struct {
	HTTP string `json:"HTTP" yaml:"HTTP"`
}

// Sinopac -.
type Sinopac struct {
	PoolMax int    `json:"PoolMax" yaml:"PoolMax"`
	URL     string `json:"URL" yaml:"URL"`
}

// Fugle -.
type Fugle struct {
	PoolMax int    `json:"PoolMax" yaml:"PoolMax"`
	URL     string `json:"URL" yaml:"URL"`
}

// RabbitMQ -.
type RabbitMQ struct {
	URL      string `json:"URL" yaml:"URL"`
	Exchange string `json:"Exchange" yaml:"Exchange"`
	WaitTime int64  `json:"WaitTime" yaml:"WaitTime"`
	Attempts int    `json:"Attempts" yaml:"Attempts"`
}

type SMTP struct {
	Host     string `json:"Host" yaml:"Host"`
	Port     int    `json:"Port" yaml:"Port"`
	Username string `json:"Username" yaml:"Username"`
	Password string `json:"Password" yaml:"Password"`
}
