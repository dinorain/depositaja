package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/dinorain/depositaja/collector"
)

var (
	brokers      = []string{"localhost:9092"}
	runCollector = flag.Bool("collector", false, "run collector processor")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	if *runCollector {
		log.Println("starting collector")
		grp.Go(collector.Run(ctx, brokers))
	}

	// Wait for SIGINT/SIGTERM
	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-waiter:
	case <-ctx.Done():
	}
	cancel()
	if err := grp.Wait(); err != nil {
		log.Println(err)
	}
	log.Println("done")
}
