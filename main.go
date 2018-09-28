package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yansal/requestr/internal/app"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app := app.New()
	g, ctx := errgroup.WithContext(context.Background())

	// Signal handler
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return nil
		case s := <-c:
			return fmt.Errorf("%v", s)
		}
	})

	if len(os.Args) == 1 || os.Args[1] != "worker" {
		g.Go(app.StartServer(ctx))
	}
	g.Go(app.StartWorker(ctx))

	if err := g.Wait(); err != nil {
		log.Fatalf("%+v", err)
	}
}
