package main

import (
	"flag"

	"github.com/dinorain/depositaja"
	"github.com/dinorain/depositaja/service"
)

var (
	withdraw = flag.Bool("withdraw", false, "emit to WithdrawStream")
	broker   = flag.String("broker", "localhost:9092", "boostrap Kafka broker")
)

func main() {
	flag.Parse()
	if *withdraw {
		service.Run([]string{*broker}, depositaja.WithdrawStream)
	} else {
		service.Run([]string{*broker}, depositaja.DepositStream)
	}
}
