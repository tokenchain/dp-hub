package types

import "github.com/tokenchain/ixo-blockchain/x/did/exported"

var (
	EmptyDid = ""
)

var ValidDidDoc = BaseDidDoc{
	Did:         "FrNMgb6xmPoVfWoFk5zDGn",
	PubKey:      "96UYka2KZEw3nNb58GfP48wPeBUjPrUFrM4AnFhoBzqx",
	Credentials: []exported.DidCredential{},
}
