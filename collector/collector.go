package collector

import (
	"context"

	"github.com/lovoo/goka"

	"github.com/dinorain/depositaja"
	"github.com/dinorain/depositaja/pb"
)

var (
	group goka.Group = "balance"
	Table goka.Table = goka.GroupTable(group)
)

func collect(ctx goka.Context, msg interface{}) {
	ml := &pb.DepositHistory{}
	if v := ctx.Value(); v != nil {
		ml = v.(*pb.DepositHistory)
	}

	m := msg.(*pb.Deposit)

	ml.WalletId = m.WalletId
	ml.Deposits = append(ml.Deposits, m)

	ctx.SetValue(ml)
}

func Run(ctx context.Context, brokers []string) func() error {
	return func() error {
		g := goka.DefineGroup(group,
			goka.Input(depositaja.DepositStream, new(depositaja.DepositCodec), collect),
			goka.Persist(new(depositaja.DepositListCodec)),
		)
		p, err := goka.NewProcessor(brokers, g)
		if err != nil {
			return err
		}
		return p.Run(ctx)
	}
}
