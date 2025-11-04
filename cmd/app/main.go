package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/basel-ax/windy-cams/configs"
	"github.com/basel-ax/windy-cams/internal/webcam"
	"github.com/basel-ax/windy-cams/pkg/database"
	"github.com/basel-ax/windy-cams/pkg/windy"
)

func main() {
	devMode := flag.Bool("dev", false, "Enable developer mode for verbose logging")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

	// 1. Load configuration
	logger.Info("loading configuration")
	cfg, err := configs.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}
	if cfg.WindyAPIKey == "" || cfg.WindyAPIKey == "your_api_key_here" {
		logger.Error("WINDY_API_KEY is not set in your .env file")
		os.Exit(1)
	}

	// 2. Initialize database
	logger.Info("initializing database client")
	db, err := database.NewClient(cfg)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	// 3. Run migrations
	logger.Info("running database migrations")
	if err := database.AutoMigrate(db); err != nil {
		logger.Error("failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("database migration completed")

	// 4. Initialize clients, repositories, and services
	windyClient := windy.NewClient(cfg.WindyAPIKey)
	if *devMode {
		windyClient.WithLogger(logger)
	}
	webcamRepo := webcam.NewRepository(db)
	webcamService := webcam.NewService(webcamRepo, windyClient, logger)

	// 5. Run the business logic
	logger.Info("fetching and storing webcams")
	if err := webcamService.FetchAndStoreWebcams(ctx); err != nil {
		logger.Error("failed to fetch and store webcams", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("process finished successfully")
}
