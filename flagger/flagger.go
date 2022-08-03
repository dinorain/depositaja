package flagger

import (
	"context"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"

	"github.com/dinorain/depositaja/proto/pb"
)

var (
	group  goka.Group  = "flagger"
	Table  goka.Table  = goka.GroupTable(group)
	Stream goka.Stream = "flag_wallet"
)

type FlagEventCodec struct{}

func (c *FlagEventCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.FlagEvent))
}

func (c *FlagEventCodec) Decode(data []byte) (interface{}, error) {
	var m pb.FlagEvent
	return &m, proto.Unmarshal(data, &m)
}

type FlagValueCodec struct{}

func (c *FlagValueCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.FlagValue))
}

func (c *FlagValueCodec) Decode(data []byte) (interface{}, error) {
	var m pb.FlagValue
	return &m, proto.Unmarshal(data, &m)
}

func flag(ctx goka.Context, msg interface{}) {
	var s *pb.FlagValue
	if v := ctx.Value(); v == nil {
		s = new(pb.FlagValue)
	} else {
		s = v.(*pb.FlagValue)
	}

	flagEvent := msg.(*pb.FlagEvent)
	if flagEvent.FlagRemoved {
		s.Flagged = false
		s.RollingPeriodStartUnix = 0
	} else {
		s.Flagged = true
		s.RollingPeriodStartUnix = flagEvent.RollingPeriodStartUnix
	}
	ctx.SetValue(s)
}

func Run(ctx context.Context, brokers []string) func() error {
	return func() error {
		g := goka.DefineGroup(group,
			goka.Input(Stream, new(FlagEventCodec), flag),
			goka.Persist(new(FlagValueCodec)),
		)
		p, err := goka.NewProcessor(brokers, g)
		if err != nil {
			return err
		}
		return p.Run(ctx)
	}
}
