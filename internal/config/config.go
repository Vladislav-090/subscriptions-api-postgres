package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

const defaultJWTTTL = 24 * time.Hour

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtTTL, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		jwtTTL = defaultJWTTTL
	}

	cfg := Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},

		Database: DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},

		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			TTL:    jwtTTL,
		},
	}

	return &cfg, nil
}
