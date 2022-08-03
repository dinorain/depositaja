package detector

import (
	"context"
	"time"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"

	"github.com/dinorain/depositaja"
	"github.com/dinorain/depositaja/flagger"
	"github.com/dinorain/depositaja/proto/pb"
)

const (
	rollingPeriod = 120 // in seconds
	maxAmount     = 10000
)

var (
	group goka.Group = "threshold"
)

type CounterCodec struct{}

func (c *CounterCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.Counter))
}

func (c *CounterCodec) Decode(data []byte) (interface{}, error) {
	var m pb.Counter
	return &m, proto.Unmarshal(data, &m)
}

func getValue(ctx goka.Context) *pb.Counter {
	if v := ctx.Value(); v != nil {
		return v.(*pb.Counter)
	}
	return &pb.Counter{}
}

func detectSpammer(ctx goka.Context, c *pb.Counter) bool {
	//log.Printf("Deposit total:\n %v starting from time %v\n", c.Received, c.RollingPeriodStartUnix)
	return c.Received >= maxAmount && c.RollingPeriodStartUnix != 0
}

func Run(ctx context.Context, brokers []string) func() error {
	return func() error {
		g := goka.DefineGroup(group,
			goka.Input(depositaja.DepositStream, new(depositaja.DepositCodec), func(ctx goka.Context, msg interface{}) {
				c := getValue(ctx)

				m := msg.(*pb.Deposit)
				c.Received += m.Amount

				if c.RollingPeriodStartUnix == 0 {
					c.RollingPeriodStartUnix = time.Now().Unix()
				} else {
					if time.Now().Unix()-c.RollingPeriodStartUnix > rollingPeriod {
						c.RollingPeriodStartUnix = 0
						c.Received = 0
					}
				}
				ctx.SetValue(c)

				// check if deposit has been above threshold
				if detectSpammer(ctx, c) {
					ctx.Emit(flagger.Stream, ctx.Key(), &pb.FlagEvent{FlagRemoved: false, RollingPeriodStartUnix: c.RollingPeriodStartUnix})
				} else {
					ctx.Emit(flagger.Stream, ctx.Key(), &pb.FlagEvent{FlagRemoved: true})
				}
			}),
			goka.Output(flagger.Stream, new(flagger.FlagEventCodec)),
			goka.Persist(new(CounterCodec)),
		)
		p, err := goka.NewProcessor(brokers, g)
		if err != nil {
			return err
		}

		return p.Run(ctx)
	}
}
