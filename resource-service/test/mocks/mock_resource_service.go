package mocks

import (
	"context"
	"github.com/google/uuid"
	"resource-service/internal/resource/domain"
)

type MockResourceService struct {
	GetResourcesFn            func(cityID uuid.UUID) ([]domain.Resource, error)
	ConsumeResourcesFn        func(ctx context.Context, cityID uuid.UUID, resources map[string]float64) error
	InitializeCityResourcesFn func(ctx context.Context, cityId uuid.UUID) error
}

func (m *MockResourceService) GetResources(cityID uuid.UUID) ([]domain.Resource, error) {
	return m.GetResourcesFn(cityID)
}

func (m *MockResourceService) ConsumeResources(ctx context.Context, cityID uuid.UUID, resources map[string]float64) error {
	return m.ConsumeResourcesFn(ctx, cityID, resources)
}
func (m *MockResourceService) InitializeCityResources(ctx context.Context, cityId uuid.UUID) error {
	return m.InitializeCityResourcesFn(ctx, cityId)
}
