package dap

import (
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
)

const (
	IxoNativeToken = types.NativeToken
)

type (
	IxoTx        = types.IxoTx
	IxoSignature = types.IxoSignature
	IxoMsg       = types.IxoMsg
	PubKeyGetter = types.PubKeyGetter
)

var (
	// Auth
	NewDefaultPubKeyGetter           = types.NewDefaultPubKeyGetter
	NewDefaultAnteHandler            = types.NewDefaultAnteHandler
	ApproximateFeeForTx              = types.ApproximateFeeForTx
	GenerateOrBroadcastMsgs          = types.GenerateOrBroadcastMsgs
	CompleteAndBroadcastTxRest       = types.CompleteAndBroadcastTxRest
	SignAndBroadcastTxFromStdSignMsg = types.SignAndBroadcastTxFromStdSignMsg
	NewSignature                     = types.NewSignature
	ProcessSig                       = types.ProcessSig
	SignAndBroadcastTxCli            = types.SignAndBroadcastTxCli
	SignAndBroadcastTxRest           = types.SignAndBroadcastTxRest
	ProcessPubKey                    = types.ProcessPubKey

	// Types
	IxoDecimals = types.IxoDecimals

	// Tx
	DefaultTxDecoder  = types.DefaultTxDecoder
	NewIxoTxSingleMsg = types.NewIxoTxSingleMsg

	//helper
	DidToAddr = types.DidToAddr
)
