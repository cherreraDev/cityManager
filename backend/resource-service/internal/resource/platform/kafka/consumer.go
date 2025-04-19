package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type ConsumerConfig struct {
	Brokers  []string
	GroupID  string
	Topic    string
	MinBytes int
	MaxBytes int
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(config ConsumerConfig) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		GroupID:     config.GroupID,
		Topic:       config.Topic,
		MinBytes:    config.MinBytes,
		MaxBytes:    config.MaxBytes,
		StartOffset: kafka.FirstOffset,
		MaxWait:     100 * time.Millisecond,
	})
	return &Consumer{reader: r}
}

func (c *Consumer) Consume(ctx context.Context, handler func(message kafka.Message) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return err
			}
			if err := handler(msg); err != nil {
				log.Printf("Error procesing message: %v", err)
			}
		}
	}
}

func (c *Consumer) ConsumeJSON(ctx context.Context, v any, handler func(key []byte, value any) error) error {
	return c.Consume(ctx, func(msg kafka.Message) error {
		if err := json.Unmarshal(msg.Value, &v); err != nil {
			return err
		}
		return handler(msg.Key, v)
	})
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
