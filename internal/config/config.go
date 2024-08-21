package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Http httpConfig
	DB   dbConfig
}

type httpConfig struct {
	Port uint16 `env:"http_port" env-default:"8000"`
}

type dbConfig struct {
	DbHost  string `env:"DB_HOST"  env-default:"localhost"`
	DbPort  string `env:"DB_PORT"  env-default:"5432"`
	DbUser  string `env:"DB_USER"  env-default:"postgres"`
	DbPass  string `env:"DB_PASS"  env-default:"root"`
	DbName  string `env:"DB_NAME"  env-default:"chat"`
	SslMode string `env:"SSL_MODE" env-default:"disable"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	return &cfg
}
