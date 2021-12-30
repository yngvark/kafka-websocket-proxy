package websocket2

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type HTTPHandler struct {
	cancelFn    context.CancelFunc
	logger      *zap.SugaredLogger
	upgrader    *websocket.Upgrader
	connection  *websocket.Conn
	subscriber  chan string
	broadcaster *broadcast.Broadcaster
}

type webfn func(r *http.Request) bool
