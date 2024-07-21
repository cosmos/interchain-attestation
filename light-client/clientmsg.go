package lightclient

import "github.com/cosmos/ibc-go/v8/modules/core/exported"

var _ exported.ClientMessage = (*PessimisticClaims)(nil)

func NewPessimisticClaims(claims []PacketCommitmentsClaim) *PessimisticClaims {
	return &PessimisticClaims{
		Claims: claims,
	}
}

func (m *PessimisticClaims) ClientType() string {
	return ModuleName
}

func (m *PessimisticClaims) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}
