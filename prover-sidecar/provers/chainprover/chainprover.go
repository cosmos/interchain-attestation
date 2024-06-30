package chainprover

import (
	"context"
)

type ChainProver interface {
	ChainID() string
	CollectProofs(ctx context.Context) error
	GetProof() []byte // TODO: Make more specific to get the correct proof (or all proofs, not sure)
}