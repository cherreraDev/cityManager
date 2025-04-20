package com.city.cityservice.city.infrastructure.messaging;

import com.city.cityservice.city.domain.model.City;
import com.city.cityservice.city.domain.ports.output.CityEventPublisher;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;

import java.util.HashMap;
import java.util.Map;
@RequiredArgsConstructor
public class KafkaCityEventPublisher implements CityEventPublisher {
    private final KafkaTemplate<String, Object> kafkaTemplate;
    @Value("${app.kafka.topics.city-events}")
    private String cityEventsTopic;
    @Override
    public void publishCityCreatedEvent(City city) {
        Map<String, Object> event = new HashMap<>();
        event.put("eventType", "city_created");
        event.put("cityId", city.getId());
        event.put("name", city.getName());
        event.put("timestamp", System.currentTimeMillis());

        kafkaTemplate.send(cityEventsTopic, city.getId().toString(), event);
    }
}
