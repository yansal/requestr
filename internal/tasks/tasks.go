package tasks

import (
	"context"

	"github.com/yansal/requestr/internal/broker"
	"github.com/yansal/requestr/internal/requestr"
)

type Tasks map[string]broker.Handler

func New(pool *requestr.Pool) Tasks {
	return Tasks{
		// "debug":       DebugHandler,
		// "sleep":       sleepHandler,
		// "error":       errorHandler,
		// "panic":       panicHandler,
		// "httprequest": httpRequestHandler,

		"add":    addHandler(pool),
		"remove": removeHandler(pool),
	}
}

func addHandler(pool *requestr.Pool) func(ctx context.Context, msg broker.Message) error {
	return func(ctx context.Context, msg broker.Message) error {
		return pool.Add(ctx, msg.Payload)
	}
}

func removeHandler(pool *requestr.Pool) func(ctx context.Context, msg broker.Message) error {
	return func(ctx context.Context, msg broker.Message) error {
		return pool.Remove(ctx, msg.Payload)
	}
}
