package gov

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/pkg/errors"
)

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

	return GetProposalID(events)
}

// GetProposalID looks for the proposal submission event in the given transaction events
// and returns the proposal id, if found.
func GetProposalID(events []abcitypes.Event) (int, error) {
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

// VoteForProposal votes for the proposal with the given ID using the given account.
func VoteForProposal(bin *utils.Binary, proposalID int, sender string) error {
	_, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand:  []string{"tx", "gov", "vote", fmt.Sprintf("%d", proposalID), "yes"},
		From:        sender,
		UseDefaults: true,
		Quiet:       true,
	})

	return errors.Wrap(err, fmt.Sprintf("failed to vote for proposal %d", proposalID))
}
