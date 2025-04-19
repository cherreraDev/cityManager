package domain

import "github.com/google/uuid"

type ResourceService interface {
	ConsumeResources(cityID uuid.UUID, resources map[string]float64) error
	GetResources(cityID uuid.UUID) ([]Resource, error)
}
