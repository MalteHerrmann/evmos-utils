package testutil

import "github.com/MalteHerrmann/upgrade-local-node-go/utils"

// NewEvmosdBinaryWithCodec returns a new Binary with the evmosd binary and the evmos codec.
func NewEvmosdBinaryWithCodec() *utils.Binary {
	cdc := utils.GetCodec()
	return &utils.Binary{Cdc: cdc}
}
