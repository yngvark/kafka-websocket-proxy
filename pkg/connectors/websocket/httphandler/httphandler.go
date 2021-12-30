// Package httphandler knows how to handle HTTP websocket connections
package httphandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/yngvark/kafka-websocket-proxy/pkg/connectors/websocket/httphandler/broadcast"
	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

// HTTPHandler knows how to handle HTTP websocket connections
type HTTPHandler struct {
	cancelFn    context.CancelFunc
	logger      *zap.SugaredLogger
	upgrader    *websocket.Upgrader
	connection  *websocket.Conn
	subscriber  chan string
	broadcaster *broadcast.Broadcaster
}

func (h *HTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	connection, err := h.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Error("could not upgrade:", err)
		return
	}

	// We only support one client
	h.connection = connection

	h.logger.Info("Client connected!")

	// Handle disconnection
	activelyCloseConnectionChannel := make(chan bool)

	defer h.closeConnectionWhenDone(activelyCloseConnectionChannel)

	go h.readIncomingMessages(activelyCloseConnectionChannel)
}

func (h *HTTPHandler) readIncomingMessages(closeConnectionChannel chan bool) {
	for {
		h.logger.Info("Reading next message...")

		_, message, err := h.connection.ReadMessage()
		if err != nil {
			// Client disconnected
			h.logger.Info("Client disconnected")

			// We need to stop both game logic and disconnect
			closeConnectionChannel <- true

			closeError, ok := err.(*websocket.CloseError)
			if ok {
				h.logger.Infof("Client disconnected OK. Code: %d", closeError.Code)
			} else {
				h.logger.Errorf("Read error: %s", err.Error())
			}

			return
		}

		h.logger.Infof("Broadcasting received message: %s", message)

		err = h.broadcast(message)
		if err != nil {
			h.logger.Errorf("Error handling incoming message. Aborting. Error: %s", err.Error())

			closeConnectionChannel <- true

			h.cancelFn()

			return
		}
	}
}

func (h *HTTPHandler) closeConnectionWhenDone(closeConnectionChannel chan bool) {
	<-closeConnectionChannel

	h.logger.Info("Closing connection from server")

	err := h.Close()

	if err != nil {
		h.logger.Info("error when closing connection: %w", err)
	} else {
		h.logger.Info("Connection closed successfully.")
	}
}

func (h *HTTPHandler) broadcast(message []byte) error {
	msgString := string(message)

	err := h.broadcaster.BroadCast(msgString)
	if err != nil {
		return fmt.Errorf("sending message with publisher: %w", err)
	}

	return nil
}

// SendMsgToConnection sends a message via the websocket
func (h *HTTPHandler) SendMsgToConnection(msg string) error {
	if h.connection == nil {
		return errors.New("could not send message, not connected")
	}

	h.logger.Infof("Sending msg: %s", msg)

	err := h.connection.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	return nil
}

// Close closes the handler
func (h *HTTPHandler) Close() error {
	h.logger.Info("Closing HTTPHandler")

	if h.connection != nil {
		return h.connection.Close()
	}

	return nil
}

type webfn func(r *http.Request) bool

func createWebsocketCheckOriginFn(logger *zap.SugaredLogger, allowedOrigins map[string]bool) webfn {
	return func(r *http.Request) bool {
		origin, ok := r.Header["Origin"]
		if !ok {
			return false
		}

		if len(origin) > 0 {
			_, ok := allowedOrigins[origin[0]]
			logger.Infof("Checking origin %s. Result: %t\n", origin[0], ok)

			return ok
		}

		return true
	}
}

// New returns a new HTTPHandler
func New(
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	allowedOrigins map[string]bool,
	subscriber chan string,
) *HTTPHandler {
	broadcaster := broadcast.New()
	broadcaster.AddSubscriber(subscriber)

	handler := &HTTPHandler{
		cancelFn:    cancelFn,
		logger:      logger,
		broadcaster: broadcaster,
		subscriber:  subscriber,
	}

	handler.upgrader = &websocket.Upgrader{
		CheckOrigin:       createWebsocketCheckOriginFn(logger, allowedOrigins),
		EnableCompression: true,
	}

	return handler
}
