package main

import (
	"log"
	"windy-cams/pkg/config"
	"windy-cams/pkg/storage"
	"windy-cams/pkg/windy"
)

func main() {
	log.Println("Starting Windy Cams data fetcher...")

	// Load configuration
	cfg := config.New()
	log.Println("Configuration loaded successfully.")

	// Initialize database
	db := storage.New(cfg)
	log.Println("Database connection established and schema migrated.")

	// Create Windy API client
	windyClient := windy.NewClient(cfg.WindyAPIKey)
	log.Println("Windy API client created.")

	// Fetch webcams from the API
	log.Println("Fetching webcams from Windy API...")
	apiWebcams, err := windyClient.GetWebcams(cfg)
	if err != nil {
		log.Fatalf("Failed to fetch webcams: %v", err)
	}
	log.Printf("Fetched %d webcams from the API.", len(apiWebcams))

	// Transform and save webcams to the database
	savedCount := 0
	for _, apiWebcam := range apiWebcams {
		dbWebcam := storage.Webcam{
			WebcamID:  apiWebcam.ID,
			Status:    apiWebcam.Status,
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
