package configs

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	WindyAPIKey string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSQLMode   string
}

// LoadConfig loads configuration from .env file.
func LoadConfig() (*Config, error) {
	// In a real app, you might not want to ignore this error.
	_ = godotenv.Load()

	return &Config{
		WindyAPIKey: os.Getenv("WINDY_API_KEY"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBSQLMode:   os.Getenv("DB_SSLMODE"),
	}, nil
}
