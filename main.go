package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cirocosta/perfer/server"
	"github.com/pkg/errors"
)

var (
	address         = flag.String("address", ":25000", "address to listen for requests")
	assetsDirectory = flag.String("assets", "/tmp", "directory to place assets")
)

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	s, err := server.NewServer(*address, *assetsDirectory)
	if err != nil {
		log.Panic(err)
	}

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- s.Listen()
	}()

	select {
	case err = <-serverDone:
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		if errors.Cause(err) == context.Canceled {
			return
		}

		log.Panic(err)
	}
}
