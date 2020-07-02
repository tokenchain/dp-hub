package ixo

import "github.com/tokenchain/ixo-blockchain/x/ixo/types"

const (
	ModuleName        = "dapx"
	NativeToken       = types.IxoNativeToken
	DefaultParamspace = ModuleName
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName

	FeeRemainderPool = "fee_remainder_pool"

	FeeIdPrefix          = "fee:"
	FeeContractIdPrefix  = FeeIdPrefix + "contract:"
	SubscriptionIdPrefix = FeeIdPrefix + "subscription:"
)
