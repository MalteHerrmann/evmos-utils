package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ParseDelegationsFromResponse parses the delegations from the given response.
func ParseDelegationsFromResponse(cdc *codec.ProtoCodec, out string) ([]stakingtypes.Delegation, error) {
	var res stakingtypes.QueryDelegatorDelegationsResponse

	err := cdc.UnmarshalJSON([]byte(out), &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling delegations: %w", err)
	}

	delegations := make([]stakingtypes.Delegation, len(res.DelegationResponses))
	for i, delegation := range res.DelegationResponses {
		delegations[i] = delegation.Delegation
	}

	return delegations, nil
}
