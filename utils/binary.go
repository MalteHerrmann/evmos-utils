package utils

import (
	"fmt"
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
}

// NewBinary returns a new Binary instance.
func NewBinary(home, appd string) (*Binary, error) {
	// check if home directory exists
	if _, err := os.Stat(home); os.IsNotExist(err) {
		return nil, errors.Wrap(err, fmt.Sprintf("home directory does not exist: %s", home))
	}

	// check if binary is installed
	_, err := exec.LookPath(appd)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("binary %q not installed", appd))
	}

	cdc := GetCodec()

	return &Binary{
		Cdc:  cdc,
		Home: home,
		Appd: appd,
	}, nil
}

// NewEvmosTestingBinary returns a new Binary instance with the default home and appd
// setup for the Evmos local node testing setup.
func NewEvmosTestingBinary() (*Binary, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user home dir")
	}

	defaultEvmosdHome := path.Join(userHome, ".tmp-evmosd")

	return NewBinary(defaultEvmosdHome, "evmosd")
}

// GetCodec returns the codec to be used for the client.
//
//nolint:ireturn // okay to return an interface here
func GetCodec() *codec.ProtoCodec {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	return encodingConfig.Codec.(*codec.ProtoCodec)
}
