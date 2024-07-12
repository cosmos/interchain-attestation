package pessimisticinterchaintest

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"cosmossdk.io/math"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/ethereum"
)

func (s *E2ETestSuite) TestDeployContract() {
	ctx := context.Background()
	s.NotNil(s.ic)

	startingEthBalance := math.NewInt(2 * ethereum.ETHER)
	users := interchaintest.GetAndFundTestUsers(s.T(), s.ctx, s.T().Name(), startingEthBalance, s.eth)
	ethUser := users[0]

	s.Require().True(s.Run("Deploy contracts", func() {
		s.Require().NoError(os.Chdir(".."))

		stdout, _, err := s.eth.ForgeScript(ctx, ethUser.KeyName(), ethereum.ForgeScriptOpts{
			ContractRootDir:  "contracts",
			SolidityContract: "script/Counter.s.sol",
			RawOptions:       []string{"--json"},
		})
		s.Require().NoError(err)

		contractAddress := s.GetEthAddressFromStdout(string(stdout))
		s.Require().NotEmpty(contractAddress)
		s.Require().True(ethcommon.IsHexAddress(contractAddress))

		_, err = ethclient.Dial(s.eth.GetHostRPCAddress())
		s.Require().NoError(err)

		// s.contract, err = sp1ics07tendermint.NewContract(ethcommon.HexToAddress(contractAddress), client)
		// s.Require().NoError(err)
	}))
}

func (s *E2ETestSuite) GetEthAddressFromStdout(stdout string) string {
	// Define the regular expression pattern
	re := regexp.MustCompile(`"value":"(0x[0-9a-fA-F]+)"`)

	// Find the first match
	matches := re.FindStringSubmatch(stdout)
	if len(matches) <= 1 {
		s.FailNow(fmt.Sprintf("no match found in stdout: %s", stdout))
	}
	// Extract the value
	return matches[1]
}
