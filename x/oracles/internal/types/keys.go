package types

import (
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

const (
	ModuleName   = "oracles"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	OracleKey = []byte{0x00}
)

func GetOraclePrefixKey(did exported.Did) []byte {
	return append(OracleKey, []byte(did)...)
}
