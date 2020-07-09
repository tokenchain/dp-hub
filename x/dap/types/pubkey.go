package types

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

// DidKeeper defines the did contract that must be fulfilled throughout the ixo module
type DidKeeper interface {
	GetDidDoc(ctx sdk.Context, did exported.Did) (exported.DidDoc, error)
	SetDidDoc(ctx sdk.Context, did exported.DidDoc) (err error)
	AddDidDoc(ctx sdk.Context, did exported.DidDoc)
	AddCredentials(ctx sdk.Context, did exported.Did, credential exported.DidCredential) (err error)
	GetAllDidDocs(ctx sdk.Context) (didDocs []exported.DidDoc)
	GetAddDids(ctx sdk.Context) (dids []exported.Did)
}

func NewDefaultPubKeyGetter(didKeeper DidKeeper) PubKeyGetter {
	return func(ctx sdk.Context, msg IxoMsg) (pubKey crypto.PubKey, res error) {
		signerDidDoc, err := didKeeper.GetDidDoc(ctx, msg.GetSignerDid())
		if err != nil {
			return pubKey, err
		}
		var pubKeyRaw ed25519tm.PubKeyEd25519
		copy(pubKeyRaw[:], base58.Decode(signerDidDoc.GetPubKey()))
		return pubKeyRaw, nil
	}
}
