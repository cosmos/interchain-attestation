package lightclient

import (
	"crypto/sha256"
	"github.com/cosmos/cosmos-sdk/codec"
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

func GetSignableBytes(cdc codec.BinaryCodec, claim PacketCommitmentsClaim) []byte {
	packetBytes := cdc.MustMarshal(&claim)
	hash := sha256.Sum256(packetBytes)
	return hash[:]
}