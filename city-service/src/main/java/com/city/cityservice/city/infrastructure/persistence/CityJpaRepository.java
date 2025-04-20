package com.city.cityservice.city.infrastructure.persistence;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface CityJpaRepository  extends JpaRepository<CityJpaEntity, UUID> {
    boolean existsByName(String name);
}
