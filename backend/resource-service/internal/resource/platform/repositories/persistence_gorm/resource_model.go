package persistence_gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceModel struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	CityID uuid.UUID `gorm:"type:uuid;index"`
	Type   string
	Amount float64
}

// BeforeCreate sets a UUID if none is provided
func (r *ResourceModel) BeforeCreate(*gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
