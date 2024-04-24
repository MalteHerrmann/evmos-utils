package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/evmos/evmos/v17/app"
	"github.com/evmos/evmos/v17/encoding"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Binary is a struct to hold the necessary information to execute commands
// using a Cosmos SDK-based binary.
type Binary struct {
	// Accounts are the accounts stored in the local keyring.
	Accounts []Account
	// Cdc is the codec to be used for the client.
	Cdc *codec.ProtoCodec

	// Config is the configuration of the binary
	Config BinaryConfig

	// Logger is a logger to be used within all commands.
	Logger zerolog.Logger
}

// BinaryConfig holds the configuration of the binary.
type BinaryConfig struct {
	// Appd is the name of the binary to be executed, e.g. "evmosd".
	Appd string
	// ChainID is the chain ID of the network.
	ChainID string
	// Denom for the fee payments on transactions
	Denom string
	// Home is the home directory of the binary.
	Home string
	// KeyringBackend defines which keyring to use
	KeyringBackend string
	// Node is the endpoint for gRPC connections
	Node string
}

// NewBinary returns a new Binary instance.
func NewBinary(config BinaryConfig) (*Binary, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user home dir")
	}

	// strip the home directory from the given home if already included
	var homeDir string
	if strings.Contains(config.Home, userHome) {
		homeDir = config.Home
	} else {
		homeDir = filepath.Join(userHome, config.Home)
	}

	if _, err = os.Stat(homeDir); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "home directory does not exist: "+homeDir)
	}

	config.Home = homeDir

	// check if binary is installed
	if _, err = exec.LookPath(config.Appd); err != nil {
		return nil, fmt.Errorf("binary %q not installed", config.Appd)
	}

	cdc, ok := GetCodec()
	if !ok {
		return nil, errors.Wrap(err, "failed to get codec")
	}

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	binary := &Binary{
		Cdc:    cdc,
		Config: config,
		Logger: logger,
	}

	if err = binary.getAccounts(); err != nil {
		return nil, err
	}

	return binary, nil
}

// GetCodec returns the codec to be used for the client.
func GetCodec() (*codec.ProtoCodec, bool) {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	protoCodec, ok := encodingConfig.Codec.(*codec.ProtoCodec)

	return protoCodec, ok
}
