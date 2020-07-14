package did

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/ante"
)

func GetPubKeyGetter(keeper Keeper) ante.PubKeyGetter {
	return func(ctx sdk.Context, msg ante.IxoMsg) (pubKey crypto.PubKey, res error) {
		// Get signer PubKey
		var pubKeyEd25519 ed25519.PubKeyEd25519
		switch msg := msg.(type) {
		case MsgAddDid:
			copy(pubKeyEd25519[:], base58.Decode(msg.DidDoc.PubKey))
			//pubKeyEd25519 = did.RecoverDidToEd25519PubKey(msg.DidDoc.)
		default:
			// For the remaining messages, the did is the signer
			didDoc, _ := keeper.GetDidDoc(ctx, msg.GetSignerDid())
			if didDoc == nil {
				return pubKey, Unauthorized("Issuer did not found")
			}
			copy(pubKeyEd25519[:], base58.Decode(didDoc.GetPubKey()))
		}
		/*
		   MsgAddDid{Did: did:dxp:VrsU9cUAcYgF7f397xtjsX, publicKey: GjKLRmDSCLALj28519q8XwKTmJTfFpobEsWCCKWHhzut}

		   fmt.Println("- json message -")
		   fmt.Println(msg)
		*/
		return pubKeyEd25519, nil
	}
}

func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}
