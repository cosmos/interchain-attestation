package attestator

import (
	"context"

	"github.com/cosmos/interchain-attestation/core/types"
)

// TODO: Document
type Attestator interface {
	ChainID() string
	CollectIBCData(ctx context.Context) (types.IBCData, error)
}
