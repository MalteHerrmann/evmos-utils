package gov_test

import (
	"testing"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParseMinDepositFromResponse(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		out           string
		expMinDeposit sdk.Coins
		expError      bool
		errContains   string
	}{
		{
			name:          "pass",
			out:           `{"min_deposit":[{"denom":"aevmos","amount":"10000000"}],"max_deposit_period":"30000000000"}`,
			expMinDeposit: sdk.Coins{sdk.NewInt64Coin("aevmos", 10000000)},
		},
		{
			name:        "fail - no min deposit",
			out:         "invalid output",
			expError:    true,
			errContains: "failed to find min deposit in params output",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			minDeposit, err := gov.ParseMinDepositFromResponse(tc.out)
			if tc.expError {
				require.Error(t, err, "expected error parsing min deposit")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error parsing min deposit")
				require.Equal(t, tc.expMinDeposit, minDeposit, "expected different min deposit")
			}
		})
	}
}
