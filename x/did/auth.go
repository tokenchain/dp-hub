package did

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/ante"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

func GetPubKeyGetter(keeper Keeper) ante.PubKeyGetter {
	return func(ctx sdk.Context, msg ante.IxoMsg) (pubKey crypto.PubKey, res error) {
		// Get signer PubKey
		var pubKeyEd25519 ed25519tm.PubKeyEd25519
		switch msg := msg.(type) {
		case MsgAddDid:
			copy(pubKeyEd25519[:], base58.Decode(msg.DidDoc.PubKey))
			//pubKeyEd25519 = did.RecoverDidToEd25519PubKey(msg.DidDoc.)

		default:
			// For the remaining messages, the did is the signer
			//fmt.Println("--- GetPubKeyGetter .1")
			fmt.Println(msg.GetSignerDid())

			didDoc, er := keeper.GetDidDoc(ctx, msg.GetSignerDid())
			//fmt.Println("--- GetPubKeyGetter .3")
			if er != nil {
				return nil, er
			}
			//fmt.Println("--- GetPubKeyGetter .4")
			if didDoc == nil {
				return pubKey, exported.Unauthorized("Issuer did not found")
			}

			copy(pubKeyEd25519[:], base58.Decode(didDoc.GetPubKey()))
			//fmt.Println("--- GetPubKeyGetter .5")
		}
		/*
		   MsgAddDid{Did: did:dxp:VrsU9cUAcYgF7f397xtjsX, publicKey: GjKLRmDSCLALj28519q8XwKTmJTfFpobEsWCCKWHhzut}

		   fmt.Println("- json message -")
		   fmt.Println(msg)
		*/
		//fmt.Println("--- GetPubKeyGetter .6")
		//fmt.Println(pubKeyEd25519)
		return pubKeyEd25519, nil
	}
}

func NewDefaultPubKeyGetter(didKeeper Keeper) ante.PubKeyGetter {
	return func(ctx sdk.Context, msg ante.IxoMsg) (pubKey crypto.PubKey, res error) {
		signerDidDoc, err := didKeeper.GetDidDoc(ctx, msg.GetSignerDid())
		if err != nil {
			return pubKey, err
		}
		var pubKeyRaw ed25519tm.PubKeyEd25519
		copy(pubKeyRaw[:], base58.Decode(signerDidDoc.GetPubKey()))
		return pubKeyRaw, nil
	}
}
