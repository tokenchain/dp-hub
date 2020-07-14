package dap

import (
	"github.com/tokenchain/ixo-blockchain/x/dap/auth"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
)

const (
	IxoNativeToken = types.NativeToken
)

var (
	// Auth

	ApproximateFeeForTx     = auth.ApproximateFeeForTx
	GenerateOrBroadcastMsgs = auth.GenerateOrBroadcastMsgs
	//CompleteAndBroadcastTxRest       = auth.CompleteAndBroadcastTxRest
	SignAndBroadcastTxFromStdSignMsg = auth.SignAndBroadcastTxFromStdSignMsg

	ProcessSig             = auth.ProcessSig
	SignAndBroadcastTxCli  = auth.SignAndBroadcastTxCli
	SignAndBroadcastTxRest = auth.SignAndBroadcastTxRest

	// Types
	IxoDecimals = types.IxoDecimals
)
