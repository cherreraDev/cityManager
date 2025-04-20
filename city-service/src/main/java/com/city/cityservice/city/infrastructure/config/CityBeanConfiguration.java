package com.city.cityservice.city.infrastructure.config;

import com.city.cityservice.city.application.service.CityService;
import com.city.cityservice.city.domain.ports.output.CityEventPublisher;
import com.city.cityservice.city.domain.ports.output.CityPersistence;
import com.city.cityservice.city.domain.service.CityUseCase;
import com.city.cityservice.city.infrastructure.messaging.KafkaCityEventPublisher;
import com.city.cityservice.city.infrastructure.persistence.CityJpaRepository;
import com.city.cityservice.city.infrastructure.persistence.CityPersistenceAdapter;
import com.city.cityservice.city.infrastructure.persistence.CityPersistenceMapper;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.core.KafkaTemplate;

@Configuration
public class CityBeanConfiguration {
    @Bean
    public CityUseCase cityUseCase(
            CityPersistence cityPersistence,
            CityEventPublisher cityEventPublisher) {
        return new CityService(cityPersistence, cityEventPublisher);
    }

    @Bean
    public CityPersistence cityPersistence(
            CityJpaRepository cityJpaRepository,
            CityPersistenceMapper cityPersistenceMapper) {
        return new CityPersistenceAdapter(cityJpaRepository, cityPersistenceMapper);
    }

    @Bean
    public CityEventPublisher cityEventPublisher(KafkaTemplate<String, Object> kafkaTemplate) {
        return new KafkaCityEventPublisher(kafkaTemplate);
    }
}
