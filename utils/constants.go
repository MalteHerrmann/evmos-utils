package utils

import (
	evmosutils "github.com/evmos/evmos/v14/utils"
)

const (
	// defaultFees is the amount of fees to be sent with a default transaction.
	defaultFees int = 1e18 // 1 aevmos
	// DeltaHeight is the amount of blocks in the future that the upgrade will be scheduled.
	DeltaHeight = 20
	// denom is the denomination used for the local node.
	denom = evmosutils.BaseDenom
)
