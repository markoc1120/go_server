package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl       string
	Platform    string
	Secret      string
	PolkaAPIKey string
	Port        string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		DBUrl:       os.Getenv("DB_URL"),
		Platform:    os.Getenv("PLATFORM"),
		Secret:      os.Getenv("SECRET"),
		PolkaAPIKey: os.Getenv("POLKA_KEY"),
		Port:        getEnvDefault("PORT", "8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBUrl == "" {
		return errors.New("DB_URL must be set")
	}
	if c.Platform == "" {
		return errors.New("PLATFORM must be set")
	}
	if c.Secret == "" {
		return errors.New("SECRET must be set")
	}
	if c.PolkaAPIKey == "" {
		return errors.New("POLKA_KEY must be set")
	}
	return nil
}

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
