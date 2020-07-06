package types

import (
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
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

func GetProjectPrefixKey(did types.Did) []byte {
	return append(ProjectKey, []byte(did)...)
}

func GetAccountPrefixKey(did types.Did) []byte {
	return append(AccountKey, []byte(did)...)
}

func GetWithdrawalPrefixKey(did types.Did) []byte {
	return append(WithdrawalKey, []byte(did)...)
}
