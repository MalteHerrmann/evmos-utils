package cmd

import (
	"github.com/MalteHerrmann/evmos-utils/gov"
	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // required by cobra
var voteCmd = &cobra.Command{
	Use:   "vote [PROPOSAL_ID]",
	Short: "Vote for a governance proposal",
	Long: `Vote for a governance proposal with all keys in the keyring.
If no proposal ID is passed, the latest proposal on chain is queried and used.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		bin, err := utils.NewEvmosTestingBinary()
		if err != nil {
			bin.Logger.Error().Msgf("error creating binary: %v", err)
			return
		}

		if err = bin.GetAccounts(); err != nil {
			bin.Logger.Error().Msgf("error getting accounts: %v", err)
			return
		}

		proposalID, err := gov.SubmitAllVotes(bin, args)
		if err != nil {
			bin.Logger.Error().Msgf("error submitting votes: %v", err)
			return
		}

		bin.Logger.Info().Msgf("successfully submitted votes for proposal %d", proposalID)
	},
}

//nolint:gochecknoinits // required by cobra
func init() {
	rootCmd.AddCommand(voteCmd)
}
