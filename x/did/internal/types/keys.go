package types

import "github.com/tokenchain/ixo-blockchain/x/did/exported"

const (
	ModuleName   = "did"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var DidKey = []byte{0x01}

func GetDidPrefixKey(did exported.Did) []byte {
	return append(DidKey, []byte(did)...)
}
