package com.city.cityservice.city.application.service;

import com.city.cityservice.city.domain.model.City;
import com.city.cityservice.city.domain.ports.output.CityEventPublisher;
import com.city.cityservice.city.domain.ports.output.CityPersistence;
import com.city.cityservice.city.domain.service.CityUseCase;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@AllArgsConstructor
@Service
public class CityService implements CityUseCase {
    private final CityPersistence cityPersistence;
    private final CityEventPublisher cityEventPublisher;
    @Override
    public City createCity(City city) {
        City savedCity = cityPersistence.saveCity(city);
        cityEventPublisher.publishCityCreatedEvent(savedCity);
        return savedCity;
    }

    @Override
    public List<City> getAllCities() {
        return cityPersistence.getAllCities();
    }

    @Override
    public Optional<City> getCityById(UUID id) {
        return cityPersistence.getCityById(id);
    }

    @Override
    public void deleteCity(UUID id) {
        cityPersistence.deleteCityById(id);
    }
}
