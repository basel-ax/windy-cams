package windy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"project/internal/domain"
)

const platformsURL = "https://api.windy.com/api/webcams/v2/platforms"

// Client is a client for the Windy Webcams API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Windy API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GetPlatforms fetches all asset platforms from the Windy API.
func (c *Client) GetPlatforms(ctx context.Context) ([]domain.Platform, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, platformsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-windy-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var platforms []domain.Platform
	if err := json.NewDecoder(resp.Body).Decode(&platforms); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return platforms, nil
}
