package relayer

import (
	"go.uber.org/zap"

	"github.com/cosmos/cosmos-sdk/codec"
)

type Relayer struct {
	logger *zap.Logger

	cdc     codec.Codec
	homedir string
}

func NewRelayer(logger *zap.Logger, cdc codec.Codec, homedir string) *Relayer {
	return &Relayer{
		logger:  logger,
		cdc:     cdc,
		homedir: homedir,
	}
}
