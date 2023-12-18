package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmosutils "github.com/evmos/evmos/v14/utils"
)

// BinaryCmdArgs are the arguments passed to be executed with the Evmos binary.
type BinaryCmdArgs struct {
	Subcommand  []string
	Home        string
	From        string
	UseDefaults bool
	Quiet       bool
}

// ExecuteBinaryCmd executes a shell command and returns the output and error.
func ExecuteBinaryCmd(bin *Binary, args BinaryCmdArgs) (string, error) {
	fullCommand := args.Subcommand
	if args.Home == "" {
		fullCommand = append(fullCommand, "--home", bin.Home)
	} else {
		fullCommand = append(fullCommand, "--home", args.Home)
	}

	if args.From != "" {
		fullCommand = append(fullCommand, "--from", args.From)
	}

	if args.UseDefaults {
		defaultFlags := getDefaultFlags()
		fullCommand = append(fullCommand, defaultFlags...)
	}

	//#nosec G204 // no risk of injection here because only internal commands are passed
	cmd := exec.Command(bin.Appd, fullCommand...)

	output, err := cmd.CombinedOutput()
	if err != nil && !args.Quiet {
		bin.Logger.Error().Msg(string(output))
	}

	return string(output), err
}

// getDefaultFlags returns the default flags to be used for the Evmos binary.
func getDefaultFlags() []string {
	chainID := evmosutils.TestnetChainID + "-1"

	defaultFlags := []string{
		"--chain-id", chainID,
		"--keyring-backend", "test",
		"--gas", "auto",
		"--fees", fmt.Sprintf("%d%s", defaultFees, denom),
		"--gas-adjustment", "1.3",
		"-b", "sync",
		"-y",
	}

	return defaultFlags
}

// GetCurrentHeight returns the current block height of the node.
func GetCurrentHeight(bin *Binary) (int, error) {
	output, err := ExecuteBinaryCmd(bin, BinaryCmdArgs{
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
func GetTxEvents(bin *Binary, out string) ([]sdk.StringEvent, error) {
	txHash, err := GetTxHashFromTxResponse(bin.Cdc, out)
	if err != nil {
		return nil, err
	}

	// Wait for the transaction to be included in a block
	var (
		txOut     string
		nAttempts = 10
	)

	for i := 0; i < nAttempts; i++ {
		txOut, err = ExecuteBinaryCmd(bin, BinaryCmdArgs{
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

	return GetEventsFromTxResponse(bin.Cdc, txOut)
}

// GetEventsFromTxResponse unpacks the transaction response into the corresponding
// SDK type and returns the events.
func GetEventsFromTxResponse(cdc *codec.ProtoCodec, out string) ([]sdk.StringEvent, error) {
	var txRes sdk.TxResponse

	err := cdc.UnmarshalJSON([]byte(out), &txRes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling transaction response: %w\n\nresponse: %s", err, out)
	}

	logs := txRes.Logs
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs found in transaction response: %s", out)
	}

	var events []sdk.StringEvent

	for _, msgLog := range logs {
		for _, event := range msgLog.Events {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetTxHashFromTxResponse parses the transaction hash from the given response.
func GetTxHashFromTxResponse(cdc *codec.ProtoCodec, out string) (string, error) {
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

// Wait waits for the specified amount of seconds.
func Wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
