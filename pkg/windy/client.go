package windy

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/basel-ax/windy-cams/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

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
	BaseURL    string
	logger     *slog.Logger
	tracer     trace.Tracer
}

// NewClient creates a new Windy API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		BaseURL:    "https://api.windy.com/webcams/",
		tracer:     otel.Tracer("windy-client"),
	}
}

// WithLogger sets a logger for the client to enable request logging.
func (c *Client) WithLogger(logger *slog.Logger) *Client {
	c.logger = logger
	return c
}

// GetWebcams fetches the list of webcams from the Windy API.
func (c *Client) GetWebcams(ctx context.Context) ([]domain.Webcam, error) {
	ctx, span := c.tracer.Start(ctx, "windy.client.GetWebcams")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"api/v3/webcams", nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.logger != nil {
		c.logger.Info("sending request to windy api", slog.String("method", req.Method), slog.String("url", req.URL.String()))
	}

	req.Header.Set("x-windy-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var windyResp windyResponse
	if err := json.NewDecoder(resp.Body).Decode(&windyResp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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

	span.SetAttributes(attribute.Int("webcams.count", len(webcams)))

	return webcams, nil
}

// ExportAllWebcams fetches a JSON file with all webcams from the Windy export endpoint.
func (c *Client) ExportAllWebcams(ctx context.Context) ([]domain.Webcam, error) {
	ctx, span := c.tracer.Start(ctx, "windy.client.ExportAllWebcams")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"export/all-webcams.json", nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.logger != nil {
		c.logger.Info("sending request to windy api for export", slog.String("method", req.Method), slog.String("url", req.URL.String()))
	}

	req.Header.Set("x-windy-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var webcams []domain.Webcam
	// The export endpoint returns a flat array of webcam objects.
	// We make an assumption here that its structure is compatible with domain.Webcam.
	if err := json.NewDecoder(resp.Body).Decode(&webcams); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	span.SetAttributes(attribute.Int("webcams.count", len(webcams)))

	return webcams, nil
}
