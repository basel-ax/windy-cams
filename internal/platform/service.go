package webcam

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/basel-ax/windy-cams/pkg/windy"
)

// Service defines the interface for platform business logic.
type Service interface {
	FetchAndStorePlatforms(ctx context.Context) error
}

type service struct {
	windyClient *windy.Client
	repo        Repository
	logger      *slog.Logger
}

// NewService creates a new platform service.
func NewService(repo Repository, windyClient *windy.Client, logger *slog.Logger) Service {
	return &service{
		repo:        repo,
		windyClient: windyClient,
		logger:      logger,
	}
}

// FetchAndStorePlatforms fetches platforms from Windy API and stores them in the repository.
func (s *service) FetchAndStorePlatforms(ctx context.Context) error {
	s.logger.Info("fetching platforms from windy api")
	platforms, err := s.windyClient.GetPlatforms(ctx)
	if err != nil {
		s.logger.Error("failed to get platforms from windy", slog.Any("error", err))
		return fmt.Errorf("failed to get platforms from windy: %w", err)
	}

	if len(platforms) == 0 {
		s.logger.Info("no new platforms to save")
		return nil
	}

	s.logger.Info("saving platforms", slog.Int("count", len(platforms)))
	if err := s.repo.SavePlatforms(ctx, platforms); err != nil {
		s.logger.Error("failed to save platforms", slog.Any("error", err))
		return fmt.Errorf("failed to save platforms: %w", err)
	}

	s.logger.Info("successfully saved platforms", slog.Int("count", len(platforms)))
	return nil
}
