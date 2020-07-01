package types

import (
	"github.com/tokenchain/dp-hub/x/ixo/types"
)

const (
	ModuleName   = "did"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var DidKey = []byte{0x01}

func GetDidPrefixKey(did types.Did) []byte {
	return append(DidKey, []byte(did)...)
}
