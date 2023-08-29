package utils_test

import (
	"testing"

	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetTxHashFromResponse(t *testing.T) {
	t.Parallel()

	cdc, ok := utils.GetCodec()
	require.True(t, ok, "unexpected error getting codec")

	testcases := []struct {
		name        string
		out         string
		expHash     string
		expError    bool
		errContains string
	}{
		{
			name: "pass - successful tx",
			//nolint:lll // line length is okay here
			out:     `{"height":"0","txhash":"F9C69496731969BDC3C03E5D65612AB07E09809E0BDC753A2758B6E70C92FD74","codespace":"","code":0,"data":"","raw_log":"[]","logs":[],"info":"","gas_wanted":"0","gas_used":"0","tx":null,"timestamp":"","events":[]}`,
			expHash: "F9C69496731969BDC3C03E5D65612AB07E09809E0BDC753A2758B6E70C92FD74",
		},
		{
			name: "fail - unsuccessful tx",
			//nolint:lll // line length is okay here
			out:         `{"height":"0","txhash":"F9C69496731969BDC3C03E5D65612AB07E09809E0BDC753A2758B6E70C92FD74","codespace":"","code":1,"data":"","raw_log":"[]","logs":[],"info":"","gas_wanted":"0","gas_used":"0","tx":null,"timestamp":"","events":[]}`,
			expError:    true,
			errContains: "transaction failed with code",
		},
		{
			name:        "fail - no tx hash",
			out:         "invalid output",
			expError:    true,
			errContains: "error unpacking transaction hash from json",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hash, err := utils.GetTxHashFromTxResponse(cdc, tc.out)
			if tc.expError {
				require.Error(t, err, "expected error getting tx hash")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error getting tx hash")
				require.Equal(t, tc.expHash, hash, "expected different transaction hash")
			}
		})
	}
}

func TestGetEventsFromTxResponse(t *testing.T) {
	t.Parallel()

	cdc, ok := utils.GetCodec()
	require.True(t, ok, "unexpected error getting codec")

	testcases := []struct {
		name        string
		out         string
		expEvents   []sdk.StringEvent
		expError    bool
		errContains string
	}{
		{
			name: "pass",
			//nolint:lll // line length is okay here
			out: `{"height":"138","txhash":"FE14C1BF8BBA55A314D7040ACA404A97D2172126ABF81C0C90D0B5C9B0CADEE6","codespace":"","code":0,"data":"12330A2D2F636F736D6F732E676F762E763162657461312E4D73675375626D697450726F706F73616C526573706F6E736512020805","raw_log":"[{\"msg_index\":0,\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/cosmos.gov.v1beta1.MsgSubmitProposal\"},{\"key\":\"sender\",\"value\":\"evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2\"},{\"key\":\"module\",\"value\":\"gov\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_id\",\"value\":\"5\"},{\"key\":\"proposal_messages\",\"value\":\",/cosmos.gov.v1.MsgExecLegacyContent\"}]},{\"type\":\"coin_spent\",\"attributes\":[{\"key\":\"spender\",\"value\":\"evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2\"},{\"key\":\"amount\",\"value\":\"100000000000000000000aevmos\"}]},{\"type\":\"coin_received\",\"attributes\":[{\"key\":\"receiver\",\"value\":\"evmos10d07y265gmmuvt4z0w9aw880jnsr700jcrztvm\"},{\"key\":\"amount\",\"value\":\"100000000000000000000aevmos\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"evmos10d07y265gmmuvt4z0w9aw880jnsr700jcrztvm\"},{\"key\":\"sender\",\"value\":\"evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2\"},{\"key\":\"amount\",\"value\":\"100000000000000000000aevmos\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"sender\",\"value\":\"evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2\"}]},{\"type\":\"proposal_deposit\",\"attributes\":[{\"key\":\"amount\",\"value\":\"100000000000000000000aevmos\"},{\"key\":\"proposal_id\",\"value\":\"5\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"voting_period_start\",\"value\":\"5\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"submit_proposal","attributes":[{"key":"proposal_id","value":"5"},{"key":"proposal_messages","value":",/cosmos.gov.v1.MsgExecLegacyContent"}]}]}],"info":"","gas_wanted":"270887","gas_used":"209242","tx":{"@type":"/cosmos.tx.v1beta1.Tx","body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgSubmitProposal","content":{"@type":"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal","title":"'Upgrade to v14.0.0-rc4'","description":"'Upgrade to v14.0.0-rc4'","plan":{"name":"v14.0.0-rc4","time":"0001-01-01T00:00:00Z","height":"151","info":"","upgraded_client_state":null}},"initial_deposit":[{"denom":"aevmos","amount":"100000000000000000000"}],"proposer":"evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9i4S3tyqlhjlLyOwVu9PNYWZNLL29RR643ae7K63VJ1"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"5"}],"fee":{"amount":[{"denom":"aevmos","amount":"1000000000000000000"}],"gas_limit":"270887","payer":"","granter":""},"tip":null},"signatures":["cOwD3//F3MuL1rBmlCRTySjJidZGJygW2u96nXL5DR8a1GkmjyfToUUUpfs2HtYZohIGOsjG8werUSxllSa6LgA="]},"timestamp":"2023-08-23T21:16:24Z","events":[{"type":"submit_proposal","attributes":[{"key":"proposal_id","value":"5","index":true},{"key":"proposal_messages","value":",/cosmos.gov.v1.MsgExecLegacyContent","index":true}]}]}`,
			expEvents: []sdk.StringEvent{{
				Type: "submit_proposal",
				Attributes: []sdk.Attribute{
					{Key: "proposal_id", Value: "5"},
					{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
				},
			}},
		},
		{
			name:        "fail - invalid output",
			out:         "invalid output",
			expError:    true,
			errContains: "error unmarshalling transaction response",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			events, err := utils.GetEventsFromTxResponse(cdc, tc.out)
			if tc.expError {
				require.Error(t, err, "expected error getting tx events")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error getting tx events")
				require.Equal(t, tc.expEvents, events, "expected different transaction events")
			}
		})
	}
}
