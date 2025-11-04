package webcam

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/basel-ax/windy-cams/internal/domain"
)

// Repository defines the interface for webcam data operations.
type Repository interface {
	SaveWebcams(ctx context.Context, webcams []domain.Webcam) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewRepository creates a new platform repository.
func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

// SaveWebcams saves a list of webcams to the database.
// It uses an "upsert" operation to avoid duplicates.
func (r *postgresRepository) SaveWebcams(ctx context.Context, webcams []domain.Webcam) error {
	if len(webcams) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "status", "view_url"}),
		}).
		Create(&webcams).Error
	if err != nil {
		return fmt.Errorf("failed to save webcams: %w", err)
	}

	return nil
}
