package chainattestor

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/types"
)


type ChainAttestor interface {
	ChainID() string
	CollectClaims(ctx context.Context) error
	GetLatestSignedClaim() *types.SignedPacketCommitmentsClaim
}