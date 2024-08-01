package pessimisticinterchaintest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/config"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/server"
	"github.com/pelletier/go-toml/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"path"
	"time"
)

func (s *E2ETestSuite) TestSidecar() {
	s.NotNil(s.ic)

	chainConfigs := []config.CosmosChainConfig{
		{
			ChainID:  rollupsimappChainID,
			RPC:      s.rollupsimapp.GetRPCAddress(),
			ClientID: "07-tendermint-0",
		},
	}

	for i, val := range s.simapp.Validators {
		s.Len(val.Sidecars, 1)
		sidecar := val.Sidecars[0]

		privKey := secp256k1.GenPrivKey()
		bz := privKey.Bytes()
		base64EncodedKey := make([]byte, base64.StdEncoding.EncodedLen(len(bz)))
		base64.StdEncoding.Encode(base64EncodedKey, bz)

		privKeyFilename := "priv.key"
		err := sidecar.WriteFile(s.ctx, base64EncodedKey, privKeyFilename)
		s.NoError(err)

		privKeyPath := path.Join("/home/sidecar", privKeyFilename)

		sidecarConfig := config.Config{
			CosmosChains:          chainConfigs,
			AttestatorID:          fmt.Sprintf("attestator-%d", i),
			SigningPrivateKeyPath: privKeyPath,
		}

		byteWriter := new(bytes.Buffer)
		err = toml.NewEncoder(byteWriter).Encode(sidecarConfig)
		s.NoError(err)
		err = sidecar.WriteFile(s.ctx, byteWriter.Bytes(), "config.toml")
		s.NoError(err)

		err = sidecar.CreateContainer(s.ctx)
		s.NoError(err)

		err = sidecar.StartContainer(s.ctx)
		s.NoError(err)

		// Wait for the sidecar to be ready
		time.Sleep(2 * time.Second)
	}

	// TODO: Put some packets into the chain

	// Give the sidecars some time to do their thing
	time.Sleep(2 * time.Second)

	for i, val := range s.simapp.Validators {
		s.Len(val.Sidecars, 1)
		sidecar := val.Sidecars[0]

		hostPorts, err := sidecar.GetHostPorts(s.ctx, "6969/tcp")
		s.NoError(err)
		s.Len(hostPorts, 1)

		client, err := grpc.NewClient(hostPorts[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
		s.NoError(err)
		defer client.Close()

		proofClient := server.NewProofClient(client)
		proof, err := proofClient.GetProof(s.ctx, &server.ProofRequest{ChainId: rollupsimappChainID})
		s.NoError(err)
		s.NotNil(proof)
		s.Equal([]byte(fmt.Sprintf("attestator-%d", i)), proof.Proof.AttestatorId)
	}
}
