package chainprover

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/types"
)

type ChainProver interface {
	ChainID() string
	CollectProofs(ctx context.Context) error
	GetProof() *types.SignedPacketCommitmentsClaim
}