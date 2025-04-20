package persistence_gorm

import "resource-service/internal/resource/domain"

func toDomain(model ResourceModel) domain.Resource {
	return domain.Resource{
		ID:     model.ID,
		CityID: model.CityID,
		Type:   model.Type,
		Amount: model.Amount,
	}
}

func toModel(resource domain.Resource) ResourceModel {
	return ResourceModel{
		ID:     resource.ID,
		CityID: resource.CityID,
		Type:   resource.Type,
		Amount: resource.Amount,
	}
}
