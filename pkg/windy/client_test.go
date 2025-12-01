package windy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/basel-ax/windy-cams/internal/domain"
)

func TestGetWebcams(t *testing.T) {
	t.Parallel()

	successResponse := windyResponse{
		Result: struct {
			Webcams []windyWebcam `json:"webcams"`
		}{
			Webcams: []windyWebcam{
				{
					ID:     "123",
					Title:  "Test Webcam",
					Status: "available",
					URL: struct {
						Current struct {
							Desktop string `json:"desktop"`
						} `json:"current"`
					}{
						Current: struct {
							Desktop string `json:"desktop"`
						}{
							Desktop: "http://example.com/webcam.jpg",
						},
					},
				},
			},
		},
	}

	successBody, _ := json.Marshal(successResponse)

	expectedWebcams := []domain.Webcam{
