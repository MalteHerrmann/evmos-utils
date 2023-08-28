package utils

import (
	"fmt"

	"github.com/evmos/evmos/v14/app"
	"github.com/evmos/evmos/v14/encoding"
	evmosutils "github.com/evmos/evmos/v14/utils"
)

const (
	// The amount of fees to be sent with a default transaction.
	defaultFees int = 1e18 // 1 aevmos
	// The denomination used for the local node.
	denom = evmosutils.BaseDenom
)

var (
	// cdc is the codec to be used for the client
	cdc = encodingConfig.Codec
	// encodingConfig specifies the encoding configuration to be used for the client
	encodingConfig = encoding.MakeConfig(app.ModuleBasics)

	// The chain ID of the node that will be upgraded.
	chainID = evmosutils.TestnetChainID + "-1"
	// defaultFlags are the default flags to be used for the client.
	defaultFlags = []string{
		"--chain-id", chainID,
		"--keyring-backend", "test",
		"--gas", "auto",
		"--fees", fmt.Sprintf("%d%s", defaultFees, denom),
		"--gas-adjustment", "1.3",
		"-b", "sync",
		"-y",
	}
)
