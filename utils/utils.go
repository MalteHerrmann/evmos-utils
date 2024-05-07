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
)

// QueryArgs are the arguments passed to a CLI query.
type QueryArgs struct {
	Subcommand []string
	Quiet      bool
}

// ExecuteQueryCmd executes a query command.
func ExecuteQuery(bin *Binary, args QueryArgs) (string, error) {
	queryCommand := args.Subcommand
	queryCommand = append(queryCommand, "--node", bin.Config.Node)

	return ExecuteBinaryCmd(bin, BinaryCmdArgs{
		Subcommand: queryCommand,
		Quiet:      args.Quiet,
	})
}

// TxArgs are the arguments passed to a CLI transaction.
type TxArgs struct {
	Subcommand []string
	From       string
	Quiet      bool
}

// ExecuteTx executes a transaction using the given binary.
func ExecuteTx(bin *Binary, args TxArgs) (string, error) {
	txCommand := args.Subcommand
	txCommand = append(txCommand,
		"--node", bin.Config.Node,
		"--home", bin.Config.Home,
		"--from", args.From,
		"--keyring-backend", bin.Config.KeyringBackend,
		"--gas", "auto",
		"--fees", fmt.Sprintf("%d%s", defaultFees, bin.Config.Denom),
		"--gas-adjustment", "1.3",
		"-b", "sync",
		"-y",
	)

	return ExecuteBinaryCmd(bin, BinaryCmdArgs{
		Subcommand: txCommand,
		Quiet:      args.Quiet,
	})
}

// BinaryCmdArgs are the arguments passed to be executed with the Evmos binary.
type BinaryCmdArgs struct {
	Subcommand []string
	Quiet      bool
}

// ExecuteBinaryCmd executes a shell command and returns the output and error.
func ExecuteBinaryCmd(bin *Binary, args BinaryCmdArgs) (string, error) {
	fullCommand := args.Subcommand

	fmt.Println("Command: ", bin.Config.Appd, strings.Join(fullCommand, " "))
	//#nosec G204 // no risk of injection here because only internal commands are passed
	cmd := exec.Command(bin.Config.Appd, fullCommand...)

	output, err := cmd.CombinedOutput()
	if err != nil && !args.Quiet {
		bin.Logger.Error().Msg(string(output))
	}

	return string(output), err
}

// GetCurrentHeight returns the current block height of the node.
//
// NOTE: Because the response contains uint64 values encoded as strings, this cannot be unmarshalled
// from the BlockResult type. Instead, we use a regex to extract the height from the response.
func GetCurrentHeight(bin *Binary) (int, error) {
	output, err := ExecuteQuery(bin, QueryArgs{
		Subcommand: []string{"q", "block"},
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

	for range nAttempts {
		txOut, err = ExecuteQuery(bin, QueryArgs{
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

// WaitNBlocks waits for the specified amount of blocks being produced
// on the connected network.
func WaitNBlocks(bin *Binary, n int) error {
	currentHeight, err := GetCurrentHeight(bin)
	if err != nil {
		return err
	}

	for {
		bin.Logger.Debug().Msgf("waiting for %d blocks\n", n)
		time.Sleep(2 * time.Second)
		height, err := GetCurrentHeight(bin)
		if err != nil {
			return err
		}

		if height >= currentHeight+n {
			break
		}
	}

	return nil
}
