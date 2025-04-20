package consumer

import (
	"context"
	"fmt"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/di"
	"sync"
)

func StartConsumers(container *di.Container, config *config.Config, ctx context.Context, wg *sync.WaitGroup) error {
	// Start resource updates consumer
	if err := startResourceConsumer(container, config, ctx, wg); err != nil {
		return fmt.Errorf("failed to start resource consumer: %w", err)
	}

	// Start city events consumer
	if err := startCityConsumer(container, config, ctx, wg); err != nil {
		return fmt.Errorf("failed to start city consumer: %w", err)
	}

	return nil
}
