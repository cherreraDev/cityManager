package mocks

import (
	"github.com/google/uuid"
	"resource-service/internal/resource/domain"
)

type MockResourceService struct {
	GetResourcesFn     func(cityID uuid.UUID) ([]domain.Resource, error)
	ConsumeResourcesFn func(cityID uuid.UUID, resources map[string]float64) error
}

func (m *MockResourceService) GetResources(cityID uuid.UUID) ([]domain.Resource, error) {
	return m.GetResourcesFn(cityID)
}

func (m *MockResourceService) ConsumeResources(cityID uuid.UUID, resources map[string]float64) error {
	return m.ConsumeResourcesFn(cityID, resources)
}
