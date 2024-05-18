package config

type EnvConfig struct {
	Database Database `json:"Database" yaml:"Database"`
	Server   Server   `json:"Server" yaml:"Server"`
	Sinopac  Sinopac  `json:"Sinopac" yaml:"Sinopac"`
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
	URL string `json:"URL" yaml:"URL"`
}

type SMTP struct {
	Host     string `json:"Host" yaml:"Host"`
	Port     int    `json:"Port" yaml:"Port"`
	Username string `json:"Username" yaml:"Username"`
	Password string `json:"Password" yaml:"Password"`
}
