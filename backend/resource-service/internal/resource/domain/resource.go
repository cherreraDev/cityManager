package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Resource struct {
	ID     uuid.UUID `json:"id"`
	CityID uuid.UUID `json:"city_id"`
	Type   string    `json:"type"`
	Amount float64   `json:"amount"`
}

func (r *Resource) Consume(quantity float64) error {
	if r.Amount < quantity {
		return fmt.Errorf("insufficient resources")
	}
	r.Amount -= quantity
	return nil
}
