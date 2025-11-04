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
		{
			ID:      "123",
			Title:   "Test Webcam",
			Status:  "available",
			ViewURL: "http://example.com/webcam.jpg",
		},
	}

	testCases := []struct {
		name    string
		handler http.HandlerFunc
		want    []domain.Webcam
		wantErr bool
	}{
		{
			name: "successful response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v3/webcams" {
					t.Errorf("expected to request '/api/v3/webcams', got: %s", r.URL.Path)
				}
				if r.Header.Get("x-windy-api-key") != "test-api-key" {
					t.Errorf("missing or incorrect x-windy-api-key header")
				}
				w.WriteHeader(http.StatusOK)
				w.Write(successBody)
			},
			want:    expectedWebcams,
			wantErr: false,
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "malformed JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{`))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := NewClient("test-api-key")
			client.BaseURL = server.URL

			webcams, err := client.GetWebcams(context.Background())

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
			}

			if !reflect.DeepEqual(webcams, tc.want) {
				t.Errorf("expected webcams %+v, got %+v", tc.want, webcams)
			}
		})
	}
}

func TestExportAllWebcams(t *testing.T) {
	t.Parallel()

	successResponse := []domain.Webcam{
		{
			ID:      "456",
			Title:   "Export Test Webcam",
			Status:  "active",
			ViewURL: "http://example.com/export_webcam.jpg",
		},
	}

	successBody, _ := json.Marshal(successResponse)

	testCases := []struct {
		name    string
		handler http.HandlerFunc
		want    []domain.Webcam
		wantErr bool
	}{
		{
			name: "successful response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/export/all-webcams.json" {
					t.Errorf("expected to request '/export/all-webcams.json', got: %s", r.URL.Path)
				}
				if r.Header.Get("x-windy-api-key") != "test-api-key" {
					t.Errorf("missing or incorrect x-windy-api-key header")
				}
				w.WriteHeader(http.StatusOK)
				w.Write(successBody)
			},
			want:    successResponse,
			wantErr: false,
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "malformed JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[`))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := NewClient("test-api-key")
			client.BaseURL = server.URL

			webcams, err := client.ExportAllWebcams(context.Background())

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
			}

			if !reflect.DeepEqual(webcams, tc.want) {
				t.Errorf("expected webcams %+v, got %+v", tc.want, webcams)
			}
		})
	}
}
