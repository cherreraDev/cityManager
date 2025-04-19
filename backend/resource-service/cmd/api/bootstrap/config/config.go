package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Kafka struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
		GroupID string   `envconfig:"KAFKA_GROUPID" default:"resource-service"`
		Topics  struct {
			ResourceUpdates string `envconfig:"TOPIC_RESOURCE_UPDATES" default:"resource-updates"`
			CityEvents      string `envconfig:"TOPIC_CITY_EVENTS" default:"city-events"`
		}
	}
	HTTP struct {
		Port            string        `envconfig:"PORT" default:"8080"`
		Host            string        `envconfig:"HOST" default:"127.0.0.1"`
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
	}
	DB struct {
		UserName string        `envconfig:"DB_USERNAME" default:"user"`
		Password string        `envconfig:"DB_PASSWORD" default:"userpassword"`
		Host     string        `envconfig:"DB_HOST" default:"localhost"`
		Port     string        `envconfig:"DB_PORT" default:"3306"`
		Name     string        `envconfig:"DB_NAME" default:"mydb"`
		Timeout  time.Duration `envconfig:"DB_TIMEOUT" default:"10s"`
	}
	SecretKey   string `envconfig:"SECRET_KEY" default:"secretos123"`
	Environment string `envconfig:"ENVIRONMENT" default:"dev"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	cfg = Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
