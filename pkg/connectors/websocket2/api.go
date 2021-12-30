package websocket2

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

func (h *HTTPHandler) Run() {
	http.HandleFunc("/v1/broker", h.ServeHTTPActual)
}

func (h *HTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	connection, err := h.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Error("could not upgrade:", err)
		return
	}

}

func (h *HTTPHandler) ServeHTTPActual(writer http.ResponseWriter, request *http.Request) {
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

func NewListener() *HTTPHandler {
	return &HTTPHandler{
		cancelFn: nil,
		logger:   nil,
		upgrader: &websocket.Upgrader{
			CheckOrigin:       createWebsocketCheckOriginFn(logger, allowedOrigins),
			EnableCompression: true,
		},
		connection:  nil,
		subscriber:  nil,
		broadcaster: nil,
	}
}

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
