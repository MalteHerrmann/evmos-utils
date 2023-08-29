package utils

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
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

// GetAccounts is a method to retrieve the binaries keys from the configured
// keyring backend and stores it in the Binary struct.
func (bin *Binary) GetAccounts() error {
	out, err := ExecuteBinaryCmd(bin, BinaryCmdArgs{
		Subcommand: []string{"keys", "list", "--output=json"},
	})
	if err != nil {
		return err
	}

	accounts, err := ParseAccountsFromOut(out)
	if err != nil {
		return err
	}

	bin.Accounts = accounts

	return nil
}

// FilterAccountsWithDelegations filters the given list of accounts for those, which are used for staking.
func FilterAccountsWithDelegations(bin *Binary) ([]Account, error) {
	var stakingAccs []Account

	if len(bin.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}

	for _, acc := range bin.Accounts {
		out, err := ExecuteBinaryCmd(bin, BinaryCmdArgs{
			Subcommand: []string{"query", "staking", "delegations", acc.Address, "--output=json"},
		})
		if err != nil {
			return nil, err
		}

		delegations, err := ParseDelegationsFromResponse(bin.Cdc, out)
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
