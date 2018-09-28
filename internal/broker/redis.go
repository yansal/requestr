package broker

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type RedisPublisher struct {
	redis *redis.Client
}

func NewRedisPublisher(redis *redis.Client) *RedisPublisher { return &RedisPublisher{redis: redis} }

func (b *RedisPublisher) Publish(ctx context.Context, queue string, payload string, opts ...PublishOptions) error {
	msg := Message{Payload: payload}
	for i := range opts {
		msg = opts[i](msg)
	}
	_, err := b.redis.LPush(queue, msg).Result()
	return errors.WithStack(err)
}

type RedisReceiver struct {
	redis             *redis.Client
	brpoplpushTimeout time.Duration
}

func NewRedisReceiver(redis *redis.Client) *RedisReceiver {
	return &RedisReceiver{redis: redis, brpoplpushTimeout: 0}
}

func (r *RedisReceiver) Receive(ctx context.Context, queue string, handler Handler) error {
	// TODO: spawn a garbage collecting goroutine to look at the procesing queue?
	processing := queue + "-processing"

	type msg struct {
		message Message
		err     error
	}
	brpoplpush := make(chan msg)

	for {
		go func() {
			var message Message
			err := r.redis.BRPopLPush(queue, processing, r.brpoplpushTimeout).Scan(&message)
			brpoplpush <- msg{message: message, err: err}
		}()

		var message Message
		select {
		case <-ctx.Done():
			return nil
		case msg := <-brpoplpush:
			if err := msg.err; err == redis.Nil {
				continue
			} else if err != nil {
				return errors.WithStack(err)
			}
			message = msg.message
		}

		if err := handler(ctx, message); err != nil {
			// TOD: check if context.Canceled
			// TODO: delete from processing? retry?
			return err
		}

		if _, err := r.redis.LRem(processing, 1, message).Result(); err != nil {
			return errors.WithStack(err)
		}
	}
}
