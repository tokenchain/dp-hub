package test

import (
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

var (
	EmptyDid = ""
)

var ValidDidDoc = types.BaseDidDoc{
	Did:         "FrNMgb6xmPoVfWoFk5zDGn",
	PubKey:      "96UYka2KZEw3nNb58GfP48wPeBUjPrUFrM4AnFhoBzqx",
	Credentials: []exported.DidCredential{},
}
