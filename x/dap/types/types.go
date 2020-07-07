package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	IxoDecimals = sdk.NewDec(1000)
)

const (
	NativeToken         = "dap"
	ed25519SignatureLen = 64
)
