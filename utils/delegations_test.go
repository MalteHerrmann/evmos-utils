package utils_test

import (
	"testing"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/stretchr/testify/require"
)

func TestParseDelegationsFromResponse(t *testing.T) {
	t.Parallel()

	cdc, ok := utils.GetCodec()
	require.True(t, ok, "unexpected error getting codec")

	testcases := []struct {
		name        string
		out         string
		expVals     []string
		expError    bool
		errContains string
	}{
		{
			name: "pass",
			//nolint:lll // line length is okay here
			out:     `{"delegation_responses":[{"delegation":{"delegator_address":"evmos1v6jyld5mcu37d3dfe7kjrw0htkc4wu2mxn9y25","validator_address":"evmosvaloper1v6jyld5mcu37d3dfe7kjrw0htkc4wu2mta25tf","shares":"1000000000000000000000.000000000000000000"},"balance":{"denom":"aevmos","amount":"1000000000000000000000"}}],"pagination":{"next_key":null,"total":"0"}}`,
			expVals: []string{"evmosvaloper1v6jyld5mcu37d3dfe7kjrw0htkc4wu2mta25tf"},
		},
		{
			name:        "fail - no keys",
			out:         "invalid output",
			expError:    true,
			errContains: "error unmarshalling delegations",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			delegations, err := utils.ParseDelegationsFromResponse(cdc, tc.out)
			if tc.expError {
				require.Error(t, err, "expected error parsing delegations")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error parsing delegations")

				var vals []string
				for _, delegation := range delegations {
					vals = append(vals, delegation.ValidatorAddress)
				}
				require.Equal(t, tc.expVals, vals, "expected different validators")
			}
		})
	}
}
