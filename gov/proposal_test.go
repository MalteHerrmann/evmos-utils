package gov_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/MalteHerrmann/upgrade-local-node-go/gov"
	"github.com/MalteHerrmann/upgrade-local-node-go/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"
)

//nolint:funlen // function length is okay for tests
func TestGetProposalID(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		events      []sdk.StringEvent
		expID       int
		expError    bool
		errContains string
	}{
		{
			name: "pass",
			events: []sdk.StringEvent{{
				Type: "submit_proposal",
				Attributes: []sdk.Attribute{
					{Key: "proposal_id", Value: "5"},
					{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
				},
			}},
			expID: 5,
		},
		{
			name: "pass - multiple events",
			events: []sdk.StringEvent{
				{
					Type: "message",
					Attributes: []sdk.Attribute{
						{Key: "action", Value: "/cosmos.gov.v1beta1.MsgSubmitProposal"},
						{Key: "sender", Value: "evmos1vv6hqcxp0w5we5rzdvf4ddhsas5gx0dep8vmv2"},
						{Key: "module", Value: "gov"},
					},
				},
				{
					Type: "submit_proposal",
					Attributes: []sdk.Attribute{
						{Key: "proposal_id", Value: "5"},
						{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
					},
				},
			},
			expID: 5,
		},
		{
			name: "fail - no submit proposal event",
			events: []sdk.StringEvent{{
				Type: "other type",
				Attributes: []sdk.Attribute{
					{Key: "proposal_id", Value: "4"},
					{Key: "proposal_messages", Value: ",/cosmos.gov.v1.MsgExecLegacyContent"},
				},
			}},
			expError:    true,
			errContains: "proposal submission event not found",
		},
		{
			name: "fail - invalid proposal ID",
			events: []sdk.StringEvent{{
				Type: "submit_proposal",
				Attributes: []sdk.Attribute{
					{Key: "proposal_id", Value: "invalid"},
				},
			}},
			expError:    true,
			errContains: "error parsing proposal id",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			propID, err := gov.GetProposalIDFromSubmitEvents(tc.events)
			if tc.expError {
				require.Error(t, err, "expected error parsing proposal ID")
				require.ErrorContains(t, err, tc.errContains, "expected different error")
			} else {
				require.NoError(t, err, "unexpected error parsing proposal ID")
				require.Equal(t, tc.expID, propID, "expected different proposal ID")
			}
		})
	}
}

// This function tests if the custom type ProposalIDsFromProposalsQueryResponse
// is compatible with the output of the proposals query.
//
//nolint:paralleltest // only one test included
func TestProposalIDsFromProposalsQueryResponse(t *testing.T) {
	cdc, ok := utils.GetCodec()
	require.True(t, ok, "failed to get codec")

	sdkProposals := govv1types.Proposals{
		&govv1types.Proposal{Id: 5},
		&govv1types.Proposal{Id: 6},
	}

	sdkRes := govv1types.QueryProposalsResponse{
		Proposals: sdkProposals,
	}

	expParsedProposals := make([]gov.Proposal, 0, len(sdkProposals))
	for _, proposal := range sdkProposals {
		expParsedProposals = append(expParsedProposals, gov.Proposal{
			ProposalID: strconv.FormatUint(proposal.Id, 10),
		})
	}

	bz, err := cdc.MarshalJSON(&sdkRes)
	require.NoError(t, err, "unexpected error marshaling proposals")

	var res gov.ProposalIDsFromProposalsQueryResponse
	err = json.Unmarshal(bz, &res)
	require.NoError(t, err, "unexpected error parsing proposal IDs")
	require.Equal(t, expParsedProposals, res.Proposals, "expected different parsed proposals")
}
