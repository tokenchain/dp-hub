package ixo

const (
	ModuleName        = "dapx"
	DefaultParamspace = ModuleName
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName

	FeeRemainderPool = "fee_remainder_pool"

	FeeIdPrefix          = "fee:"
	FeeContractIdPrefix  = FeeIdPrefix + "contract:"
	SubscriptionIdPrefix = FeeIdPrefix + "subscription:"
)