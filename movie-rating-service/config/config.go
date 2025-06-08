package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"dev"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
	DbConfig    DatabaseConfig
	Port        int    `env:"PORT" envDefault:"8080"`
	JWTSecret   string `env:"JWT_SECRET" envDefault:"secret"`
	Migrate     bool   `env:"MIGRATE" envDefault:"true"`
}

type DatabaseConfig struct {
	User     string `env:"DB_USER" envDefault:"thermondo_user"`
	Password string `env:"DB_PASSWORD" envDefault:"thermondo_pass"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Name     string `env:"DB_NAME" envDefault:"thermondo"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

var Cfg Config

func Init() error {
	err := env.Parse(&Cfg)
	if err != nil {
		return fmt.Errorf("error occurred while parsing environment variables: %w", err)
	}
	return nil
}
