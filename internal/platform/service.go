package webcam

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/basel-ax/windy-cams/pkg/windy"
)

// Service defines the interface for webcam business logic.
type Service interface {
	FetchAndStoreWebcams(ctx context.Context) error
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

// FetchAndStoreWebcams fetches webcams from Windy API and stores them in the repository.
func (s *service) FetchAndStoreWebcams(ctx context.Context) error {
	s.logger.Info("fetching webcams from windy api")
	webcams, err := s.windyClient.GetWebcams(ctx)
	if err != nil {
		s.logger.Error("failed to get webcams from windy", slog.Any("error", err))
		return fmt.Errorf("failed to get webcams from windy: %w", err)
	}

	if len(webcams) == 0 {
		s.logger.Info("no new webcams to save")
		return nil
	}

	s.logger.Info("saving webcams", slog.Int("count", len(webcams)))
	if err := s.repo.SaveWebcams(ctx, webcams); err != nil {
		s.logger.Error("failed to save webcams", slog.Any("error", err))
		return fmt.Errorf("failed to save webcams: %w", err)
	}

	s.logger.Info("successfully saved webcams", slog.Int("count", len(webcams)))
	return nil
}
