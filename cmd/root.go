package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
//
//nolint:gochecknoglobals // required by cobra
var rootCmd = &cobra.Command{
	Use:   "evmos-utils",
	Short: "A collection of utilities to interact with an Evmos node during development.",
	Long: `The evmos-utils collection offers helpers to interact with an Evmos node during development.
It can be used to test upgrades, deposit or vote for specific or the latest proposals, etc.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
