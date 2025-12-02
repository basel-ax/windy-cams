package main

import (
	"log"

	"github.com/basel-ax/windy-cams/pkg/config"
	"github.com/basel-ax/windy-cams/pkg/storage"
	"github.com/basel-ax/windy-cams/pkg/windy"
	"gorm.io/gorm"
)

func runFetchAll(devMode bool) {
	log.Println("Starting to fetch all Windy Cams data...")

	// Load configuration
	cfg := config.New()
	log.Println("Configuration loaded successfully.")

	// Initialize database
	db := storage.New(cfg, devMode)
	log.Println("Database connection established and schema migrated.")

	// Create Windy API client
	windyClient := windy.NewClientWithDevMode(cfg.WindyAPIKey, devMode)
	log.Println("Windy API client created.")

	continents := []string{"AF", "AN", "AS", "EU", "NA", "OC", "SA"}
	totalSavedCount := 0
	const limit = 50 // The number of webcams to fetch per API call

	for _, continent := range continents {
		log.Printf("Fetching webcams for continent: %s", continent)
		offset := 0
		continentSavedCount := 0

		for {
			apiWebcams, total, err := windyClient.GetWebcamsWithParams(continent, limit, offset)
			if err != nil {
				log.Printf("Failed to fetch webcams for continent %s at offset %d: %v", continent, offset, err)
				break // Stop processing this continent on error
			}

			if len(apiWebcams) == 0 {
				log.Printf("No more webcams found for continent: %s", continent)
				break // No more webcams to fetch for this continent
			}

			log.Printf("Fetched %d webcams for continent %s (total: %d)", len(apiWebcams), continent, total)

			batchSavedCount := 0
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

				// Check if a webcam with the same ID already exists.
				var existingWebcam storage.Webcam
				if err := db.First(&existingWebcam, dbWebcam.WebcamID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// Webcam does not exist, so create it.
						log.Printf("Creating new record for webcam ID: %d", dbWebcam.WebcamID)
						if createErr := db.Create(&dbWebcam).Error; createErr != nil {
							log.Printf("Failed to create webcam %d: %v", dbWebcam.WebcamID, createErr)
							continue
						}
						batchSavedCount++
					} else {
						// Another error occurred.
						log.Printf("Failed to query for webcam %d: %v", dbWebcam.WebcamID, err)
					}
				}
			}

			continentSavedCount += batchSavedCount
			offset += len(apiWebcams)

			if offset >= total {
				log.Printf("Finished fetching all %d webcams for continent: %s", total, continent)
				break
			}
		}
		log.Printf("Saved %d new webcams for continent %s.", continentSavedCount, continent)
		totalSavedCount += continentSavedCount
	}

	log.Printf("Successfully saved a total of %d new webcams from all continents.", totalSavedCount)
	log.Println("Windy Cams full data fetch finished successfully.")
}
