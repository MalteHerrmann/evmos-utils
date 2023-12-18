package gov

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/pkg/errors"
)

// SubmitAllVotes submits a vote for the given proposal ID using all testing accounts.
func SubmitAllVotes(bin *utils.Binary, args []string) error {
	proposalID, err := GetProposalIDFromInput(bin, args)
	if err != nil {
		return err
	}

	return SubmitAllVotesForProposal(bin, proposalID)
}

// SubmitAllVotesForProposal submits a vote for the given proposal ID using all testing accounts.
func SubmitAllVotesForProposal(bin *utils.Binary, proposalID int) error {
	accsWithDelegations, err := utils.FilterAccountsWithDelegations(bin)
	if err != nil {
		return errors.Wrap(err, "error filtering accounts")
	}

	if len(accsWithDelegations) == 0 {
		return errors.New("no accounts with delegations found")
	}

	utils.Wait(1)
	log.Printf("Voting for proposal %d...\n", proposalID)

	var out string

	for _, acc := range accsWithDelegations {
		out, err = VoteForProposal(bin, proposalID, acc.Name)
		if err != nil {
			if strings.Contains(out, fmt.Sprintf("%d: unknown proposal", proposalID)) {
				return fmt.Errorf("no proposal with ID %d found", proposalID)
			}

			if strings.Contains(out, fmt.Sprintf("%d: inactive proposal", proposalID)) {
				return fmt.Errorf("proposal with ID %d is inactive", proposalID)
			}

			log.Printf("  - could NOT vote using key: %s\n", acc.Name)
		} else {
			log.Printf("  - voted using key: %s\n", acc.Name)
		}
	}

	return nil
}

// VoteForProposal votes for the proposal with the given ID using the given account.
func VoteForProposal(bin *utils.Binary, proposalID int, sender string) (string, error) {
	out, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand:  []string{"tx", "gov", "vote", strconv.Itoa(proposalID), "yes"},
		From:        sender,
		UseDefaults: true,
		Quiet:       true,
	})
	if err != nil {
		return out, errors.Wrap(err, fmt.Sprintf("failed to vote for proposal %d", proposalID))
	}

	return out, nil
}
