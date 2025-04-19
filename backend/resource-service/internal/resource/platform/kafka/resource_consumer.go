package kafka

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

func startResourceConsumer(container *di.Container, config *config.Config, ctx context.Context, wg *sync.WaitGroup) error {
	consumerConfig := ConsumerConfig{
		Brokers:  config.Kafka.Brokers,
		GroupID:  "resource-service-group",
		Topic:    config.Kafka.Topics.ResourceUpdates,
		MinBytes: 1e3, // 1KB
		MaxBytes: 1e6, // 1MB
	}

	consumer := NewConsumer(consumerConfig)

	// Message processing function
	processMessage := func(key []byte, value any) error {
		// Parse the message into the expected format
		// Assuming the message contains cityID and required resources
		type ResourceRequest struct {
			CityID   uuid.UUID          `json:"city_id"`
			Required map[string]float64 `json:"required"`
		}

		var request ResourceRequest
		if err := json.Unmarshal(key, &request); err != nil {
			return fmt.Errorf("error unmarshaling resource request: %v", err)
		}

		// Call the service method
		return container.ResourceService.ConsumeResources(request.CityID, request.Required)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func(consumer *Consumer) {
			err := consumer.Close()
			if err != nil {

			}
		}(consumer)

		var msgData struct {
			CityID   uuid.UUID          `json:"city_id"`
			Required map[string]float64 `json:"required"`
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping resource consumer...")
				return
			default:
				err := consumer.ConsumeJSON(ctx, &msgData, processMessage)
				if err != nil {
					log.Printf("Error in resource consumer: %v. Retrying in 5 seconds...", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}
		}
	}()

	log.Println("Resource consumer started successfully")
	return nil
}
