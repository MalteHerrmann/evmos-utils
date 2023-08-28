package utils

import (
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
	// cdc is the codec to be used for the client.
	cdc = encodingConfig.Codec
	// encodingConfig specifies the encoding configuration to be used for the client
	encodingConfig = encoding.MakeConfig(app.ModuleBasics)
)
