package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Address string `yaml:"address"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	Secret          string        `env:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
}

type Config struct {
	Env      string         `env:"APP_ENV"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

func Load() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("internal/config/config.yaml", cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	return cfg
}
