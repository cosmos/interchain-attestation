package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/cosmos/interchain-attestation/simapp"
	"github.com/cosmos/interchain-attestation/simapp/params"
	"github.com/cosmos/interchain-attestation/simapp/simappd/cmd"
)

func main() {
	params.InitSDKConfig()
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", simapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
