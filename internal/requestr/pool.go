package requestr

import (
	"context"
	"log"

	"github.com/yansal/requestr/internal/broker"
	"github.com/yansal/requestr/internal/store"
)

type Pool struct {
	store     *store.Store
	receiver  broker.Receiver
	handler   broker.Handler
	publisher broker.Publisher
	id        int
}

func NewPool(ctx context.Context, store *store.Store, receiver broker.Receiver, handler broker.Handler, publisher broker.Publisher) (*Pool, error) {
	id, err := store.AddPool(ctx)
	if err != nil {
		return nil, err
	}
	return &Pool{store: store, receiver: receiver, handler: handler, publisher: publisher, id: id}, nil
}

func (p *Pool) Add(ctx context.Context, queue string) error {
	_, err := p.store.AddJob(ctx, queue, p.id)
	go func() {
		if err := p.receiver.Receive(ctx, queue, func(ctx context.Context, msg broker.Message) error {
			msg.PoolID = p.id
			err := p.handler(ctx, msg)
			if err != nil {
				return err
			}
			if msg.Loop {
				return p.publisher.Publish(ctx, queue, msg.Payload, broker.Loop)
			}
			return nil
		}); err != nil {
			log.Print(err)
		}
	}()

	return err
}

func (p *Pool) Remove(ctx context.Context, queue string) error {
	return nil
}
