package gov

import (
	"fmt"
	"strconv"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/pkg/errors"
)

// DepositForProposal deposits the given amount of Evmos for the proposal with the given proposalID
// from the given account.
func DepositForProposal(bin *utils.Binary, proposalID int, sender string, amount int) (string, error) {
	out, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand: []string{
			"tx", "gov", "deposit", strconv.Itoa(proposalID), strconv.Itoa(amount) + "aevmos",
		},
		From:        sender,
		UseDefaults: true,
		Quiet:       true,
	})
	if err != nil {
		return out, errors.Wrap(err, fmt.Sprintf("failed to deposit for proposal %d", proposalID))
	}

	return out, nil
}
