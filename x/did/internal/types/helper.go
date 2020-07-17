package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func CastTypeSdkTx(tx sdk.Tx) (IxoTx, bool) {
	mx, ok:= tx.(IxoTx)
	return mx, ok
}
