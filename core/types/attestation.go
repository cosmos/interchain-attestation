package types

import (
	"crypto/sha256"

	"github.com/cosmos/cosmos-sdk/codec"
)

func GetSignableBytes(cdc codec.BinaryCodec, dataToAttestTo IBCData) []byte {
	packetBytes := cdc.MustMarshal(&dataToAttestTo)
	hash := sha256.Sum256(packetBytes)
	return hash[:]
}
