package webcam

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/basel-ax/windy-cams/internal/domain"
)

// Repository defines the interface for platform data operations.
type Repository interface {
	SavePlatforms(ctx context.Context, platforms []domain.Platform) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewRepository creates a new platform repository.
func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

// SavePlatforms saves a list of platforms to the database.
// It uses an "upsert" operation to avoid duplicates.
func (r *postgresRepository) SavePlatforms(ctx context.Context, platforms []domain.Platform) error {
	if len(platforms) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "slug", "homepage"}),
		}).
		Create(&platforms).Error
	if err != nil {
		return fmt.Errorf("failed to save platforms: %w", err)
	}

	return nil
}
