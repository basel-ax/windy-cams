package main

import (
	"context"
	"log"

	"project/configs"
	"project/internal/platform"
	"project/pkg/database"
	"project/pkg/windy"
)

func main() {
	ctx := context.Background()

	// 1. Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if cfg.WindyAPIKey == "" || cfg.WindyAPIKey == "your_api_key_here" {
		log.Fatal("WINDY_API_KEY is not set in your .env file")
	}

	// 2. Initialize database
	db, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migration completed.")

	// 4. Initialize clients, repositories, and services
	windyClient := windy.NewClient(cfg.WindyAPIKey)
	platformRepo := platform.NewRepository(db)
	platformService := platform.NewService(platformRepo, windyClient)

	// 5. Run the business logic
	log.Println("Fetching and storing platforms...")
	if err := platformService.FetchAndStorePlatforms(ctx); err != nil {
		log.Fatalf("failed to fetch and store platforms: %v", err)
	}
	log.Println("Process finished successfully.")
}
