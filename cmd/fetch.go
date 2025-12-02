package main

import (
	"log"

	"github.com/basel-ax/windy-cams/pkg/config"
	"github.com/basel-ax/windy-cams/pkg/storage"
	"github.com/basel-ax/windy-cams/pkg/windy"
)

func runFetch(devMode bool) {
	log.Println("Starting Windy Cams data fetcher...")

	// Load configuration
	cfg := config.New()
	log.Println("Configuration loaded successfully.")

	// Initialize database
	db := storage.New(cfg)
	log.Println("Database connection established and schema migrated.")

	// Create Windy API client
	windyClient := windy.NewClientWithDevMode(cfg.WindyAPIKey, devMode)
	log.Println("Windy API client created.")

	// Fetch webcams from the API
	log.Println("Fetching webcams from Windy API...")
	apiWebcams, _, err := windyClient.GetWebcams(cfg)
	if err != nil {
		log.Fatalf("Failed to fetch webcams: %v", err)
	}
	log.Printf("Fetched %d webcams from the API.", len(apiWebcams))

	// Transform and save webcams to the database
	savedCount := 0
	for _, apiWebcam := range apiWebcams {
		dbWebcam := storage.Webcam{
			WebcamID:  apiWebcam.WebcamID,
			ApiStatus: apiWebcam.Status,
			Status:    "New",
			Title:     apiWebcam.Title,
			Latitude:  apiWebcam.Location.Latitude,
			Longitude: apiWebcam.Location.Longitude,
			City:      apiWebcam.Location.City,
			Country:   apiWebcam.Location.Country,
			Continent: apiWebcam.Location.Continent,
		}

		// Using FirstOrCreate to avoid duplicates based on WebcamID
		result := db.Where(storage.Webcam{WebcamID: dbWebcam.WebcamID}).FirstOrCreate(&dbWebcam)
		if result.Error != nil {
			log.Printf("Failed to save webcam %d: %v", dbWebcam.WebcamID, result.Error)
			continue
		}
		if result.RowsAffected > 0 {
			savedCount++
		}
	}

	log.Printf("Successfully saved %d new webcams to the database.", savedCount)
	log.Println("Windy Cams data fetcher finished successfully.")
}
