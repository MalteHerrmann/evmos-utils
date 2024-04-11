package cmd

import (
	"github.com/MalteHerrmann/evmos-utils/gov"
	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // required by cobra
var depositCmd = &cobra.Command{
	Use:   "deposit [PROPOSAL_ID]",
	Short: "Deposit for a governance proposal",
	Long: `Deposit the minimum needed deposit for a given governance proposal.
If no proposal ID is given by the user, the latest proposal is queried and deposited for.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		bin, err := utils.NewEvmosTestingBinary()
		if err != nil {
			bin.Logger.Error().Msgf("error creating binary: %v", err)

			return
		}

		proposalID, err := gov.Deposit(bin, args)
		if err != nil {
			bin.Logger.Error().Msgf("error depositing: %v", err)

			return
		}

		bin.Logger.Info().Msgf("successfully deposited for proposal %d", proposalID)
	},
}
