package main

import (
	"fmt"
	"github.com/gjermundgaraba/interchain-attestation/rollupsimapp/params"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/gjermundgaraba/interchain-attestation/rollupsimapp"
	"github.com/gjermundgaraba/interchain-attestation/rollupsimapp/rollupsimappd/cmd"
)

func main() {
	params.InitSDKConfig()
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", rollupsimapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
