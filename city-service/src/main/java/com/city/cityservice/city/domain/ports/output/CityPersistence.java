package com.city.cityservice.city.domain.ports.output;

import com.city.cityservice.city.domain.model.City;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface CityPersistence {
    City saveCity(City city);
    List<City> getAllCities();
    Optional<City> getCityById(UUID id);
    void deleteCityById(UUID id);
}
