package gov

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	abcitypes "github.com/cometbft/cometbft/abci/types"
)

// SubmitUpgradeProposal submits a software upgrade proposal with the given target version and upgrade height.
func SubmitUpgradeProposal(targetVersion string, upgradeHeight int) (int, error) {
	upgradeProposal := buildUpgradeProposalCommand(targetVersion, upgradeHeight)
	out, err := utils.ExecuteBinaryCmd(utils.BinaryCmdArgs{
		Subcommand:  upgradeProposal,
		From:        "dev0",
		UseDefaults: true,
	})
	if err != nil {
		return 0, err
	}

	// Clean gas estimate output and only leave json output
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	out = lines[len(lines)-1] // last line is json output

	events, err := utils.GetTxEvents(out)
	if err != nil {
		return 0, fmt.Errorf("error getting tx events: %w", err)
	}

	return getProposalID(events)
}

// getProposalID looks for the proposal submission event in the given transaction events
// and returns the proposal id, if found.
func getProposalID(events []abcitypes.Event) (int, error) {
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
func VoteForProposal(proposalID int, sender string) error {
	_, err := utils.ExecuteBinaryCmd(utils.BinaryCmdArgs{
		Subcommand:  []string{"tx", "gov", "vote", fmt.Sprintf("%d", proposalID), "yes"},
		From:        sender,
		UseDefaults: true,
		Quiet:       true,
	})
	return err
}
