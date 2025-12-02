package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	// Windy API configuration
	WindyAPIKey      string `env:"WINDY_API_KEY,required"`
	APILimit         int    `env:"API_LIMIT" envDefault:"50"`
	APIOffset        int    `env:"API_OFFSET" envDefault:"0"`
	APISortKey       string `env:"API_SORT_KEY" envDefault:"createdOn"`
	APISortDirection string `env:"API_SORT_DIRECTION" envDefault:"desc"`
	APIContinents    string `env:"API_CONTINENTS" envDefault:"AF"`

	// Database configuration
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSchema   string `env:"DB_SCHEMA" envDefault:"public"`
}

// New loads the configuration from environment variables.
func New() *Config {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse configuration: %+v", err)
	}
	return cfg
}
