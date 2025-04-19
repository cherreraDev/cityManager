package consumer

import (
	"context"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/di"
	"sync"
)

func StartConsumers(container *di.Container, config *config.Config, ctx context.Context, wg *sync.WaitGroup) error {
	// Consumer for resource processing
	err := startResourceConsumer(container, config, ctx, wg)
	if err != nil {
		return err
	}

	return nil
}
