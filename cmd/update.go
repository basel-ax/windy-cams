package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/basel-ax/windy-cams/pkg/config"
	"github.com/basel-ax/windy-cams/pkg/storage"
	"github.com/basel-ax/windy-cams/pkg/windy"

	"gorm.io/datatypes"
)

func runUpdate(devMode bool) {
	log.Println("Starting to update all Windy Cams data...")

	// Load configuration
	cfg := config.New()
	log.Println("Configuration loaded successfully.")

	// Initialize database
	db := storage.New(cfg, devMode)
	log.Println("Database connection established.")

	// Create Windy API client
	windyClient := windy.NewClientWithDevMode(cfg.WindyAPIKey, devMode)
	log.Println("Windy API client created.")

	// 1. Fetch all webcams from the database
	var webcams []storage.Webcam
	if err := db.Find(&webcams).Error; err != nil {
		log.Fatalf("Failed to fetch webcams from database: %v", err)
	}

	totalWebcams := len(webcams)
	log.Printf("Found %d webcams to update.", totalWebcams)

	updatedCount := 0
	for i, webcam := range webcams {
		log.Printf("[%d/%d] Fetching details for webcam ID: %d", i+1, totalWebcams, webcam.WebcamID)

		// 2. Fetch details for each webcam
		details, err := windyClient.GetWebcamDetails(webcam.WebcamID)
		if err != nil {
			log.Printf("Failed to get details for webcam %d: %v", webcam.WebcamID, err)
			continue // Skip to the next webcam
		}

		// 3. Marshal nested JSON objects
		categoriesJSON, err := json.Marshal(details.Categories)
		if err != nil {
			log.Printf("Failed to marshal categories for webcam %d: %v", webcam.WebcamID, err)
			continue
		}

		playerJSON, err := json.Marshal(details.Player)
		if err != nil {
			log.Printf("Failed to marshal player for webcam %d: %v", webcam.WebcamID, err)
			continue
		}

		urlsJSON, err := json.Marshal(details.Urls)
		if err != nil {
			log.Printf("Failed to marshal urls for webcam %d: %v", webcam.WebcamID, err)
			continue
		}

		// 4. Update the webcam record in the database
		result := db.Model(&webcam).Updates(storage.Webcam{
			Title:         details.Title,
			ApiStatus:     details.Status,
			ViewCount:     details.ViewCount,
			LastUpdatedOn: details.LastUpdatedOn,
			Latitude:      details.Location.Latitude,
			Longitude:     details.Location.Longitude,
			City:          details.Location.City,
			Region:        details.Location.Region,
			RegionCode:    details.Location.RegionCode,
			Country:       details.Location.Country,
			CountryCode:   details.Location.CountryCode,
			Continent:     details.Location.Continent,
			ContinentCode: details.Location.ContinentCode,
			Categories:    datatypes.JSON(categoriesJSON),
			Player:        datatypes.JSON(playerJSON),
			Urls:          datatypes.JSON(urlsJSON),
			UpdatedAt:     time.Now(),
		})

		if result.Error != nil {
			log.Printf("Failed to update webcam %d in database: %v", webcam.WebcamID, result.Error)
			continue
		}

		log.Printf("Successfully updated webcam ID: %d", webcam.WebcamID)
		updatedCount++
	}

	log.Printf("Update process finished. Successfully updated %d out of %d webcams.", updatedCount, totalWebcams)
}
