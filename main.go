package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
)

// The amount of blocks in the future that the upgrade will be scheduled.
const deltaHeight = 15

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
	currentHeight, err := utils.GetCurrentHeight()
	if err != nil {
		log.Fatalf("Error getting current height: %v", err)
	}

	upgradeHeight := currentHeight + deltaHeight
	fmt.Println("Submitting upgrade proposal...")
	proposalID, err := gov.SubmitUpgradeProposal(targetVersion, upgradeHeight)
	if err != nil {
		log.Fatalf("Error executing upgrade proposal: %v", err)
	}
	fmt.Printf("Scheduled upgrade to %s at height %d.\n", targetVersion, upgradeHeight)

	availableAccounts, err := utils.GetAccounts()
	if err != nil {
		log.Fatalf("Error getting available keys: %v", err)
	}

	accsWithDelegations, err := utils.FilterAccountsWithDelegations(availableAccounts)
	if err != nil {
		log.Fatalf("Error filtering accounts: %v", err)
	}
	wait(1)

	fmt.Println("Voting for upgrade...")
	for _, acc := range accsWithDelegations {
		if err = gov.VoteForProposal(proposalID, acc.Name); err != nil {
			fmt.Printf("  - could NOT vote using key: %s\n", acc.Name)
		} else {
			fmt.Printf("  - voted using key: %s\n", acc.Name)
		}
	}
}

// wait waits for the specified amount of seconds.
func wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
