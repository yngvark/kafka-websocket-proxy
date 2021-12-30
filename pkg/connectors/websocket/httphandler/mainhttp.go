package httphandler

/*
import (
	"flag"
	"fmt"
	"github.com/yngvark/gr-zombie/pkg/gamelogic"
	"github.com/yngvark/gr-zombie/pkg/pubsub"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type mainHelp struct {
	log *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) *mainHelp {
	return &mainHelp{
		log: logger,
	}
}

func (m *mainHelp) SetupGame(allowedCorsOrigins map[string]bool) {
	broker := pubsub.NewBroadcaster()
	stopGamelogicChannel := make(chan bool)

	var publisher pubsub.Publisher = broker
	httpHandler := NewHTTPHandler(m.log, allowedCorsOrigins, publisher, stopGamelogicChannel)
	http.Handle("/zombie", httpHandler)

	var messageSender pubsub.Publisher = httpHandler
	gameLogic := gamelogic.NewGameLogic(m.log, messageSender, stopGamelogicChannel, nil)

	broker.Subscribe(gameLogic)
}

func (m *mainHelp) HttpListen(port string, lg *zap.SugaredLogger) {
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		out := []byte("OK")
		_, err := writer.Write(out)

		if err != nil {
			m.log.Errorf("error when responding on /health: %s\n", err)
		}
	})

	serverAddr := flag.String("addr", fmt.Sprintf(":%s", port), "http service address")
	lg.Infof("Running on %s\n", *serverAddr)

	log.Fatal(http.ListenAndServe(*serverAddr, nil))
}
*/
