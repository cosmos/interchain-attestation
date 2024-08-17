package attestator

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
)

// TODO: Document
type Attestator interface {
	ChainID() string
	CollectAttestation(ctx context.Context) (types.Attestation, error)
}