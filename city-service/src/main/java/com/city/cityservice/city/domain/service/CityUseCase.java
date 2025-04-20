package com.city.cityservice.city.domain.service;

import com.city.cityservice.city.domain.model.City;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface CityUseCase {
    City createCity(City city);
    List<City> getAllCities();
    Optional<City> getCityById(UUID id);
    void deleteCity(UUID id);
}
