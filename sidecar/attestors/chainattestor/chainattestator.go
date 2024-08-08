package chainattestor

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
)

// TODO: Document
type ChainAttestor interface {
	ChainID() string
	CollectAttestation(ctx context.Context) (types.Attestation, error)
}