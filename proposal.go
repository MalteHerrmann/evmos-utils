package main

import (
	"fmt"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/evmos/evmos/v14/app"
	"github.com/evmos/evmos/v14/encoding"
	"strconv"
	"strings"
)

var (
	// cdc is the codec to be used for the client
	cdc = encodingConfig.Codec
	// encodingConfig specifies the encoding configuration to be used for the client
	encodingConfig = encoding.MakeConfig(app.ModuleBasics)
)

// submitUpgradeProposal submits a software upgrade proposal with the given target version and upgrade height.
func submitUpgradeProposal(targetVersion string, upgradeHeight int) (int, error) {
	upgradeProposal := buildUpgradeProposalCommand(targetVersion, upgradeHeight)
	out, err := executeShellCommand(upgradeProposal, evmosdHome, "dev0", true)
	if err != nil {
		return 0, err
	}

	// Clean gas estimate output and only leave json output
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	out = lines[len(lines)-1] // last line is json output

	events, err := getTxEvents(out)
	if err != nil {
		panic(err)
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

// voteForProposal votes for the proposal with the given ID using the given account.
func voteForProposal(proposalID int, sender string) error {
	_, err := executeShellCommand([]string{"tx", "gov", "vote", fmt.Sprintf("%d", proposalID), "yes"}, evmosdHome, sender, true)
	return err
}
