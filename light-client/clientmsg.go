package lightclient

import (
	"crypto/sha256"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ClientMessage = (*PessimisticClaims)(nil)

func NewPessimisticClaims(claims []SignedPacketCommitmentsClaim) *PessimisticClaims {
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

func GetSignableBytes(packetCommitements [][]byte) []byte {
	var packetBytes []byte

	for _, packetCommitement := range packetCommitements {
		packetBytes = append(packetBytes, packetCommitement...)
	}

	hash := sha256.Sum256(packetBytes)
	return hash[:]
}