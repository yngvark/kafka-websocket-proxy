// Package kafka handles publishing and subscribing with Kafka
package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/pubsub"
	"go.uber.org/zap"
	"time"
)

type kafkaPublisher struct {
	log      *zap.SugaredLogger
	ctx      context.Context
	cancelFn context.CancelFunc
	conn     *kafka.Conn
}

func (p kafkaPublisher) SendMsg(msg string) error {
	_, err := p.conn.WriteMessages(kafka.Message{
		Value: []byte(msg),
	})
	if err != nil {
		p.cancelFn()
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (p kafkaPublisher) Close() error {
	return p.conn.Close()
}

// NewPublisher returns a kafka publisher
func NewPublisher(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	topic string,
) (pubsub.Publisher, error) {
	config := kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  100 * time.Millisecond,
	}

	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", config.Topic, config.Partition)
	if err != nil {
		return nil, fmt.Errorf("connecting to Kafka: %w", err)
	}

	return kafkaPublisher{
		log:      logger,
		ctx:      ctx,
		cancelFn: cancelFn,
		conn:     conn,
	}, nil
}
