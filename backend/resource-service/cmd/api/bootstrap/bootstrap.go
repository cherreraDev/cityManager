package bootstrap

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/di"
	"resource-service/internal/resource/platform/kafka/consumer"
	"resource-service/internal/resource/platform/server"
	"sync"
)

func Run() error {
	// Carga de configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Inicialización de dependencias
	container, err := di.InitializeContainer(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}

	// Creación del servidor HTTP con contexto
	ctx, cancel, srv := server.NewServer(context.Background(), cfg.HTTP.Host, cfg.HTTP.Port, cfg.HTTP.ShutdownTimeout)
	defer cancel()

	// Grupo de espera para consumidores
	var wg sync.WaitGroup

	// Inicio de consumidores Kafka
	if err := consumer.StartConsumers(container, &cfg, ctx, &wg); err != nil {
		return fmt.Errorf("failed to start kafka consumers: %w", err)
	}

	// Canal para errores del servidor
	serverErr := make(chan error, 1)

	// Ejecución del servidor
	go func() {
		log.Println("Starting HTTP server...")
		serverErr <- srv.Run(ctx, container)
	}()

	select {
	case err := <-serverErr:
		wg.Wait()
		closeResources(container)
		return err

	case <-ctx.Done():
		log.Println("Starting graceful shutdown...")
		wg.Wait()
		closeResources(container)
		log.Println("Shutdown completed")
		return nil
	}
}

func closeResources(container *di.Container) {
	if container.Producer != nil {
		log.Println("Closing Kafka producer...")
		if err := container.Producer.Close(); err != nil {
			log.Printf("Warning: producer close error: %v", err)
		}
	}

	if container.Db != nil {
		if db, err := container.Db.DB(); err == nil {
			err := db.Close()
			if err != nil {
				log.Printf("Warning: db close error: %v", err)
			}
		}
	}
}
