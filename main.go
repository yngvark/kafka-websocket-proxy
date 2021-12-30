package main

import (
	"fmt"
	proxyPkg "github.com/yngvark/kafka-websocket-proxy/pkg/proxy"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(fmt.Errorf("error running game: %w", err))
	}

	fmt.Println("Main ended.")
}

func run() error {
	ctx, cancelFn := context.WithCancel(context.Background())
	osInterruptChan := make(chan os.Signal, 1)

	signal.Notify(osInterruptChan, os.Interrupt)

	// Don't listen for interrupts after program quits
	defer func() {
		signal.Stop(osInterruptChan)
		cancelFn()
	}()

	// Listen in the background (i.e. goroutine) if the OS interrupts our program.
	go cancelProgramIfOsInterrupts(ctx, osInterruptChan, cancelFn)

	// Setup app
	appOpts, err := newAppOpts(ctx, cancelFn, os.Getenv)
	if err != nil {
		cancelFn()
		return fmt.Errorf("creating dependencies: %w", err)
	}

	return runHTTPServer(appOpts)
}

func runHTTPServer(appOpts *AppOpts) error {
	port, ok := os.LookupEnv("GAME_PORT")
	if !ok {
		port = "8080"
	}

	serverAddr := ":" + port
	appOpts.logger.Infof("Running on %s\n", serverAddr)

	go func() {
		appOpts.logger.Info("Now attempting to listen on port " + port)

		err := http.ListenAndServe(serverAddr, nil)
		appOpts.logger.Errorf("HTTP listen and serve: %s", err.Error())
		appOpts.cancelFn()
	}()

	http.HandleFunc("/health", health)

	p := proxyPkg.New(appOpts.context, appOpts.cancelFn, appOpts.logger)
	return p.Run()
}

func cancelProgramIfOsInterrupts(ctx context.Context, osInterruptChan chan os.Signal, cancelFn context.CancelFunc) {
	func() {
		select {
		case <-osInterruptChan:
			cancelFn()
		case <-ctx.Done():
			// Stop listening
			return
		}
	}()
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}
