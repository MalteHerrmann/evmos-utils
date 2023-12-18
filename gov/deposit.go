package gov

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// Deposit deposits the minimum needed deposit for a given governance proposal.
func Deposit(bin *utils.Binary, args []string) error {
	deposit, err := GetMinDeposit(bin)
	if err != nil {
		log.Fatalf("error getting minimum deposit: %v", err)
	}

	proposalID, err := GetProposalIDFromInput(bin, args)
	if err != nil {
		log.Fatalf("error getting proposal ID: %v", err)
	}

	return DepositForProposal(
		bin, proposalID, bin.Accounts[0].Name, deposit.String(),
	)
}

// DepositForProposal deposits the given amount for the proposal with the given proposalID
// from the given account.
func DepositForProposal(bin *utils.Binary, proposalID int, sender, deposit string) error {
	_, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand: []string{
			"tx", "gov", "deposit", strconv.Itoa(proposalID), deposit,
		},
		From:        sender,
		UseDefaults: true,
		Quiet:       true,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to deposit for proposal %d", proposalID))
	}

	return nil
}

// GetMinDeposit returns the minimum deposit necessary for a proposal from the governance parameters of
// the running chain.
func GetMinDeposit(bin *utils.Binary) (sdk.Coins, error) {
	out, err := utils.ExecuteBinaryCmd(bin, utils.BinaryCmdArgs{
		Subcommand: []string{"q", "gov", "param", "deposit", "--output=json"},
		Quiet:      true,
	})
	if err != nil {
		return sdk.Coins{}, errors.Wrap(err, "failed to query governance parameters")
	}

	return ParseMinDepositFromResponse(out)
}

// ParseMinDepositFromResponse parses the minimum deposit from the given output of the governance
// parameters query.
//
// FIXME: It wasn't possible to unmarshal the JSON output of the query because of a missing unit in the max_deposit_period
// parameter. This should rather be done using GRPC.
func ParseMinDepositFromResponse(out string) (sdk.Coins, error) {
	depositPatternRaw := `min_deposit":\[{"denom":"(\w+)","amount":"(\d+)`
	depositPattern := regexp.MustCompile(depositPatternRaw)

	minDepositMatch := depositPattern.FindStringSubmatch(out)
	if len(minDepositMatch) == 0 {
		return sdk.Coins{}, fmt.Errorf("failed to find min deposit in params output: %q", out)
	}

	minDepositDenom := minDepositMatch[1]

	minDepositAmount, err := strconv.Atoi(minDepositMatch[2])
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("failed to find min deposit in params output: %q", out)
	}

	return sdk.Coins{sdk.NewInt64Coin(minDepositDenom, int64(minDepositAmount))}, nil
}
