package utils

import (
	"encoding/json"
	"fmt"

	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Account is the type for a single account.
type Account struct {
	Name        string                    `json:"name"`
	Type        string                    `json:"type"`
	Address     string                    `json:"address"`
	PubKey      string                    `json:"pubkey"`
	Delegations []stakingtypes.Delegation `json:"delegations"`
}

// GetAccounts returns the list of keys from the current running local node.
func GetAccounts(bin *Binary) ([]Account, error) {
	out, err := ExecuteBinaryCmd(bin, BinaryCmdArgs{
		Subcommand: []string{"keys", "list", "--output=json"},
	})
	if err != nil {
		return nil, err
	}

	accounts, err := ParseAccountsFromOut(out)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// FilterAccountsWithDelegations filters the given list of accounts for those, which are used for staking.
func FilterAccountsWithDelegations(bin *Binary, accounts []Account) ([]Account, error) {
	var stakingAccs []Account

	for _, acc := range accounts {
		out, err := ExecuteBinaryCmd(bin, BinaryCmdArgs{
			Subcommand: []string{"query", "staking", "delegations", acc.Address, "--output=json"},
		})
		if err != nil {
			return nil, err
		}

		delegations, err := ParseDelegationsFromResponse(bin, out)
		if err != nil {
			continue
		}

		acc.Delegations = delegations
		if len(delegations) > 0 {
			stakingAccs = append(stakingAccs, acc)
		}
	}

	return stakingAccs, nil
}

// ParseDelegationsFromResponse parses the delegations from the given response.
func ParseDelegationsFromResponse(bin *Binary, out string) ([]stakingtypes.Delegation, error) {
	var res stakingtypes.QueryDelegatorDelegationsResponse

	err := bin.cdc.UnmarshalJSON([]byte(out), &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling delegations: %w", err)
	}

	delegations := make([]stakingtypes.Delegation, len(res.DelegationResponses))
	for i, delegation := range res.DelegationResponses {
		delegations[i] = delegation.Delegation
	}

	return delegations, nil
}

// ParseAccountsFromOut parses the keys from the given output from the keys list command.
func ParseAccountsFromOut(out string) ([]Account, error) {
	var (
		accounts = make([]Account, 0)
		keys     = make([]cryptokeyring.KeyOutput, 0)
	)

	err := json.Unmarshal([]byte(out), &keys)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling keys: %w", err)
	}

	for _, key := range keys {
		accounts = append(accounts, Account{
			Name:    key.Name,
			Type:    key.Type,
			Address: key.Address,
			PubKey:  key.PubKey,
		})
	}

	return accounts, nil
}
