package proxy

import (
	"context"
	"go.uber.org/zap"
)

type Proxy struct {
	context  context.Context
	cancelFn context.CancelFunc
	logger   *zap.SugaredLogger
}
