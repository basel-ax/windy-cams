package platform

import (
	"context"
	"fmt"

	"project/pkg/windy"
)

// Service defines the interface for platform business logic.
type Service interface {
	FetchAndStorePlatforms(ctx context.Context) error
}

type service struct {
	windyClient *windy.Client
	repo        Repository
}

// NewService creates a new platform service.
func NewService(repo Repository, windyClient *windy.Client) Service {
	return &service{
		repo:        repo,
		windyClient: windyClient,
	}
}

// FetchAndStorePlatforms fetches platforms from Windy API and stores them in the repository.
func (s *service) FetchAndStorePlatforms(ctx context.Context) error {
	platforms, err := s.windyClient.GetPlatforms(ctx)
	if err != nil {
		return fmt.Errorf("failed to get platforms from windy: %w", err)
	}

	if err := s.repo.SavePlatforms(ctx, platforms); err != nil {
		return fmt.Errorf("failed to save platforms: %w", err)
	}

	fmt.Printf("Successfully saved %d platforms.\n", len(platforms))
	return nil
}
