package interchaintest

import (
	"context"
	"testing"

	testifysuite "github.com/stretchr/testify/suite"

	"github.com/cosmos/interchain-attestation/interchaintest/suite"
)

type InterchainAttestationTestSuite struct {
	suite.E2ETestSuite
}

func TestInterchainAttestationTestSuite(t *testing.T) {
	testifysuite.Run(t, new(InterchainAttestationTestSuite))
}

func (s *InterchainAttestationTestSuite) TestDeploy() {
	ctx := context.Background()

	cosmosHeight, err := s.Simapp.Height(ctx)
	s.Require().NoError(err)
	s.Require().NotZero(cosmosHeight)

	evmHeight, err := s.EVM.Height(ctx)
	s.Require().NoError(err)
	s.Require().NotZero(evmHeight)
}
