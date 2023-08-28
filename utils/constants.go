package utils

import (
	evmosutils "github.com/evmos/evmos/v14/utils"
)

const (
	// The amount of fees to be sent with a default transaction.
	defaultFees int = 1e18 // 1 aevmos
	// The denomination used for the local node.
	denom = evmosutils.BaseDenom
)
