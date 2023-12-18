package utils_test

import (
	"testing"

	"github.com/MalteHerrmann/evmos-utils/utils"
	"github.com/stretchr/testify/require"
)

func TestParseKeysFromOut(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		out         string
		expKeys     []string
		expError    bool
		errContains string
	}{
		{
			name: "pass",
			//nolint:lll // line length is okay here
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
