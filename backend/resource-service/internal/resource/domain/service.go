package domain

import (
	"context"
	"github.com/google/uuid"
)

type ResourceService interface {
	ConsumeResources(ctx context.Context, cityID uuid.UUID, resources map[string]float64) error
	GetResources(cityID uuid.UUID) ([]Resource, error)
}
