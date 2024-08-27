package main

import (
	"fmt"
	"github.com/cosmos/interchain-attestation/simapp/params"
	"os"

	"github.com/cosmos/interchain-attestation/simapp"
	"github.com/cosmos/interchain-attestation/simapp/simappd/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	params.InitSDKConfig()
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", simapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
