package tasks

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/yansal/requestr/internal/broker"
)

func DebugHandler(ctx context.Context, msg broker.Message) error {
	log.Printf("%+v", msg)
	return nil
}

func sleepHandler(ctx context.Context, msg broker.Message) error {
	duration, err := time.ParseDuration(msg.Payload)
	if err != nil {
		return errors.WithStack(err)
	}
	select {
	case <-time.After(duration):
	case <-ctx.Done():
	}
	return nil
}

func errorHandler(ctx context.Context, msg broker.Message) error {
	return errors.New(msg.Payload)
}

func panicHandler(ctx context.Context, msg broker.Message) error {
	panic(msg.Payload)
}

func republishHandler(publisher broker.Publisher, queue string) func(context.Context, broker.Message) error {
	return func(ctx context.Context, msg broker.Message) error {
		return publisher.Publish(ctx, queue, msg.Payload)
	}
}
