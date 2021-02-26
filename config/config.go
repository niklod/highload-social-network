package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var SessionName = "hsn-session"

type Config struct {
	DB        *DBConfig
	Server    *HTTPServerConfig
	Tarantool *Tarantool
	SecretKey string `envconfig:"SESSION_SECRET_KEY" default:"verysecretkey"`
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"3306"`
	Login    string `envconfig:"DB_LOGIN" default:"niklod"`
	Password string `envconfig:"DB_PASSWORD" default:"VLQi4Vttuo6wFRqm"`
	DBName   string `envconfig:"DB_NAME" default:"hsn"`
}

func (d *DBConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Login, d.Password, d.Host, d.Port, d.DBName)
}

type Tarantool struct {
	Host     string `envconfig:"TARANTOOL_HOST" default:"localhost"`
	Port     int    `envconfig:"TARANTOOL_PORT" default:"3013"`
	Login    string `envconfig:"TARANTOOL_LOGIN" default:"niklod"`
	Password string `envconfig:"TARANTOOL_PASSWORD" default:"VLQi4Vttuo6wFRqm"`
}

type HTTPServerConfig struct {
	Port int `envconfig:"HTTP_SERVER_PORT" default:"8080"`
}

func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("creating config: %w", err)
	}

	return &cfg, nil
}
