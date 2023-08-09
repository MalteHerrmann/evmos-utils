package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseKeysFromOut(t *testing.T) {
	testcases := []struct {
		name     string
		out      string
		expKeys  []string
		expError bool
	}{
		{
			name: "pass",
			out: `  - address: evmos19mx9kcksequm4m4xume5h0k9fquwgmea3yvu89
						name: dev0
						pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"AmquZBW+CPcgHKx6D4YRDICzr0MNcRvl9Wm/jJn8wJxs"}'
						type: local
					- address: evmos18z7xfs864u49jcv6gkgajpteesjl5d7krpple6
						name: dev1
						pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"AtY/rqJrmhKbXrQ02xSxq/t9JGgbP2T7HPGTZJIbuT8I"}'
						type: local
					- address: evmos12rrt7vcnxvhxad6gzz0vt5psdlnurtldety57n
						name: dev2
						pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A544btlGjv4zB/qpWT8dQqlAHrcmgZEvrFSgJnp7Yjt4"}'
						type: local
					- address: evmos1dln2gjtsfd2sny6gwdxzyxcsr0uu8sh5nwajun
						name: testKey1
						pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"Amja5pRiVw+5vPkozo6Eo20AEbYVVBqOKBi5yP7EbxyJ"}'
						type: local
					- address: evmos1qdxgxz9g2la8g9eyjdq4srlpxgrmuqd6ty88zm
						name: testKey2
						pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A+ytKfWmkQiW0c6iOCXSL71e4b5njmJVUd1msONsPEnA"}'
						type: local
					- address: evmos1hduvvhjvu0pqu7m97pajymdsupqx3us3ntey9a
						name: testKey3
						pubkey: '{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"AsdAPndEVttzhUz5iSm0/FoFxkzB0oZE7DuKf3NjzXkS"}'
						type: local`,
			expKeys: []string{"dev0", "dev1", "dev2", "testKey1", "testKey2", "testKey3"},
		},
		{
			name:     "fail - no keys",
			out:      "invalid output",
			expError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			keys, err := parseKeysFromOut(tc.out)
			if tc.expError {
				require.Error(t, err, "expected error parsing keys")
			} else {
				require.NoError(t, err, "unexpected error parsing keys")
				require.Equal(t, tc.expKeys, keys)
			}
		})
	}
}
