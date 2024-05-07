package gov

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/pkg/errors"
)

// SubmitAllVotes submits a vote for the given proposal ID using all testing accounts.
func SubmitAllVotes(bin *utils.Binary, args []string) (int, error) {
	proposalID, err := GetProposalIDFromInput(bin, args)
	if err != nil {
		return 0, err
	}

	return proposalID, SubmitAllVotesForProposal(bin, proposalID)
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

	if err := utils.WaitNBlocks(bin, 1); err != nil {
		return errors.Wrapf(err, "error waiting for blocks")
	}

	bin.Logger.Info().Msgf("voting for proposal %d", proposalID)

	var (
		out             string
		successfulVotes int
	)

	for _, acc := range accsWithDelegations {
		out, err = VoteForProposal(bin, proposalID, acc.Name)
		if err != nil {
			if strings.Contains(out, fmt.Sprintf("%d: unknown proposal", proposalID)) {
				return fmt.Errorf("no proposal with ID %d found", proposalID)
			}

			if strings.Contains(out, fmt.Sprintf("%d: inactive proposal", proposalID)) {
				return fmt.Errorf("proposal with ID %d is inactive", proposalID)
			}

			bin.Logger.Error().Msgf("could not vote using key %s: %v", acc.Name, err)
		} else {
			bin.Logger.Info().Msgf("voted using key %s", acc.Name)

			successfulVotes++
		}
	}

	if successfulVotes == 0 {
		return errors.New("there were no successful votes for the proposal, please check logs")
	}

	return nil
}

// VoteForProposal votes for the proposal with the given ID using the given account.
func VoteForProposal(bin *utils.Binary, proposalID int, sender string) (string, error) {
	out, err := utils.ExecuteTx(bin, utils.TxArgs{
		Subcommand: []string{"tx", "gov", "vote", strconv.Itoa(proposalID), "yes"},
		From:       sender,
		Quiet:      true,
	})
	if err != nil {
		return out, errors.Wrap(err, fmt.Sprintf("failed to vote for proposal %d", proposalID))
	}

	return out, nil
}
