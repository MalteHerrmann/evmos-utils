package cmd

import (
	"log"

	"github.com/MalteHerrmann/evmos-utils/gov"
	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // required by cobra
var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Vote for a governance proposal",
	Long: `Vote for a governance proposal with all keys in the keyring.
If no proposal ID is passed, the latest proposal on chain is queried and used.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		bin, err := utils.NewEvmosTestingBinary()
		if err != nil {
			log.Fatalf("error creating binary: %v", err)
		}

		if err = bin.GetAccounts(); err != nil {
			log.Fatalf("error getting accounts: %v", err)
		}

		if err = gov.SubmitAllVotes(bin, args); err != nil {
			log.Fatal(err)
		}
	},
}

//nolint:gochecknoinits // required by cobra
func init() {
	rootCmd.AddCommand(voteCmd)
}
