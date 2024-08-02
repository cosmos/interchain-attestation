package main

import (
	"github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("error: %#+v", err)
	}
}
