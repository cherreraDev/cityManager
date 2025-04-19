package application

import (
	"fmt"
	"github.com/google/uuid"
	"resource-service/internal/resource/domain"
)

type ResourceService struct {
	repo domain.ResourceRepository
}

func NewResourceService(repo domain.ResourceRepository) *ResourceService {
	return &ResourceService{repo: repo}
}

func (s *ResourceService) GetResources(cityID uuid.UUID) ([]domain.Resource, error) {
	return s.repo.FindByCityID(cityID)
}

func (s *ResourceService) ConsumeResources(cityID uuid.UUID, required map[string]float64) error {
	return s.repo.WithTransaction(func(txRepo domain.ResourceRepository) error {
		resources, err := txRepo.FindByCityID(cityID)
		if err != nil {
			return fmt.Errorf("failed to fetch resources: %w", err)
		}

		resourceMap := make(map[string]*domain.Resource)
		for i := range resources {
			r := &resources[i]
			resourceMap[r.Type] = r
		}

		for typ, amount := range required {
			res, ok := resourceMap[typ]
			if !ok {
				return fmt.Errorf("resource type %s not found", typ)
			}
			if err := res.Consume(amount); err != nil {
				return fmt.Errorf("failed to consume resource %s: %w", typ, err)
			}
		}

		for _, res := range resourceMap {
			if err := txRepo.Save(res); err != nil {
				return fmt.Errorf("failed to persist resource %s: %w", res.Type, err)
			}
		}

		return nil
	})
}
