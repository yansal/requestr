package worker

import (
	"context"

	"github.com/yansal/requestr/internal/broker"
	"github.com/yansal/requestr/internal/tasks"
	"golang.org/x/sync/errgroup"
)

func Start(ctx context.Context, receiver broker.Receiver, tasks tasks.Tasks) func() error {
	return func() error {
		g, ctx := errgroup.WithContext(ctx)
		for queue, handler := range tasks {
			queue := queue
			handler := handler
			g.Go(func() error {
				return receiver.Receive(ctx, queue, handler)
			})
		}
		return g.Wait()
	}
}
