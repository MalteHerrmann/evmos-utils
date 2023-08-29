package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/pkg/errors"
)

// The amount of blocks in the future that the upgrade will be scheduled.
const deltaHeight = 15

func main() {
	if len(os.Args) < 2 {
		log.Printf(
			"Possible usages:\n" +
				"  upgrade-local-node-go <target_version>\n" +
				"  upgrade-local-node-go vote [proposal-id]\n",
		)
		os.Exit(1)
	}

	bin, err := utils.NewEvmosTestingBinary()
	if err != nil {
		log.Fatalf("Error creating binary: %v", err)
	}

	//nolint:nestif // nesting complexity is fine here, will be reworked with Cobra commands anyway
	if os.Args[1] == "vote" {
		proposalID, err := getProposalIDForVoting(bin, os.Args)
		if err != nil {
			log.Fatalf("Error getting proposal ID: %v", err)
		}

		err = gov.SubmitAllVotesForProposal(bin, proposalID)
		if err != nil {
			log.Fatalf("Error submitting votes for proposal %d: %v", proposalID, err)
		}
	} else {
		targetVersion := os.Args[1]
		if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d(-rc\d+)?`, targetVersion); !matched {
			log.Fatalf("Invalid target version: %s. Please use the format vX.Y.Z(-rc*).\n", targetVersion)
		}

		err := upgradeLocalNode(bin, targetVersion)
		if err != nil {
			log.Fatalf("Error upgrading local node: %v", err)
		}
	}
}

// getProposalIDForVoting gets the proposal ID from the command line arguments.
func getProposalIDForVoting(bin *utils.Binary, args []string) (int, error) {
	var (
		err        error
		proposalID int
	)

	switch len(args) {
	case 2:
		proposalID, err = gov.QueryLatestProposalID(bin)
		if err != nil {
			return 0, errors.Wrap(err, "Error querying latest proposal ID")
		}
	case 3:
		proposalID, err = strconv.Atoi(args[2])
		if err != nil {
			return 0, errors.Wrapf(err, "Error converting proposal ID %s to integer", args[2])
		}
	default:
		return 0, errors.New("Invalid number of arguments")
	}

	return proposalID, nil
}

// upgradeLocalNode prepares upgrading the local node to the target version
// by submitting the upgrade proposal and voting on it using all testing accounts.
func upgradeLocalNode(bin *utils.Binary, targetVersion string) error {
	currentHeight, err := utils.GetCurrentHeight(bin)
	if err != nil {
		return errors.Wrap(err, "Error getting current height")
	}

	upgradeHeight := currentHeight + deltaHeight

	log.Println("Submitting upgrade proposal...")

	proposalID, err := gov.SubmitUpgradeProposal(bin, targetVersion, upgradeHeight)
	if err != nil {
		return errors.Wrap(err, "Error executing upgrade proposal")
	}

	log.Printf("Scheduled upgrade to %s at height %d.\n", targetVersion, upgradeHeight)

	err = gov.SubmitAllVotesForProposal(bin, proposalID)
	if err != nil {
		return errors.Wrapf(err, "Error submitting votes for proposal %d", proposalID)
	}

	return nil
}
