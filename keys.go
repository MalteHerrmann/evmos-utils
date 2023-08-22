package main

import (
	"encoding/json"
	"fmt"
)

// Account is the type for a single account.
type Account struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Address     string `json:"address"`
	PubKey      string `json:"pubkey"`
	Delegations []string
}

// getAccounts returns the list of keys from the current running local node
func getAccounts() ([]Account, error) {
	out, err := executeShellCommand([]string{"keys", "list", "--output=json"}, evmosdHome, "", false, false)
	if err != nil {
		return nil, err
	}

	accounts, err := parseAccountsFromOut(out)
	if err != nil {
		return nil, err
	}

	return stakingAccounts(accounts)
}

// stakingAccounts filters the given list of accounts for those, which are used for staking.
func stakingAccounts(accounts []Account) ([]Account, error) {
	var stakingAccs []Account
	for _, acc := range accounts {
		out, err := executeShellCommand([]string{"query", "staking", "delegations", acc.Address, "--output=json"}, evmosdHome, "", false, false)
		if err != nil {
			return nil, err
		}

		delegations, err := parseDelegationsFromResponse(out)
		if err != nil {
			continue
		}

		acc.Delegations = delegations
		stakingAccs = append(stakingAccs, acc)
	}

	return stakingAccs, nil
}

// parseDelegationsFromResponse parses the delegations from the given response.
func parseDelegationsFromResponse(out string) ([]string, error) {
	return nil, nil
}

// parseAccountsFromOut parses the keys from the given output from the keys list command.
func parseAccountsFromOut(out string) ([]Account, error) {
	// Unmarshal the output into a slice of accounts
	var accounts []Account
	err := json.Unmarshal([]byte(out), &accounts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling accounts: %w", err)
	}

	return accounts, nil
}
