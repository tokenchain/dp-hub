package types

import (
	"github.com/tokenchain/ixo-blockchain/x/ixo"
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

func GetBondPrefixKey(did ixo.Did) []byte {
	return append(BondKey, []byte(did)...)
}
