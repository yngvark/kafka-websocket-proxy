// Package websocket provides publish subscribe using websockets
package websocket

import (
	"github.com/yngvark/kafka-websocket-proxy/pkg/connectors/websocket/httphandler"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/pubsub"
	"go.uber.org/zap"
)

type websocketPublisher struct {
	logger      *zap.SugaredLogger
	httpHandler *httphandler.HTTPHandler
}

// SendMsg sends messages
func (p websocketPublisher) SendMsg(msg string) error {
	return p.httpHandler.SendMsgToConnection(msg)
}

// Close closes the publisher
func (p websocketPublisher) Close() error {
	p.logger.Info("websocketPublisher websocketConsumer")

	if p.httpHandler != nil {
		return p.httpHandler.Close()
	}

	return nil
}

// NewPublisher returns a new publisher for websockets
func newPublisher(logger *zap.SugaredLogger, httpHandler *httphandler.HTTPHandler) pubsub.Publisher {
	return websocketPublisher{
		logger:      logger,
		httpHandler: httpHandler,
	}
}
