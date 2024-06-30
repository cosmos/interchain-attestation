package provers

import (
	"context"
)

type ChainProver interface {
	ChainID() string
	CollectProofs(ctx context.Context) error
}