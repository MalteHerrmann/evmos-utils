package main

import (
	"encoding/json"
	"fmt"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

// executeShellCommand executes a shell command and returns the output and error.
func executeShellCommand(command []string, home string, sender string, defaults bool) (string, error) {
	fullCommand := command
	if home != "" {
		fullCommand = append(fullCommand, "--home", home)
	}
	if sender != "" {
		fullCommand = append(fullCommand, "--from", sender)
	}
	if defaults {
		fullCommand = append(fullCommand, defaultFlags...)
	}

	cmd := exec.Command("evmosd", fullCommand...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
	}
	return string(output), err
}

// getCurrentHeight returns the current block height of the node.
func getCurrentHeight() (int, error) {
	output, err := executeShellCommand([]string{"q", "block", "--node", "http://localhost:26657"}, evmosdHome, "", false)
	if err != nil {
		return 0, fmt.Errorf("error executing command: %w", err)
	}

	heightPattern := regexp.MustCompile(`"last_commit":{"height":"(\d+)"`)
	match := heightPattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return 0, fmt.Errorf("did not find block height in output: \n%s", output)
	}

	height, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("error converting height to integer: %w", err)
	}

	return height, nil
}

// getTxEvents returns the transaction events associated with the transaction, whose hash is contained
// in the given output from a transaction command.
//
// It tries to get the transaction hash from the output
// and then waits for the transaction to be included in a block.
// It then returns the transaction events.
func getTxEvents(out string) (txEvents []abcitypes.Event, err error) {
	txHash, err := getTxHashFromResponse(out)
	if err != nil {
		return nil, err
	}

	// Wait for the transaction to be included in a block
	var txOut string
	nAttempts := 10
	for i := 0; i < nAttempts; i++ {
		txOut, err = executeShellCommand([]string{"q", "tx", txHash, "--output=json"}, evmosdHome, "", false)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if txOut == "" {
		return nil, fmt.Errorf("transaction %q not found after %d attempts", txHash, nAttempts)
	}

	return getEventsFromTxResponse(txOut)
}

// getEventsFromTxResponse unpacks the transaction response into the corresponding
// SDK type and returns the events.
func getEventsFromTxResponse(out string) ([]abcitypes.Event, error) {
	var txRes sdk.TxResponse
	err := cdc.UnmarshalJSON([]byte(out), &txRes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling transaction response: %w", err)
	}
	return txRes.Events, nil
}

// TxHashFromResponse is a helper struct for parsing the transaction hash from the response.
type TxHashFromResponse struct {
	TxHash string `json:"txhash"`
}

// getTxHashFromResponse parses the transaction hash from the given response.
func getTxHashFromResponse(out string) (string, error) {
	var txHash TxHashFromResponse
	err := json.Unmarshal([]byte(out), &txHash)
	if err != nil {
		return "", fmt.Errorf("error unpacking transaction hash from json: %w", err)
	}

	return txHash.TxHash, nil
}
