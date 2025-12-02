package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/basel-ax/windy-cams/pkg/config"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Webcam defines the database model for a webcam.
type Webcam struct {
	WebcamID      uint64 `gorm:"primaryKey"`
	ApiStatus     string
	Status        string
	Title         string `gorm:"index"`
	ViewCount     int
	LastUpdatedOn time.Time
	Latitude      float64
	Longitude     float64
	City          string
	Region        string
	RegionCode    string
	Country       string
	CountryCode   string
	Continent     string
	ContinentCode string
	Categories    datatypes.JSON
	Player        datatypes.JSON
	Urls          datatypes.JSON
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// New initializes a new PostgreSQL database connection.
func New(cfg *config.Config, devMode bool) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	logLevel := logger.Silent
	if devMode {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: !devMode,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the Webcam schema
	if err := db.AutoMigrate(&Webcam{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
