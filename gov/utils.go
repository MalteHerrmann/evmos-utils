package gov

import (
	"fmt"
	"strconv"

	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/pkg/errors"
)

// GetProposalIDFromInput gets the proposal ID from the command line arguments.
func GetProposalIDFromInput(bin *utils.Binary, args []string) (int, error) {
	var (
		err        error
		proposalID int
	)

	switch len(args) {
	case 0:
		proposalID, err = QueryLatestProposalID(bin)
		if err != nil {
			return 0, errors.Wrap(err, "error querying latest proposal ID")
		}
	case 1:
		proposalID, err = strconv.Atoi(args[0])
		if err != nil {
			return 0, errors.Wrapf(err, "error converting proposal ID %s to integer", args[0])
		}
	default:
		return 0, fmt.Errorf("invalid number of arguments; expected 0 or 1; got %d", len(args))
	}

	return proposalID, nil
}
