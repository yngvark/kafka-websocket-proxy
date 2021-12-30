package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/pubsub"
	"go.uber.org/zap"
	"time"
)

type kafkaConsumer struct {
	log        *zap.SugaredLogger
	ctx        context.Context
	subscriber chan string
	reader     *kafka.Reader
}

func (c kafkaConsumer) ListenForMessages() error {
	for i := 0; i < 10; i++ {
		msg, err := c.reader.ReadMessage(c.ctx)
		if err != nil {
			return fmt.Errorf("reading message. Failed: %w", err)
		}

		fmt.Printf("message: %s\n", string(msg.Value))
		c.subscriber <- string(msg.Value)

		if string(msg.Value) == "/quit" {
			fmt.Println("Got /quit msg, quitting")
			break
		}
	}

	c.log.Info("Kafka reading done")

	return nil
}

func (c kafkaConsumer) SubscriberChannel() chan string {
	return c.subscriber
}

func (c kafkaConsumer) Close() error {
	return c.reader.Close()
}

// NewConsumer returns a KAFKA consumer
func NewConsumer(
	ctx context.Context,
	logger *zap.SugaredLogger,
	topic string,
	subscriber chan string,
) (pubsub.Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  100 * time.Millisecond,
	})

	return kafkaConsumer{
		log:        logger,
		ctx:        ctx,
		subscriber: subscriber,
		reader:     reader,
	}, nil
}
