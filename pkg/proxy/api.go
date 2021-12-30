package proxy

import (
	"context"
	"fmt"
	"github.com/yngvark/kafka-websocket-proxy/pkg/connectors/kafka"
	"github.com/yngvark/kafka-websocket-proxy/pkg/connectors/websocket"
	"github.com/yngvark/kafka-websocket-proxy/pkg/connectors/websocket2"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/oslookup"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/pubsub"
	"go.uber.org/zap"
	"os"
)

func (p Proxy) Run() error {
	subscriber := make(chan string)

	// Handle Kafka
	kafkaPublisher, kafkaConsumer, err := pubSubForKafka(p.context, p.cancelFn, p.logger, subscriber)
	if err != nil {
		return fmt.Errorf("creating kafka connectors: %w", err)
	}

	// Handle websockets
	websocketListener := websocket2.NewListener("/v1/broker")
	websocketListener.Run()

	// TODO how to connect websockets with kafka

	// Close producer and consumer when done
	defer func() {
		var err error

		err = kafkaConsumer.Close()
		if err != nil {
			p.logger.Errorf("error closing kafka consumer: %s", err.Error())
		}

		err = kafkaPublisher.Close()
		if err != nil {
			p.logger.Errorf("error closing kafka publisher: %s", err.Error())
		}
	}()

	go func() {
		err2 := kafkaConsumer.ListenForMessages()
		if err2 != nil {
			p.logger.Errorf("Error listening for messages, stopping proxy. Details: %s", err.Error())
			p.cancelFn()
		}
	}()

	for {
		select {
		case msg := <-kafkaConsumer.SubscriberChannel():
			p.logger.Debugf("Received message: %s\n", msg)

			// TODO broadast msg to websockets
		case <-p.context.Done():
			p.logger.Info("Kafka proxy stopped.")
			return nil
		}
	}
}

const allowedCorsOriginsEnvVarKey = "ALLOWED_CORS_ORIGINS"

func pubSubForKafka(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriber chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	p, err := kafka.NewPublisher(ctx, cancelFn, logger, "zombie")
	if err != nil {
		return nil, nil, fmt.Errorf("creating publisher: %w", err)
	}

	c, err := kafka.NewConsumer(ctx, logger, "gameinit", subscriber)

	return p, c, nil
}

func pubSubForWebsocket(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriber chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	corsHelper := oslookup.NewCORSHelper(logger)

	allowedCorsOrigins, err := corsHelper.GetAllowedCorsOrigins(os.LookupEnv, allowedCorsOriginsEnvVarKey)
	if err != nil {
		return nil, nil, fmt.Errorf("getting allowed CORS origins: %w", err)
	}

	corsHelper.PrintAllowedCorsOrigins(allowedCorsOrigins)

	p, c := websocket.New(ctx, cancelFn, logger, subscriber, allowedCorsOrigins)

	return p, c, nil
}

func New(ctx context.Context, cancelFn context.CancelFunc, logger *zap.SugaredLogger) Proxy {
	return Proxy{
		context:  ctx,
		cancelFn: cancelFn,
		logger:   logger,
	}
}
