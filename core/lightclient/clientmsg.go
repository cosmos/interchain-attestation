package lightclient

import (
	"github.com/cosmos/ibc-go/v9/modules/core/exported"

	"github.com/cosmos/interchain-attestation/core/types"
)

var _ exported.ClientMessage = (*AttestationClaim)(nil)

func NewAttestationClaim(attestation []types.Attestation) *AttestationClaim {
	return &AttestationClaim{
		Attestations: attestation,
	}
}

func (m *AttestationClaim) ClientType() string {
	return ModuleName
}

func (m *AttestationClaim) ValidateBasic() error {
	// TODO: implement me
	panic("implement me")
}
