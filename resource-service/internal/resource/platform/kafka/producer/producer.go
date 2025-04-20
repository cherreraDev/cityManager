package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

type ProducerConfig struct {
	Brokers       []string
	RequiredAcks  kafka.RequiredAcks
	Async         bool
	Compression   kafka.Compression
	BatchTimeout  time.Duration
	BatchSize     int
	QueueCapacity int
}

func NewProducer(brokers []string) *Producer {
	return NewProducerWithConfig(ProducerConfig{
		Brokers:      brokers,
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		Compression:  kafka.Snappy,
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	})
}

func NewProducerWithConfig(config ProducerConfig) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: config.RequiredAcks,
		Async:        config.Async,
		Compression:  config.Compression,
		BatchTimeout: config.BatchTimeout,
		BatchSize:    config.BatchSize,
	}
	return &Producer{writer: w}
}

func (p *Producer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	// Retry mechanism (simple)
	var lastErr error
	for i := 0; i < 3; i++ {
		if err := p.writer.WriteMessages(ctx, msg); err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed after 3 attempts: %w", lastErr)
}

func (p *Producer) SendJSON(ctx context.Context, topic string, key []byte, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}
	return p.SendMessage(ctx, topic, key, jsonValue)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
