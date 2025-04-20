package persistence_gorm

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"resource-service/internal/resource/domain"
)

type ResourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

func (r *ResourceRepository) WithTransaction(fn func(r domain.ResourceRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &ResourceRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *ResourceRepository) FindByCityID(cityID uuid.UUID) ([]domain.Resource, error) {
	var models []ResourceModel
	if err := r.db.Where("city_id = ?", cityID).Find(&models).Error; err != nil {
		return nil, err
	}

	resources := make([]domain.Resource, len(models))
	for i, m := range models {
		resources[i] = toDomain(m)
	}

	return resources, nil
}

func (r *ResourceRepository) Save(resource *domain.Resource) error {
	model := toModel(*resource)
	return r.db.Save(&model).Error
}

func (r *ResourceRepository) InitializeResources(ctx context.Context, cityID uuid.UUID, resources map[string]float64) error {
	// Usar transacción para asegurar atomicidad
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for resourceType, amount := range resources {
			resource := &ResourceModel{
				ID:     uuid.New(),
				CityID: cityID,
				Type:   resourceType,
				Amount: amount,
			}
			if err := tx.Create(resource).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ResourceRepository) DeleteResources(ctx context.Context, cityID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("city_id = ?", cityID).Delete(&ResourceModel{}).Error
}
