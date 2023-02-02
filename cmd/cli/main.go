package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pduzinki/fpl-price-checker/cmd/cli/fpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	fpc.Execute(ctx)
}
