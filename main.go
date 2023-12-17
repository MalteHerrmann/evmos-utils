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

func main() {
	if len(os.Args) < 2 {
		log.Printf(
			"Possible usages:\n" +
				"  upgrade-local-node-go <target_version>\n" +
				"  upgrade-local-node-go vote [proposal-id]\n" +
				"  upgrade-local-node-go deposit [proposal-id]\n",
		)
		os.Exit(1)
	}

	bin, err := utils.NewEvmosTestingBinary()
	if err != nil {
		log.Fatalf("error creating binary: %v", err)
	}

	if err = bin.GetAccounts(); err != nil {
		log.Fatalf("error getting accounts: %v", err)
	}

	// TODO: use with Cobra CLI
	switch os.Args[1] {
	case "vote":
		proposalID, err := getProposalIDFromInput(bin, os.Args)
		if err != nil {
			log.Fatalf("error getting proposal ID: %v", err)
		}

		if err = gov.SubmitAllVotesForProposal(bin, proposalID); err != nil {
			log.Fatalf("error submitting votes for proposal %d: %v", proposalID, err)
		}

	case "deposit":
		deposit, err := gov.GetMinDeposit(bin)
		if err != nil {
			log.Fatalf("error getting minimum deposit: %v", err)
		}

		proposalID, err := getProposalIDFromInput(bin, os.Args)
		if err != nil {
			log.Fatalf("error getting proposal ID: %v", err)
		}

		if _, err = gov.DepositForProposal(bin, proposalID, bin.Accounts[0].Name, deposit.String()); err != nil {
			log.Fatalf("error depositing for proposal %d: %v", proposalID, err)
		}

	default:
		targetVersion := os.Args[1]
		if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d(-rc\d+)?`, targetVersion); !matched {
			log.Fatalf("invalid target version: %s; please use the format vX.Y.Z(-rc*).\n", targetVersion)
		}

		if err := upgradeLocalNode(bin, targetVersion); err != nil {
			log.Fatalf("error upgrading local node: %v", err)
		}
	}
}

// getProposalIDFromInput gets the proposal ID from the command line arguments.
func getProposalIDFromInput(bin *utils.Binary, args []string) (int, error) {
	var (
		err        error
		proposalID int
	)

	switch len(args) {
	case 2:
		proposalID, err = gov.QueryLatestProposalID(bin)
		if err != nil {
			return 0, errors.Wrap(err, "error querying latest proposal ID")
		}
	case 3:
		proposalID, err = strconv.Atoi(args[2])
		if err != nil {
			return 0, errors.Wrapf(err, "error converting proposal ID %s to integer", args[2])
		}
	default:
		return 0, errors.New("invalid number of arguments")
	}

	return proposalID, nil
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
