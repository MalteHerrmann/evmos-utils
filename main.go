package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
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

	if os.Args[1] == "vote" {
		var (
			err        error
			proposalID int
		)

		switch len(os.Args) {
		case 2:
			proposalID, err = gov.QueryLatestProposalID(bin)
			if err != nil {
				log.Fatalf("Error querying latest proposal ID: %v", err)
			}
		case 3:
			proposalID, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Printf("Invalid proposal ID: %s. Please provide an integer.\n", os.Args[2])
				os.Exit(2)
			}
		default:
			log.Println("Please provide the proposal ID.")
			os.Exit(2)
		}

		gov.SubmitAllVotesForProposal(bin, proposalID)
	} else {
		targetVersion := os.Args[1]
		if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d(-rc\d+)?`, targetVersion); !matched {
			log.Println("Invalid target version. Please use the format vX.Y.Z(-rc*).")
			os.Exit(2)
		}

		upgradeLocalNode(bin, targetVersion)
	}
}

// upgradeLocalNode prepares upgrading the local node to the target version
// by submitting the upgrade proposal and voting on it using all testing accounts.
func upgradeLocalNode(bin *utils.Binary, targetVersion string) {
	currentHeight, err := utils.GetCurrentHeight(bin)
	if err != nil {
		log.Fatalf("Error getting current height: %v", err)
	}

	upgradeHeight := currentHeight + deltaHeight

	log.Println("Submitting upgrade proposal...")

	proposalID, err := gov.SubmitUpgradeProposal(bin, targetVersion, upgradeHeight)
	if err != nil {
		log.Fatalf("Error executing upgrade proposal: %v", err)
	}

	log.Printf("Scheduled upgrade to %s at height %d.\n", targetVersion, upgradeHeight)

	gov.SubmitAllVotesForProposal(bin, proposalID)
}
