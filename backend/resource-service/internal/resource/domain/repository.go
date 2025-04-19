package domain

import "github.com/google/uuid"

type ResourceRepository interface {
	FindByCityID(cityID uuid.UUID) ([]Resource, error)
	Save(resource *Resource) error
}
