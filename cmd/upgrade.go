package cmd

import (
	"log"
	"regexp"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // required by cobra
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Prepare an upgrade of a node",
	Long: `Prepare an upgrade of a node by submitting a governance proposal, 
voting for it using all keys of in the keyring and having it pass.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bin, err := utils.NewEvmosTestingBinary()
		if err != nil {
			log.Fatalf("error creating binary: %v", err)
		}

		if err = bin.GetAccounts(); err != nil {
			log.Fatalf("error getting accounts: %v", err)
		}

		targetVersion := args[0]
		if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d(-rc\d+)?`, targetVersion); !matched {
			log.Fatalf("invalid target version: %s; please use the format vX.Y.Z(-rc*).\n", targetVersion)
		}

		if err := upgradeLocalNode(bin, targetVersion); err != nil {
			log.Fatalf("error upgrading local node: %v", err)
		}
	},
}

//nolint:gochecknoinits // required by cobra
func init() {
	rootCmd.AddCommand(upgradeCmd)
}

// upgradeLocalNode prepares upgrading the local node to the target version
// by submitting the upgrade proposal and voting on it using all testing accounts.
func upgradeLocalNode(bin *utils.Binary, targetVersion string) error {
	currentHeight, err := utils.GetCurrentHeight(bin)
	if err != nil {
		return errors.Wrap(err, "error getting current height")
	}

	upgradeHeight := currentHeight + utils.DeltaHeight

	log.Println("Submitting upgrade proposal...")

	proposalID, err := gov.SubmitUpgradeProposal(bin, targetVersion, upgradeHeight)
	if err != nil {
		return errors.Wrap(err, "error executing upgrade proposal")
	}

	log.Printf("Scheduled upgrade to %s at height %d.\n", targetVersion, upgradeHeight)

	if err = gov.SubmitAllVotesForProposal(bin, proposalID); err != nil {
		return errors.Wrapf(err, "error submitting votes for proposal %d", proposalID)
	}

	return nil
}
