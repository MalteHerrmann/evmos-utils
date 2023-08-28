package gov

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

func TestGetProposalID(t *testing.T) {
	testcases := []struct {
		name        string
		events      []abcitypes.Event
		expID       int
		expError    bool
		errContains string
	}{
		{
			name: "pass",
			events: []abcitypes.Event{{
				Type: "submit_proposal",
				Attributes: []abcitypes.EventAttribute{
					{Key: "proposal_id", Value: "5"},
					{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
				},
			}},
			expID: 5,
		},
		{
			name: "pass - multiple events",
			events: []abcitypes.Event{
				{
					Type: "message",
					Attributes: []abcitypes.EventAttribute{
						{Key: "action", Value: "/cosmos.gov.v1beta1.MsgSubmitProposal"},
						{Key: "sender", Value: "evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2"},
						{Key: "module", Value: "gov"},
					},
				},
				{
					Type: "submit_proposal",
					Attributes: []abcitypes.EventAttribute{
						{Key: "proposal_id", Value: "5"},
						{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
					},
				},
			},
			expID: 5,
		},
		{
			name: "fail - no submit proposal event",
			events: []abcitypes.Event{{
				Type: "other type",
				Attributes: []abcitypes.EventAttribute{
					{Key: "proposal_id", Value: "4"},
					{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
				},
			}},
			expError:    true,
			errContains: "proposal submission event not found",
		},
		{
			name: "fail - invalid proposal ID",
			events: []abcitypes.Event{{
				Type: "submit_proposal",
				Attributes: []abcitypes.EventAttribute{
					{Key: "proposal_id", Value: "invalid"},
				},
			}},
			expError:    true,
			errContains: "error parsing proposal id",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := getProposalID(tc.events)
			if tc.expError {
				require.Error(t, err, "expected error parsing proposal ID")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error parsing proposal ID")
				require.Equal(t, tc.expID, id, "expected different proposal ID")
			}
		})
	}
}
