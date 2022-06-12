package depositaja

import (
	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"

	"github.com/dinorain/depositaja/pb"
)

var (
	DepositStream  goka.Stream = "deposit"
	WithdrawStream goka.Stream = "withdraw"
)

type DepositCodec struct{}

func (c *DepositCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.Deposit))
}

func (c *DepositCodec) Decode(data []byte) (interface{}, error) {
	var m pb.Deposit
	return &m, proto.Unmarshal(data, &m)
}

type DepositListCodec struct{}

func (c *DepositListCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*pb.DepositHistory))
}

func (c *DepositListCodec) Decode(data []byte) (interface{}, error) {
	var m pb.DepositHistory
	return &m, proto.Unmarshal(data, &m)
}
