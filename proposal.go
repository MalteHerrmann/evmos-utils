package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// submitUpgradeProposal submits a software upgrade proposal with the given target version and upgrade height.
func submitUpgradeProposal(targetVersion string, upgradeHeight int) (int, error) {
	upgradeProposal := buildUpgradeProposalCommand(targetVersion, upgradeHeight)
	out, err := executeShellCommand(upgradeProposal, evmosdHome, "dev0", true, false)
	if err != nil {
		return 0, err
	}

	return getProposalID(out)
}

// getProposalID parses the proposal ID from the given output from submitting an upgrade proposal.
func getProposalID(out string) (int, error) {
	// Define the regular expression pattern
	pattern := `- key:\s*proposal_id\s*\n\s*value:\s*"([^"]+)"`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(out)
	if len(match) != 2 {
		return 0, fmt.Errorf("proposal ID not found in output")
	}

	return strconv.Atoi(match[1])
}

// buildUpgradeProposalCommand builds the command to submit a software upgrade proposal.
func buildUpgradeProposalCommand(targetVersion string, upgradeHeight int) []string {
	return []string{
		"tx", "gov", "submit-legacy-proposal", "software-upgrade", targetVersion,
		"--title", fmt.Sprintf("'Upgrade to %s'", targetVersion),
		"--description", fmt.Sprintf("'Upgrade to %s'", targetVersion),
		"--upgrade-height", fmt.Sprintf("%d", upgradeHeight),
		"--deposit", "100000000000000000000aevmos",
		"--no-validate",
	}
}

// voteForProposal votes for the proposal with the given ID using the given account.
func voteForProposal(proposalID int, sender string) error {
	_, err := executeShellCommand(
		[]string{"tx", "gov", "vote", fmt.Sprintf("%d", proposalID), "yes"},
		evmosdHome,
		sender,
		true, true,
	)
	return err
}
