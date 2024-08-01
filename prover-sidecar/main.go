package main

import (
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("error: %#+v", err)
	}
}
