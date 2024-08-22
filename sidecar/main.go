package main

import (
	"github.com/gjermundgaraba/interchain-attestation/sidecar/cmd"
	"os"
)

func main() {
	rootCmd := cmd.RootCmd()

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("error: %#+v", err)
		os.Exit(1)
	}
}
