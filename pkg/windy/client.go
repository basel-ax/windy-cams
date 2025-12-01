package windy

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

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
func (c *Client) GetWebcams(ctx context.Context, limit, offset int) ([]domain.Webcam, error) {
	ctx, span := c.tracer.Start(ctx, "windy.client.GetWebcams")
	defer span.End()

	baseURL, err := url.Parse(c.BaseURL + "api/v3/webcams")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
