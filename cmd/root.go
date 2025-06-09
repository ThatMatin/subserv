package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "subserv",
	Short: "Subserv is a subscription management service",
	Long:  `Subserv is a subscription management service that allows users to manage their subscriptions easily.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
