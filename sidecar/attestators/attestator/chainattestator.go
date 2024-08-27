package attestator

import (
	"context"

	"github.com/cosmos/interchain-attestation/core/types"
)

// TODO: Document
type Attestator interface {
	ChainID() string
	CollectAttestation(ctx context.Context) (types.Attestation, error)
}
