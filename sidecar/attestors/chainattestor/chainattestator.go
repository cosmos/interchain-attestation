package chainattestor

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
)


type ChainAttestor interface {
	ChainID() string
	CollectAttestations(ctx context.Context) error
	GetLatestAttestation() *types.Attestation
}