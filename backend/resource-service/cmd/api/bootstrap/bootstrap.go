package bootstrap

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/di"
	"resource-service/internal/resource/platform/kafka"
	"resource-service/internal/resource/platform/server"
	"sync"
)

func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	container, err := di.InitializeContainer(cfg)
	if err != nil {
		return err
	}

	ctx, srv := server.NewServer(context.Background(), cfg.HTTP.Host, cfg.HTTP.Port, cfg.HTTP.ShutdownTimeout)
	var wg sync.WaitGroup
	err = kafka.StartConsumers(container, &cfg, ctx, &wg)
	if err != nil {
		return err
	}
	return srv.Run(ctx, container)
}
