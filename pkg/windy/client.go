package windy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/basel-ax/windy-cams/internal/domain"
)

const webcamsURL = "https://api.windy.com/api/webcams/v2/list"

// windyWebcam defines the structure of a webcam object from the Windy API response.
type windyWebcam struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	URL    struct {
		Current struct {
			Desktop string `json:"desktop"`
		} `json:"current"`
	} `json:"url"`
}

// windyResponse defines the structure of the top-level API response from Windy.
type windyResponse struct {
	Result struct {
		Webcams []windyWebcam `json:"webcams"`
	} `json:"result"`
}

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

// GetWebcams fetches the list of webcams from the Windy API.
func (c *Client) GetWebcams(ctx context.Context) ([]domain.Webcam, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, webcamsURL, nil)
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

	var windyResp windyResponse
	if err := json.NewDecoder(resp.Body).Decode(&windyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	webcams := make([]domain.Webcam, 0, len(windyResp.Result.Webcams))
	for _, windyCam := range windyResp.Result.Webcams {
		webcams = append(webcams, domain.Webcam{
			ID:      windyCam.ID,
			Title:   windyCam.Title,
			Status:  windyCam.Status,
			ViewURL: windyCam.URL.Current.Desktop,
		})
	}

	return webcams, nil
}
