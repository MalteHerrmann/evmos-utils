package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

const (
	// The chain ID of the node that will be upgraded.
	chainID = "evmos_9000-1"
	// The home directory of the node that will be upgraded.
	evmosdHome = "~/.tmp-evmosd"
	// The amount of blocks in the future that the upgrade will be scheduled.
	deltaHeight = 15
	// The amount of fees to be sent with a default transaction.
	defaultFees int = 1e18 // 1 aevmos
	// The denomination used for the local node.
	denom = "aevmos"
)

var (
	// The default flags that will be used for all commands related to governance.
	defaultFlags = []string{
		"--chain-id", chainID,
		"--keyring-backend", "test",
		"--gas", "auto",
		"--fees", fmt.Sprintf("%d%s", defaultFees, denom),
		"--gas-adjustment", "1.3",
		"-b", "block",
		"-y",
	}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run upgrade-local-node.go <target_version>")
		os.Exit(1)
	}

	targetVersion := os.Args[1]
	if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d(-rc\d+)?`, targetVersion); !matched {
		fmt.Println("Invalid target version. Please use the format vX.Y.Z(-rc*).")
		os.Exit(2)
	}

	upgradeLocalNode(targetVersion)
}

// upgradeLocalNode prepares upgrading the local node to the target version
// by submitting the upgrade proposal and voting on it using all testing accounts.
func upgradeLocalNode(targetVersion string) {
	currentHeight, err := getCurrentHeight()
	if err != nil {
		log.Fatalf("Error getting current height: %v", err)
	}

	upgradeHeight := currentHeight + deltaHeight
	upgradeProposal := buildUpgradeProposalCommand(targetVersion, upgradeHeight)
	_, err = executeShellCommand(upgradeProposal, evmosdHome, "dev0", true)
	if err != nil {
		log.Fatalf("Error executing upgrade proposal: %v", err)
	}

	fmt.Printf("Scheduled upgrade to %s at height %d.\n", targetVersion, upgradeHeight)
	wait(2)
	if err = voteForProposal(1, "dev0"); err != nil {
		log.Fatalf("Error voting for upgrade: %v", err)
	}

	wait(2)
	if err = voteForProposal(1, "dev1"); err != nil {
		log.Fatalf("Error voting for upgrade: %v", err)
	}

	wait(2)
	if err = voteForProposal(1, "dev2"); err != nil {
		log.Fatalf("Error voting for upgrade: %v", err)
	}
}

// voteForProposal votes for the proposal with the given ID using the given account.
func voteForProposal(proposalID int, sender string) error {
	_, err := executeShellCommand([]string{"tx", "gov", "vote", fmt.Sprintf("%d", proposalID), "yes"}, evmosdHome, sender, true)
	return err
}

// wait waits for the specified amount of seconds.
func wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
