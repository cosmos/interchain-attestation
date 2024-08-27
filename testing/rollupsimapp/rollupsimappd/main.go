package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/cosmos/interchain-attestation/rollupsimapp"
	"github.com/cosmos/interchain-attestation/rollupsimapp/params"
	"github.com/cosmos/interchain-attestation/rollupsimapp/rollupsimappd/cmd"
)

func main() {
	params.InitSDKConfig()
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", rollupsimapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
