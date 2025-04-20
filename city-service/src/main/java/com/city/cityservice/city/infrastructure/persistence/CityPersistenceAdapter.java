package com.city.cityservice.city.infrastructure.persistence;

import com.city.cityservice.city.domain.model.City;
import com.city.cityservice.city.domain.ports.output.CityPersistence;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@Component
@AllArgsConstructor
public class CityPersistenceAdapter implements CityPersistence {
    private final CityJpaRepository cityJpaRepository;
    private final CityPersistenceMapper cityPersistenceMapper;

    @Override
    public City saveCity(City city) {
        CityJpaEntity entity = cityPersistenceMapper.toEntity(city);
        CityJpaEntity savedEntity = cityJpaRepository.save(entity);
        return cityPersistenceMapper.toDomain(savedEntity);
    }

    @Override
    public List<City> getAllCities() {
        return cityJpaRepository.findAll().stream()
                .map(cityPersistenceMapper::toDomain)
                .collect(Collectors.toList());
    }

    @Override
    public Optional<City> getCityById(UUID id) {
        return cityJpaRepository.findById(id)
                .map(cityPersistenceMapper::toDomain);
    }

    @Override
    public void deleteCityById(UUID id) {
        cityJpaRepository.deleteById(id);
    }

    public boolean existsByName(String name) {
        return cityJpaRepository.existsByName(name);
    }

}
