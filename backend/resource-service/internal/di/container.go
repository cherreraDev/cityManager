package di

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/resource/domain"
	"resource-service/internal/resource/platform/kafka"
	"sync"
	"time"
)

type Container struct {
	// Infrastructure
	producer *kafka.Producer
	db       *gorm.DB

	// Repositories
	resourceRepo domain.ResourceRepository

	//Services
	resourceService domain.ResourceService
}

var (
	container *Container
	initErr   error
	once      sync.Once
)

func InitializeContainer(cfg config.Config) (*Container, error) {
	once.Do(func() {
		// Database setup
		DBCfg := cfg.DB
		dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBCfg.UserName, DBCfg.Password, DBCfg.Host, DBCfg.Port, DBCfg.Name)

		db, err := waitForDB(dbURI, DBCfg.Timeout)
		if err != nil {
			initErr = fmt.Errorf("failed to init database: %w", err)
			return
		}

		//Producer
		producer := kafka.NewProducer(cfg.Kafka.Brokers)

		//Repositories
		//ToDo
		//Services
		//ToDo

		container = &Container{
			producer: producer,
			db:       db,
		}
	})
	return container, initErr
}

func GetContainer() (*Container, error) {
	if container == nil && initErr != nil {
		return nil, fmt.Errorf("container initialization failed: %w", initErr)
	}
	if container == nil {
		return nil, errors.New("container not initialized")
	}
	return container, nil
}

func waitForDB(dbURI string, timeout time.Duration) (*gorm.DB, error) {
	start := time.Now()
	for {
		db, err := gorm.Open(mysql.Open(dbURI))
		if err == nil {
			return db, nil
		}
		if time.Since(start) > timeout {
			return nil, fmt.Errorf("database not ready: %w", err)
		}
		fmt.Println("Waiting for database to be ready...")
		time.Sleep(2 * time.Second)
	}
}
