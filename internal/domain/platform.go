package domain

import "gorm.io/gorm"

// Webcam represents a webcam entity in the database.
type Webcam struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Status    string         `json:"status"`
	ViewURL   string         `json:"view_url"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Adds soft delete support
}
