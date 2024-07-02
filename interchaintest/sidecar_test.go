package pessimisticinterchaintest

import (
	"bytes"
	"github.com/pelletier/go-toml/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"proversidecar/config"
	"proversidecar/server"
	"time"
)

func (s *E2ETestSuite) TestSidecar() {
	s.NotNil(s.ic)

	sidecarConfig := config.Config{
		CosmosChains: []config.CosmosChainConfig{
			{
				ChainID:  rollyChainID,
				RPC:      s.rolly.GetRPCAddress(),
				ClientID: "07-tendermint-0",
			},
		},
	}
	byteWriter := new(bytes.Buffer)
	err := toml.NewEncoder(byteWriter).Encode(sidecarConfig)
	s.NoError(err)
	sidecarConfigBz := byteWriter.Bytes()

	for _, val := range s.hub.Validators {
		s.Len(val.Sidecars, 1)

		sidecar := val.Sidecars[0]
		err = sidecar.WriteFile(s.ctx, sidecarConfigBz, "config.toml")
		s.NoError(err)

		err = sidecar.CreateContainer(s.ctx)
		s.NoError(err)

		err = sidecar.StartContainer(s.ctx)
		s.NoError(err)

		// Wait for the sidecar to be ready
		time.Sleep(2 * time.Second)

		hostPorts, err := sidecar.GetHostPorts(s.ctx, "6969/tcp")
		s.NoError(err)
		s.Len(hostPorts, 1)

		client, err := grpc.NewClient(hostPorts[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
		s.NoError(err)
		defer client.Close()

		proofClient := server.NewProofClient(client)
		proof, err := proofClient.GetProof(s.ctx, &server.ProofRequest{ChainId: rollyChainID})
		s.NoError(err)
		s.NotNil(proof)
		s.NotEmpty(proof.Proof)
	}
}
