package broker

import (
	"context"
	"encoding/json"
)

type Publisher interface {
	Publish(ctx context.Context, queue string, payload string, opts ...PublishOptions) error
}
type PublishOptions func(Message) Message

func Loop(msg Message) Message {
	msg.Loop = true
	return msg
}

type Receiver interface {
	Receive(ctx context.Context, queue string, handler Handler) error
}

type Handler func(context.Context, Message) error
type Message struct {
	PoolID  int
	Loop    bool
	Payload string
}

func (message Message) MarshalBinary() ([]byte, error)     { return json.Marshal(message) }
func (message *Message) UnmarshalBinary(data []byte) error { return json.Unmarshal(data, message) }
