package windy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"windy-cams/pkg/config"
)

const baseURL = "https://api.windy.com/webcams/api/v3/webcams"

// Client manages communication with the Windy Webcams API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Windy API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// APIResponse represents the top-level structure of the API response.
type APIResponse struct {
	Result Result `json:"result"`
}

// Result holds the list of webcams from the API.
type Result struct {
	Webcams []Webcam `json:"webcams"`
}

// Webcam represents a single webcam's data from the API.
type Webcam struct {
	ID       uint64   `json:"id"`
	Status   string   `json:"status"`
	Title    string   `json:"title"`
	Location Location `json:"location"`
}

// Location represents the geographical details of a webcam.
type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Continent string  `json:"continent"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetWebcams fetches a list of webcams from the Windy API based on configuration.
func (c *Client) GetWebcams(cfg *config.Config) ([]Webcam, error) {
	url := fmt.Sprintf("%s?limit=%d&offset=%d&sortKey=%s&sortDirection=%s&continents=%s",
		baseURL,
		cfg.APILimit,
		cfg.APIOffset,
		cfg.APISortKey,
		cfg.APISortDirection,
		cfg.APIContinents,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("x-windy-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return apiResponse.Result.Webcams, nil
}
