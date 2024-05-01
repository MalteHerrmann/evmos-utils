package cmd

import (
	"fmt"
	"regexp"

	"github.com/MalteHerrmann/evmos-utils/gov"
	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // required by cobra
var upgradeCmd = &cobra.Command{
	Use:   "upgrade TARGET_VERSION",
	Short: "Prepare an upgrade of a node",
	Long: `Prepare an upgrade of a node by submitting a governance proposal, 
voting for it using all keys of in the keyring and having it pass.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		bin, err := utils.NewBinary(collectConfig())
		if err != nil {
			return errors.Wrap(err, "error creating binary")
		}

		targetVersion := args[0]
		if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d(-rc\d+)?`, targetVersion); !matched {
			return fmt.Errorf("invalid target version: %s; please use the format vX.Y.Z(-rc*)", targetVersion)
		}

		if err = upgradeLocalNode(bin, targetVersion); err != nil {
			return errors.Wrap(err, "error upgrading local node")
		}

		bin.Logger.Info().Msgf("successfully prepared upgrade to %s", targetVersion)

		return nil
	},
}

// upgradeLocalNode prepares upgrading the local node to the target version
// by submitting the upgrade proposal and voting on it using all testing accounts.
func upgradeLocalNode(bin *utils.Binary, targetVersion string) error {
	currentHeight, err := utils.GetCurrentHeight(bin)
	if err != nil {
		return errors.Wrap(err, "error getting current height")
	}

	upgradeHeight := currentHeight + utils.DeltaHeight

	bin.Logger.Info().Msg("submitting upgrade proposal...")

	proposalID, err := gov.SubmitUpgradeProposal(bin, targetVersion, upgradeHeight)
	if err != nil {
		return errors.Wrap(err, "error executing upgrade proposal")
	}

	bin.Logger.Info().Msgf("scheduled upgrade to %s at height %d.\n", targetVersion, upgradeHeight)

	if _, err = gov.Deposit(bin, []string{}); err != nil {
		return errors.Wrapf(err, "error depositing for proposal %d", proposalID)
	}

	if err = gov.SubmitAllVotesForProposal(bin, proposalID); err != nil {
		return errors.Wrapf(err, "error submitting votes for proposal %d", proposalID)
	}

	return nil
}
