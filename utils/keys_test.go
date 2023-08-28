package utils_test

import (
	"testing"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	"github.com/stretchr/testify/require"
)

func TestParseKeysFromOut(t *testing.T) {
	testcases := []struct {
		name        string
		out         string
		expKeys     []string
		expError    bool
		errContains string
	}{
		{
			name: "pass",
			//nolint lll -- line length is okay here
			out:     `[{"name":"dev0","type":"local","address":"evmos16qljjgus9zevcxdjscuf502zy6en427nty78c0","pubkey":"{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A7YjISvuApMJ/OGKVifuVqrUnJYryXPcVAR5zPzP5yz5\"}"},{"name":"dev1","type":"local","address":"evmos16cqwxv4hcqpzc7zd9fd4pw3jr4yf9jxrfr6tj0","pubkey":"{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A+VsC7GstX+ItZDKvWSmbQrjuvmZ0GenWB46Pi6F0fwL\"}"},{"name":"dev2","type":"local","address":"evmos1ecamqksjl7erx89lextmru88mpy669psjcehlz","pubkey":"{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"Aha/x6t6Uaiw+md5F4XjaPleHTw6toUU9egkWCPm50wk\"}"},{"name":"testKey","type":"local","address":"evmos17slw9hdyxvxypzsdwj9vjg7uedhfw26ksqydye","pubkey":"{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"ApDf/TgsVwangM3CciQuAoIgBvo5ZXxPHkA7K2XpeAae\"}"}]`,
			expKeys: []string{"dev0", "dev1", "dev2", "testKey"},
		},
		{
			name:        "fail - no keys",
			out:         "invalid output",
			expError:    true,
			errContains: "error unmarshalling keys",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			accounts, err := utils.ParseAccountsFromOut(tc.out)
			if tc.expError {
				require.Error(t, err, "expected error parsing accounts")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error parsing accounts")

				var keys []string
				for _, account := range accounts {
					keys = append(keys, account.Name)
				}
				require.Equal(t, tc.expKeys, keys, "expected different keys")
			}
		})
	}
}

func TestParseDelegationsFromResponse(t *testing.T) {
	testcases := []struct {
		name        string
		out         string
		expVals     []string
		expError    bool
		errContains string
	}{
		{
			name:    "pass",
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
		t.Run(tc.name, func(t *testing.T) {
			delegations, err := utils.ParseDelegationsFromResponse(tc.out)
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
