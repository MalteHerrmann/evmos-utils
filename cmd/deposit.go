package cmd

import (
	"log"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // required by cobra
var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit for a governance proposal",
	Long: `Deposit the minimum needed deposit for a given governance proposal.
If no proposal ID is given by the user, the latest proposal is queried and deposited for.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		bin, err := utils.NewEvmosTestingBinary()
		if err != nil {
			log.Fatalf("error creating binary: %v", err)
		}

		if err = bin.GetAccounts(); err != nil {
			log.Fatalf("error getting accounts: %v", err)
		}

		err = gov.Deposit(bin, args)
		if err != nil {
			log.Fatalf("error depositing: %v", err)
		}
	},
}

//nolint:gochecknoinits // required by cobra
func init() {
	rootCmd.AddCommand(depositCmd)
}
