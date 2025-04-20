package com.city.cityservice.city.domain.ports.output;

import com.city.cityservice.city.domain.model.City;

public interface CityEventPublisher {
    void publishCityCreatedEvent(City city);
}
