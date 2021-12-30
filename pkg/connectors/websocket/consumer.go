package websocket

import (
	"context"
	"errors"
	"net/http"

	"github.com/yngvark/kafka-websocket-proxy/pkg/connectors/websocket/httphandler"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/pubsub"
	"go.uber.org/zap"
)

type websocketConsumer struct {
	ctx        context.Context
	logger     *zap.SugaredLogger
	subscriber chan string

	httpHandler *httphandler.HTTPHandler
	listening   bool
}

// ListenForMessages starts to receive messages which will be available by reading SubscriberChannel(). It blocks
// until the websocketConsumer's context is canceled, so you should start it as a goroutine.
func (c *websocketConsumer) ListenForMessages() error {
	if !c.listening {
		c.listening = true
	} else {
		return errors.New("already listening for messages. Can listen for messages only once")
	}

	http.Handle("/zombie", c.httpHandler)

	<-c.ctx.Done()

	return nil
}

// SubscriberChannel returns a channel which can be used for reading incoming messages
func (c *websocketConsumer) SubscriberChannel() chan string {
	return c.subscriber
}

// Close closes the websocketConsumer
func (c *websocketConsumer) Close() error {
	c.logger.Info("Closing websocketConsumer")

	if c.httpHandler != nil {
		return c.httpHandler.Close()
	}

	return nil
}

// NewConsumer returns a new consumer for websockets
func newConsumer(
	ctx context.Context,
	logger *zap.SugaredLogger,
	subscriber chan string,
	httphandler *httphandler.HTTPHandler,
) pubsub.Consumer {
	return &websocketConsumer{
		ctx:         ctx,
		logger:      logger,
		subscriber:  subscriber,
		httpHandler: httphandler,
	}
}
