package types

import "regexp"

const (
	ModuleName   = "did"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	ValidDid   = regexp.MustCompile(`^did:(dxp:|sov:)([a-zA-Z0-9]){21,22}([/][a-zA-Z0-9:]+|)$`)
	IsValidDid = ValidDid.MatchString
	DidKey     = []byte{0x01}
)

func GetDidPrefixKey(did Did) []byte {
	return append(DidKey, []byte(did)...)
}
