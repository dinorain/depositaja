package main

import (
	"flag"

	"github.com/dinorain/depositaja"
	"github.com/dinorain/depositaja/service"
)

var (
	broker = flag.String("broker", "localhost:9092", "boostrap Kafka broker")
)

func main() {
	flag.Parse()
	service.Run([]string{*broker}, depositaja.DepositStream)
}
