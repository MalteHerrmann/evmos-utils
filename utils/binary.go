package utils

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/evmos/evmos/v14/app"
	"github.com/evmos/evmos/v14/encoding"
	"github.com/pkg/errors"
)

// Binary is a struct to hold the necessary information to execute commands
// using a Cosmos SDK-based binary.
type Binary struct {
	// Cdc is the codec to be used for the client.
	Cdc *codec.ProtoCodec
	// Home is the home directory of the binary.
	Home string
	// Appd is the name of the binary to be executed, e.g. "evmosd".
	Appd string
	// Accounts are the accounts stored in the local keyring.
	Accounts []Account
	// Logger is a logger to be used within all commands.
	Logger zerolog.Logger
}

// NewBinary returns a new Binary instance.
func NewBinary(home, appd string, logger zerolog.Logger) (*Binary, error) {
	// check if home directory exists
	if _, err := os.Stat(home); os.IsNotExist(err) {
		return nil, errors.Wrap(err, fmt.Sprintf("home directory does not exist: %s", home))
	}

	// check if binary is installed
	_, err := exec.LookPath(appd)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("binary %q not installed", appd))
	}

	cdc, ok := GetCodec()
	if !ok {
		return nil, errors.Wrap(err, "failed to get codec")
	}

	return &Binary{
		Cdc:    cdc,
		Home:   home,
		Appd:   appd,
		Logger: logger,
	}, nil
}

// NewEvmosTestingBinary returns a new Binary instance with the default home and appd
// setup for the Evmos local node testing setup.
func NewEvmosTestingBinary() (*Binary, error) {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	userHome, err := os.UserHomeDir()
	if err != nil {
		return &Binary{Logger: logger}, errors.Wrap(err, "failed to get user home dir")
	}

	defaultEvmosdHome := path.Join(userHome, ".tmp-evmosd")

	return NewBinary(defaultEvmosdHome, "evmosd", logger)
}

// GetCodec returns the codec to be used for the client.
func GetCodec() (*codec.ProtoCodec, bool) {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	protoCodec, ok := encodingConfig.Codec.(*codec.ProtoCodec)

	return protoCodec, ok
}
