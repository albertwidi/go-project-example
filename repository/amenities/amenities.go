package amenities

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	entity "github.com/albertwidi/kothak/entity/amenities"
	"github.com/lib/pq"
)

// Repository of amenities
type Repository struct {
	db *sqlx.DB
}

// Amenities struct
type Amenities struct {
	ID        string      `db:"id"`
	Name      string      `db:"name"`
	Type      int         `db:"type"`
	ImagePath string      `db:"image_path"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at"`
	Deleted   bool        `db:"deleted"`
	IsTest    bool        `db:"is_test"`
}

// New amenities repo
func New() *Repository {
	r := Repository{}
	return &r
}

// Create a new amenities
func (r Repository) Create(ctx context.Context, amenities entity.Amenities) error {
	return nil
}

// Get amenities
func (r Repository) Get(ctx context.Context, amenitiesID ...string) ([]entity.Amenities, error) {
	a := []entity.Amenities{}
	return a, nil
}

// Update amenities
func (r Repository) Update(ctx context.Context, amenities entity.Amenities) error {
	return nil
}

// Delete the amenities
func (r Repository) Delete(ctx context.Context, amenitiesID string) error {
	return nil
}
