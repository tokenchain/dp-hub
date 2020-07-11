package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	IxoDecimals = sdk.NewDec(1000)
)

const (
	NativeToken         = "dap"
	Ed25519SignatureLen = 64
)
