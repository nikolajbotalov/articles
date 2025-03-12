package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"sync"
)

type Config struct {
	Listen     Listen
	PostgreSQL PostgreSQL
}

type Listen struct {
	BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
	Port   string `env:"PORT" env-default:"8080"`
}

type PostgreSQL struct {
	Username string `env:"PSQL_USERNAME" env-default:"postgres"`
	Password string `env:"PSQL_PASSWORD" env-default:"admin"`
	Host     string `env:"PSQL_HOST" env-default:"host.docker.internal"`
	Port     string `env:"PSQL_PORT" env-default:"5432"`
	Database string `env:"PSQL_DB" env-default:"blog"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := slog.Logger{}

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			logger.Error(err.Error())
		}
	})

	return instance
}
