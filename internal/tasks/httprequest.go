package tasks

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yansal/requestr/internal/broker"
)

func httpRequestHandler(ctx context.Context, msg broker.Message) error {
	resp, err := http.Get(msg.Payload)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()
	log.Print(resp.Status)
	return nil
}
