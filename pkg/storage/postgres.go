package storage

import (
	"fmt"
	"log"

	"github.com/basel-ax/windy-cams/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Webcam defines the database model for a webcam.
type Webcam struct {
	gorm.Model
	WebcamID  uint64 `gorm:"uniqueIndex"`
	ApiStatus string
	Status    string
	Title     string `gorm:"index"`
	Latitude  float64
	Longitude float64
	City      string
	Country   string
	Continent string
}

// New initializes a new PostgreSQL database connection.
func New(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the Webcam schema
	if err := db.AutoMigrate(&Webcam{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
