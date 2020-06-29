package types

import (
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

const (
	ModuleName   = "bonddoc"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	BondKey = []byte{0x01}
)

func GetBondPrefixKey(did types.Did) []byte {
	return append(BondKey, []byte(did)...)
}
