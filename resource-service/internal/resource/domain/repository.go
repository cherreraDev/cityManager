package domain

import (
	"context"
	"github.com/google/uuid"
)

type ResourceRepository interface {
	FindByCityID(cityID uuid.UUID) ([]Resource, error)
	Save(resource *Resource) error
	WithTransaction(fn func(r ResourceRepository) error) error

	InitializeResources(ctx context.Context, cityID uuid.UUID, resources map[string]float64) error
	DeleteResources(ctx context.Context, cityID uuid.UUID) error
}
