package gov

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/pkg/errors"
)

// buildUpgradeProposalCommand builds the command to submit a software upgrade proposal.
func buildUpgradeProposalCommand(targetVersion string, upgradeHeight int) []string {
	return []string{
		"tx", "gov", "submit-legacy-proposal", "software-upgrade", targetVersion,
		"--title", fmt.Sprintf("'Upgrade to %s'", targetVersion),
		"--description", fmt.Sprintf("'Upgrade to %s'", targetVersion),
		"--upgrade-height", fmt.Sprintf("%d", upgradeHeight),
		"--deposit", "100000000000000000000aevmos",
		"--output", "json",
		"--no-validate",
	}
}

// GetProposalIDFromSubmitEvents looks for the proposal submission event in the given transaction events
// and returns the proposal id, if found.
func GetProposalIDFromSubmitEvents(events []sdk.StringEvent) (int, error) {
	for _, event := range events {
		if event.Type != "submit_proposal" {
			continue
		}

		for _, attribute := range event.Attributes {
			if attribute.Key == "proposal_id" {
				proposalID, err := strconv.Atoi(attribute.Value)
				if err != nil {
					return 0, fmt.Errorf("error parsing proposal id: %w", err)
				}

				return proposalID, nil
			}
		}
	}

	return 0, fmt.Errorf("proposal submission event not found")
}

// QueryLatestProposalID queries the latest proposal ID.
func QueryLatestProposalID(bin *utils.Binary) (int, error) {
	out, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand: []string{"q", "gov", "proposals", "--output=json"},
		Quiet:      true,
	})
	if err != nil {
		if strings.Contains(out, "no proposals found") {
			return 0, errors.New("no proposals found")
		}

		return 0, errors.Wrap(err, "error querying proposals")
	}

	// NOTE: the SDK CLI command uses the x/gov v1 package
	// see: https://github.com/cosmos/cosmos-sdk/blob/v0.47.4/x/gov/client/cli/query.go#L151-L159
	var res govv1types.QueryProposalsResponse

	err = bin.Cdc.UnmarshalJSON([]byte(out), &res)
	if err != nil {
		return 0, errors.Wrap(err, "error unmarshalling proposals")
	}

	if len(res.Proposals) == 0 {
		return 0, errors.New("no proposals found")
	}

	return int(res.Proposals[len(res.Proposals)-1].Id), nil
}

// SubmitAllVotesForProposal submits a vote for the given proposal ID using all testing accounts.
func SubmitAllVotesForProposal(bin *utils.Binary, proposalID int) error {
	accsWithDelegations, err := utils.FilterAccountsWithDelegations(bin)
	if err != nil {
		return errors.Wrap(err, "Error filtering accounts")
	}

	if len(accsWithDelegations) == 0 {
		return errors.New("No accounts with delegations found")
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

// SubmitUpgradeProposal submits a software upgrade proposal with the given target version and upgrade height.
func SubmitUpgradeProposal(bin *utils.Binary, targetVersion string, upgradeHeight int) (int, error) {
	upgradeProposal := buildUpgradeProposalCommand(targetVersion, upgradeHeight)

	out, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand:  upgradeProposal,
		From:        "dev0",
		UseDefaults: true,
	})
	if err != nil {
		return 0, errors.Wrap(err,
			fmt.Sprintf("failed to submit upgrade proposal to %s at height %d", targetVersion, upgradeHeight),
		)
	}

	// Clean gas estimate output and only leave json output
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	out = lines[len(lines)-1] // last line is json output

	events, err := utils.GetTxEvents(bin, out)
	if err != nil {
		return 0, fmt.Errorf("error getting tx events: %w", err)
	}

	if len(events) == 0 {
		return 0, fmt.Errorf("no events found in transaction to submit proposal")
	}

	return GetProposalIDFromSubmitEvents(events)
}

// VoteForProposal votes for the proposal with the given ID using the given account.
func VoteForProposal(bin *utils.Binary, proposalID int, sender string) (string, error) {
	out, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand:  []string{"tx", "gov", "vote", fmt.Sprintf("%d", proposalID), "yes"},
		From:        sender,
		UseDefaults: true,
		Quiet:       true,
	})
	if err != nil {
		return out, errors.Wrap(err, fmt.Sprintf("failed to vote for proposal %d", proposalID))
	}

	return out, nil
}
