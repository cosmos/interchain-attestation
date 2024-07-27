package lightclient

import (
	"context"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto/types"
)

// TODO: Document the interface and its methods
type AttestatorsHandler interface{
	GetPublicKey(ctx context.Context, attestatorId []byte) (sdkcrypto.PubKey, error)
	SufficientAttestations(ctx context.Context, attestatorIds [][]byte) (bool, error)
}
