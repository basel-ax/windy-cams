package windy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/basel-ax/windy-cams/pkg/config"
)

const baseURL = "https://api.windy.com/webcams/api/v3/webcams"

// Client manages communication with the Windy Webcams API.
type Client struct {
	apiKey     string
	httpClient *http.Client
	DevMode    bool
}

// NewClient creates a new Windy API client.
func NewClient(apiKey string) *Client {
	return NewClientWithDevMode(apiKey, false)
}

// NewClientWithDevMode creates a new Windy API client with dev mode.
func NewClientWithDevMode(apiKey string, devMode bool) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		DevMode: devMode,
	}
}

// Result holds the list of webcams and total count from the API.
type Result struct {
	Total   int      `json:"total"`
	Webcams []Webcam `json:"webcams"`
}

// Webcam represents a single webcam's data from the API.
type Webcam struct {
	WebcamID uint64   `json:"webcamId"`
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
func (c *Client) GetWebcams(cfg *config.Config) ([]Webcam, int, error) {
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
		return nil, 0, fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("x-windy-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read API response body: %w", err)
	}

	if c.DevMode {
		log.Printf("Windy API Response:\n%s\n", string(body))
	}

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode API response: %w", err)
	}

	return result.Webcams, result.Total, nil
}

// GetWebcamsWithParams fetches a list of webcams with specific parameters.
func (c *Client) GetWebcamsWithParams(continent string, limit, offset int) ([]Webcam, int, error) {
	url := fmt.Sprintf("%s?limit=%d&offset=%d&sortKey=createdOn&sortDirection=desc&continents=%s",
		baseURL,
		limit,
		offset,
		continent,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("x-windy-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read API response body: %w", err)
	}

	if c.DevMode {
		log.Printf("Windy API Response:\n%s\n", string(body))
	}

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode API response: %w", err)
	}

	return result.Webcams, result.Total, nil
}
