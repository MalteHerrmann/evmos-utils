package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// defaultEvmosdHome is the home directory for the Evmos binary.
var defaultEvmosdHome string

func init() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultEvmosdHome = path.Join(userHome, ".tmp-evmosd")
}

// BinaryCmdArgs are the arguments passed to be executed with the Evmos binary.
type BinaryCmdArgs struct {
	Subcommand  []string
	Home        string
	From        string
	UseDefaults bool
	Quiet       bool
}

// ExecuteBinaryCmd executes a shell command and returns the output and error.
func ExecuteBinaryCmd(args BinaryCmdArgs) (string, error) {
	fullCommand := args.Subcommand
	if args.Home == "" {
		fullCommand = append(fullCommand, "--home", defaultEvmosdHome)
	} else {
		fullCommand = append(fullCommand, "--home", args.Home)
	}
	if args.From != "" {
		fullCommand = append(fullCommand, "--from", args.From)
	}
	if args.UseDefaults {
		fullCommand = append(fullCommand, defaultFlags...)
	}

	cmd := exec.Command("evmosd", fullCommand...)
	output, err := cmd.CombinedOutput()
	if err != nil && !args.Quiet {
		fmt.Println(string(output))
	}
	return string(output), err
}

// GetCurrentHeight returns the current block height of the node.
func GetCurrentHeight() (int, error) {
	output, err := ExecuteBinaryCmd(BinaryCmdArgs{
		Subcommand: []string{"q", "block", "--node", "http://localhost:26657"},
	})
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

// GetTxEvents returns the transaction events associated with the transaction, whose hash is contained
// in the given output from a transaction command.
//
// It tries to get the transaction hash from the output
// and then waits for the transaction to be included in a block.
// It then returns the transaction events.
func GetTxEvents(out string) (txEvents []abcitypes.Event, err error) {
	txHash, err := getTxHashFromResponse(out)
	if err != nil {
		return nil, err
	}

	// Wait for the transaction to be included in a block
	var txOut string
	nAttempts := 10
	for i := 0; i < nAttempts; i++ {
		txOut, err = ExecuteBinaryCmd(BinaryCmdArgs{
			Subcommand: []string{"q", "tx", txHash, "--output=json"},
			Quiet:      true,
		})

		if err == nil {
			break
		}
		if !strings.Contains(txOut, fmt.Sprintf("tx (%s) not found", txHash)) {
			return nil, fmt.Errorf("unexpected error while querying transaction %s: %w", txHash, err)
		}
		time.Sleep(2 * time.Second)
	}

	if strings.Contains(txOut, fmt.Sprintf("tx (%s) not found", txHash)) {
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
		return nil, fmt.Errorf("error unmarshalling transaction response: %w\n\nresponse: %s", err, out)
	}
	return txRes.Events, nil
}

// getTxHashFromResponse parses the transaction hash from the given response.
func getTxHashFromResponse(out string) (string, error) {
	var txHash sdk.TxResponse
	err := cdc.UnmarshalJSON([]byte(out), &txHash)
	if err != nil {
		return "", fmt.Errorf("error unpacking transaction hash from json: %w", err)
	}

	if txHash.Code != 0 {
		return "", fmt.Errorf("transaction failed with code %d", txHash.Code)
	}
	return txHash.TxHash, nil
}
