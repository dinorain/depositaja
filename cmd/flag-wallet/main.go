package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lovoo/goka"

	"github.com/dinorain/depositaja/flagger"
	"github.com/dinorain/depositaja/pb"
)

var (
	wallet = flag.String("wallet", "", "wallet_id to be flagged as above-threshold deposit")
	remove = flag.Bool("remove", false, "remove flag of wallet_id instead")
	broker = flag.String("broker", "localhost:9092", "boostrap Kafka broker")
)

func main() {
	flag.Parse()
	if *wallet == "" {
		fmt.Println("cannot remove flag of wallet_id ''")
		os.Exit(1)
	}
	emitter, err := goka.NewEmitter([]string{*broker}, flagger.Stream, new(flagger.FlagEventCodec))
	if err != nil {
		panic(err)
	}
	defer emitter.Finish()

	err = emitter.EmitSync(*wallet, &pb.FlagEvent{FlagRemoved: *remove})
	if err != nil {
		panic(err)
	}
}
