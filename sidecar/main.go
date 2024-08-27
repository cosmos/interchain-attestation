package main

import (
	"os"

	"github.com/cosmos/interchain-attestation/sidecar/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("error: %#+v", err)
		os.Exit(1)
	}
}
