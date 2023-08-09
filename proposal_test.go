package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProposalID(t *testing.T) {
	testcases := []struct {
		name     string
		out      string
		expID    int
		expError bool
	}{
		{
			name: "pass",
			out: `gas estimate: 850456
			code: 0
			codespace: ""
			data: 12330A2D2F636F736D6F732E676F762E763162657461312E4D73675375626D697450726F706F73616C526573706F6E736512020804
			events:
			logs:
			- events:
			  - attributes:
				- key: amount
				  value: 1000000000000aevmos
				- key: proposal_id
				  value: "4"
				type: proposal_deposit
			  - attributes:
				- key: proposal_id
				  value: "4"
				- key: proposal_messages
				  value: ',/cosmos.gov.v1.MsgExecLegacyContent'
				- key: voting_period_start
				  value: "4"
				type: submit_proposal
				type: transfer
			  log: ""
			  msg_index: 0
			timestamp: ""
			tx: null
			txhash: A505158FF9EFB4E939CD4A9A94F731E0E34AEEF50C7E53A723226EEF33A1A89B`,
			expID: 4,
		},
		{
			name:     "fail - no proposal ID",
			out:      "invalid output",
			expError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := getProposalID(tc.out)
			if tc.expError {
				require.Error(t, err, "expected error parsing proposal ID")
			} else {
				require.NoError(t, err, "unexpected error parsing proposal ID")
				require.Equal(t, tc.expID, id, "expected different proposal ID")
			}
		})
	}
}
