package relayer

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"go.uber.org/zap"
)

type Relayer struct {
	logger *zap.Logger

	cdc    codec.Codec
	homedir string
}

func NewRelayer(logger *zap.Logger, cdc codec.Codec, homedir string) *Relayer {
	return &Relayer{
		logger: logger,
		cdc:    cdc,
		homedir: homedir,
	}
}
