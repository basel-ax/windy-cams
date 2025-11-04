package domain

import "gorm.io/gorm"

// Platform represents an asset platform from the Windy Webcams API.
type Platform struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Homepage  string         `json:"homepage"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Adds soft delete support
}
