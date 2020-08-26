package dap

import (
	"github.com/tokenchain/ixo-blockchain/x/dap/auth"
)

var (
	// Auth
	ApproximateFeeForTx              = auth.ApproximateFeeForTx
	SignAndBroadcastTxFromStdSignMsg = auth.SignAndBroadcastTxFromStdSignMsg
	ProcessSig                       = auth.ProcessSig
	SignAndBroadcastTxCli            = auth.SignAndBroadcastTxCli
	SignAndBroadcastTxRest           = auth.SignAndBroadcastTxRest
)
