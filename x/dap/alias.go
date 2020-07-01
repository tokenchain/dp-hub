package dap

import "github.com/tokenchain/ixo-blockchain/x/dap/internal/types"

const (
	DarkpoolNativeToken = types.DarkpoolNativeToken
)

type (
	PubKeyGetter = types.PubKeyGetter
	IxoTx        = types.IxoTx
	IxoSignature = types.IxoSignature
	IxoMsg       = types.IxoMsg
)

var (
	// Auth
	NewDefaultPubKeyGetter           = types.NewDefaultPubKeyGetter
	ProcessSig                       = types.ProcessSig
	NewDefaultAnteHandler            = types.NewDefaultAnteHandler
	ApproximateFeeForTx              = types.ApproximateFeeForTx
	GenerateOrBroadcastMsgs          = types.GenerateOrBroadcastMsgs
	CompleteAndBroadcastTxRest       = types.CompleteAndBroadcastTxRest
	SignAndBroadcastTxFromStdSignMsg = types.SignAndBroadcastTxFromStdSignMsg

	// Types
	IxoDecimals = types.IxoDecimals

	// Tx
	DefaultTxDecoder  = types.DefaultTxDecoder
	NewIxoTxSingleMsg = types.NewIxoTxSingleMsg
)
