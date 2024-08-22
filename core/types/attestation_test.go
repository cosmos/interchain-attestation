package types_test

import (
	"crypto/rand"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/gjermundgaraba/interchain-attestation/core/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	mockChainID  = "testchain-1"
	mockClientID = "testclient-1"
)

func TestGetSignableBytes(t *testing.T) {
	cdc := testutil.MakeTestEncodingConfig().Codec

	for i := 0; i < 10; i++ {
		var packetCommitments [][]byte
		for j := 0; j < i; j++ {
			packetCommitments = append(packetCommitments, getRandomBytes(j))
		}

		attestationData := types.IBCData{
			ChainId:           mockChainID,
			ClientId:          mockClientID,
			Height:            clienttypes.NewHeight(1, 42),
			Timestamp:         time.Now(),
			PacketCommitments: packetCommitments,
		}
		expectedSignableBytes := types.GetSignableBytes(cdc, attestationData)

		var signers []*secp256k1.PrivKey
		for j := 0; j < i; j++ {
			signers = append(signers, secp256k1.GenPrivKey())
		}

		for j := 0; j < 10; j++ {
			for _ , signer := range signers {
				bz := types.GetSignableBytes(cdc, attestationData)
				require.NotNil(t, bz)

				// verify bytes are the same every time
				require.Equal(t, expectedSignableBytes, bz)

				signature, err := signer.Sign(bz)
				require.NoError(t, err)

				pubKey := signer.PubKey()
				verified := pubKey.VerifySignature(bz, signature)
				require.True(t, verified)
			}
		}
	}
}

func getRandomBytes(extraBytes int) []byte {
	numBytes := 16 + extraBytes

	// Create a byte slice to hold the random bytes
	randomBytes := make([]byte, numBytes)

	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return randomBytes
}