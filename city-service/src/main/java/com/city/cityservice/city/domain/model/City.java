package com.city.cityservice.city.domain.model;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;

import java.time.LocalDateTime;
import java.util.UUID;

@Getter
@Builder
@AllArgsConstructor
public class City {
    private final UUID id;
    private final String name;
    private  Integer population;
    private final LocalDateTime createdAt;
    private final LocalDateTime updatedAt;

    public static City create(String name, Integer population) {
        return City.builder()
                .id(UUID.randomUUID())
                .name(name)
                .population(population)
                .createdAt(LocalDateTime.now())
                .updatedAt(LocalDateTime.now())
                .build();
    }

    public void changePopulation(Integer populationChange){
        this.population = this.population + populationChange;
    }
}


