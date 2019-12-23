package property

import (
	"context"

	entity "github.com/albertwidi/go-project-example/internal/entity/property"
	"github.com/albertwidi/go-project-example/internal/pkg/sqldb"
)

// Repository of property
type Repository struct {
	db *sqldb.DB
}

// New property repository
func New(db *sqldb.DB) *Repository {
	r := Repository{
		db: db,
	}
	return &r
}

// Create new property
func (r Repository) Create(ctx context.Context, property entity.Property, detail entity.Detail, addressMap entity.AddressMap, pricings []entity.Pricing) error {
	return nil
}

// Update property
func (r Repository) Update(ctx context.Context) error {
	return nil
}

// UpdateDetail of property
func (r Repository) UpdateDetail(ctx context.Context) error {
	return nil
}

// UpdateAddress update property address
func (r Repository) UpdateAddress(ctx context.Context) error {
	return nil
}

// Delete property
func (r Repository) Delete(ctx context.Context, id string) error {
	return nil
}
