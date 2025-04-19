package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"resource-service/internal/resource/domain"
	"resource-service/internal/resource/platform/kafka/producer"
	"time"
)

type ResourceService struct {
	repo     domain.ResourceRepository
	producer *producer.Producer
}

func NewResourceService(repo domain.ResourceRepository, producer *producer.Producer) *ResourceService {
	return &ResourceService{repo: repo, producer: producer}
}

func (s *ResourceService) GetResources(cityID uuid.UUID) ([]domain.Resource, error) {
	return s.repo.FindByCityID(cityID)
}

func (s *ResourceService) ConsumeResources(ctx context.Context, cityID uuid.UUID, required map[string]float64) error {
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
		event := map[string]interface{}{
			"event_type": "resources_consumed",
			"city_id":    cityID,
			"resources":  required,
			"timestamp":  time.Now().UTC(),
		}

		key, err := cityID.MarshalBinary()
		if err != nil {
			return fmt.Errorf("failed to marshal cityID: %w", err)
		}

		if err := s.producer.SendJSON(ctx, "resource-updates", key, event); err != nil {
			return fmt.Errorf("failed to send kafka message: %w", err)
		}

		return nil
	})
}

func (s *ResourceService) InitializeCityResources(ctx context.Context, cityID uuid.UUID) error {
	initialResources := map[string]float64{
		"water":    0,
		"food":     0,
		"energy":   0,
		"minerals": 0,
	}
	if err := s.repo.InitializeResources(ctx, cityID, initialResources); err != nil {
		return fmt.Errorf("failed to initialize resources in DB: %w", err)
	}
	return nil
}
