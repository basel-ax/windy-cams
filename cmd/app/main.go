package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/basel-ax/windy-cams/configs"
	"github.com/basel-ax/windy-cams/internal/webcam"
	"github.com/basel-ax/windy-cams/pkg/database"
	"github.com/basel-ax/windy-cams/pkg/observability"
	"github.com/basel-ax/windy-cams/pkg/windy"
)

func main() {
	devMode := flag.Bool("dev", false, "Enable developer mode for verbose logging")
	exportAll := flag.Bool("export-all", false, "Export all webcams and store them")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

	// Initialize OpenTelemetry
	tp, err := observability.InitTracerProvider()
	if err != nil {
		logger.Error("failed to initialize tracer provider", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error("failed to shutdown tracer provider", slog.Any("error", err))
		}
	}()

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
	if *exportAll {
		logger.Info("exporting and storing all webcams")
		webcams, err := windyClient.ExportAllWebcams(ctx)
		if err != nil {
			logger.Error("failed to export all webcams", slog.Any("error", err))
			os.Exit(1)
		}
		if *devMode {
			logger.Info("successfully exported webcams", slog.Int("count", len(webcams)))
		}
		if err := webcamRepo.SaveWebcams(ctx, webcams); err != nil {
			logger.Error("failed to save all exported webcams", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("exported webcams have been stored successfully")
	} else {
		logger.Info("fetching and storing webcams")
		if err := webcamService.FetchAndStoreWebcams(ctx); err != nil {
			logger.Error("failed to fetch and store webcams", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("process finished successfully")
	}
}
