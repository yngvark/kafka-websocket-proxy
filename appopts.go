package main

import (
	"context"
	"fmt"
	"github.com/yngvark/kafka-websocket-proxy/pkg/lib/log2"
	"go.uber.org/zap"
)

// AppOpts contains various dependencies
type AppOpts struct {
	context  context.Context
	logger   *zap.SugaredLogger
	cancelFn context.CancelFunc
}

type getEnv func(key string) string

func newAppOpts(ctx context.Context, cancelFn context.CancelFunc, getEnv getEnv) (*AppOpts, error) {
	logger, err := log2.New()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	return &AppOpts{
		context: ctx,
		logger:  logger,
	}, nil
}
