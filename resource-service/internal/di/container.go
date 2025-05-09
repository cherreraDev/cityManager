package di

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"resource-service/cmd/api/bootstrap/config"
	"resource-service/internal/resource/application"
	"resource-service/internal/resource/domain"
	"resource-service/internal/resource/platform/kafka/producer"
	"resource-service/internal/resource/platform/repositories/persistence_gorm"
	"sync"
	"time"
)

type Container struct {
	// Infrastructure
	Db       *gorm.DB
	Producer *producer.Producer

	// Repositories
	ResourceRepo domain.ResourceRepository

	//Services
	ResourceService domain.ResourceService
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

		prod := producer.NewProducer(cfg.Kafka.Brokers)

		// Repositories
		var resourceRepo domain.ResourceRepository = persistence_gorm.NewResourceRepository(db)

		// Services
		var resourceService domain.ResourceService = application.NewResourceService(resourceRepo, prod)

		container = &Container{
			Db:              db,
			Producer:        prod,
			ResourceService: resourceService,
			ResourceRepo:    resourceRepo,
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
			err = db.AutoMigrate(&persistence_gorm.ResourceModel{})
			if err != nil {
				return nil, fmt.Errorf("failed to auto-migrate ResourceModel: %w", err)
			}
			return db, nil
		}
		if time.Since(start) > timeout {
			return nil, fmt.Errorf("database not ready: %w", err)
		}
		fmt.Println("Waiting for database to be ready...")
		time.Sleep(2 * time.Second)
	}
}
