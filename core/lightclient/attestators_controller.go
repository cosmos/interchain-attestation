package lightclient

import (
	"context"
)

// TODO: Document the interface and its methods
type AttestatorsController interface {
	SufficientAttestations(ctx context.Context, attestatorIds [][]byte) (bool, error)
}
