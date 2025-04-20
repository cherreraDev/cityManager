package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/di"
	"sync"
	"time"
)

func startCityConsumer(container *di.Container, config *config.Config, ctx context.Context, wg *sync.WaitGroup) error {
	consumerConfig := ConsumerConfig{
		Brokers:  config.Kafka.Brokers,
		GroupID:  "resource-service-city",
		Topic:    config.Kafka.Topics.CityEvents,
		MinBytes: 1e3, // 1KB
		MaxBytes: 1e6, // 1MB
	}

	consumer := NewConsumer(consumerConfig)

	processMessage := func(key []byte, value any) error {
		type CityEvent struct {
			EventType string             `json:"event_type"`
			CityID    uuid.UUID          `json:"city_id"`
			Resources map[string]float64 `json:"resources"`
			Timestamp time.Time          `json:"timestamp"`
		}

		var event CityEvent
		if err := json.Unmarshal(value.([]byte), &event); err != nil {
			return fmt.Errorf("error unmarshaling city event: %v", err)
		}

		// Only process city creation events
		if event.EventType != "city_created" {
			return nil
		}
		return container.ResourceService.InitializeCityResources(context.Background(), event.CityID)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := consumer.Close(); err != nil {
				log.Printf("Error closing city consumer: %v", err)
			}
		}()

		var msgData struct {
			EventType string             `json:"event_type"`
			CityID    uuid.UUID          `json:"city_id"`
			Resources map[string]float64 `json:"resources"`
			Timestamp time.Time          `json:"timestamp"`
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping city consumer...")
				return
			default:
				err := consumer.ConsumeJSON(ctx, &msgData, processMessage)
				if err != nil {
					log.Printf("Error in city consumer: %v. Retrying in 5 seconds...", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}
		}
	}()

	log.Println("City consumer started successfully")
	return nil
}
