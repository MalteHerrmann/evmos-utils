package cmd

import (
	"os"

	"github.com/MalteHerrmann/evmos-utils/utils"
	evmosutils "github.com/evmos/evmos/v17/utils"
	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands.
	//
	//nolint:gochecknoglobals // required by cobra
	rootCmd = &cobra.Command{
		Use:   "evmos-utils",
		Short: "A collection of utilities to interact with an Evmos node during development.",
		Long: `The evmos-utils collection offers helpers to interact with an Evmos node during development.
It can be used to test upgrades, deposit or vote for specific or the latest proposals, etc.`,
	}

	// appd is the name of the binary to execute commands with.
	appd string
	// chainID is the chain ID of the network.
	chainID string
	// denom of the chain's fee token.
	denom string
	// home is the home directory of the binary.
	home string
	// keyringBackend is the keyring to use.
	keyringBackend string
	// node to post requests and transactions to.
	node string
)

//nolint:gochecknoinits // required by cobra
func init() {
	rootCmd.PersistentFlags().StringVar(
		&appd,
		"bin",
		"evmosd",
		"Name of the binary to be executed",
	)
	rootCmd.PersistentFlags().StringVar(
		&chainID,
		"chain-id",
		evmosutils.TestnetChainID+"-1",
		"Chain ID of the network",
	)
	rootCmd.PersistentFlags().StringVar(
		&denom,
		"denom",
		"aevmos",
		"Fee token denomination of the network",
	)
	rootCmd.PersistentFlags().StringVar(
		&home,
		"home",
		".tmp-evmosd",
		"Home directory of the binary",
	)
	rootCmd.PersistentFlags().StringVar(
		&keyringBackend,
		"keyring-backend",
		"test",
		"Keyring to use",
	)
	rootCmd.PersistentFlags().StringVar(
		&node,
		"node",
		"http://localhost:26657",
		"Node to post queries and transactions to",
	)

	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(depositCmd)
	rootCmd.AddCommand(voteCmd)
}

// collectConfig returns a BinaryConfig filled with the current configuration options
// that depend on the passed flags to the given CLI commands.
func collectConfig() utils.BinaryConfig {
	return utils.BinaryConfig{
		Appd:           appd,
		ChainID:        chainID,
		Denom:          denom,
		Home:           home,
		KeyringBackend: keyringBackend,
		Node:           node,
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
