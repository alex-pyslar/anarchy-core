package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv" // Для загрузки переменных из .env файла
)

// Config struct holds all application configurations.
type Config struct {
	AppPort      string // Порт, на котором будет работать приложение
	DatabaseURL  string // URL для подключения к PostgreSQL
	JWTSecretKey string // Секретный ключ для подписи JWT токенов
}

// LoadConfig loads configuration from environment variables.
// It tries to load from a .env file first.
func LoadConfig() (*Config, error) {
	// Try to load .env file, ignore if not found
	godotenv.Load()

	cfg := &Config{
		AppPort:      os.Getenv("APP_PORT"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	// Validate required configurations
	if cfg.AppPort == "" {
		cfg.AppPort = "8080" // Default port if not set
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	if cfg.JWTSecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is not set")
	}

	return cfg, nil
}
