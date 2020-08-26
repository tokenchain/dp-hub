package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	StakingToken = sdk.DefaultBondDenom
	NativeToken  = "mdap"
	TestToken    = "xxb"
)

var (
	NativeDecimals = sdk.NewDec(1000)
)
