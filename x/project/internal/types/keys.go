package types

import (
	"github.com/tokenchain/dp-hub/x/did/exported"
)

const (
	ModuleName        = "project"
	DefaultParamspace = ModuleName
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
)

var (
	ProjectKey    = []byte{0x01}
	AccountKey    = []byte{0x02}
	WithdrawalKey = []byte{0x03}
)

func GetProjectPrefixKey(did exported.Did) []byte {
	return append(ProjectKey, []byte(did)...)
}

func GetAccountPrefixKey(did exported.Did) []byte {
	return append(AccountKey, []byte(did)...)
}

func GetWithdrawalPrefixKey(did exported.Did) []byte {
	return append(WithdrawalKey, []byte(did)...)
}
