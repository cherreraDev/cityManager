package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,

		Async:        false,
		Compression:  kafka.Snappy,
		BatchTimeout: 10 * time.Millisecond,
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
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) SendJSON(ctx context.Context, topic string, key []byte, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.SendMessage(ctx, topic, key, jsonValue)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
