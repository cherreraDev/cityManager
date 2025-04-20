package com.city.cityservice.city.infrastructure.persistence;

import com.city.cityservice.city.domain.model.City;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface CityPersistenceMapper {

    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    CityJpaEntity toEntity(City city);

    City toDomain(CityJpaEntity entity);
}
