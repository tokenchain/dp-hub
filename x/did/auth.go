package did

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"


	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

func GetPubKeyGetter(keeper Keeper) ixo.PubKeyGetter {
	return func(ctx sdk.Context, msg ixo.IxoMsg) ([32]byte, sdk.Result) {

		// Get signer PubKey
		var pubKey [32]byte
		switch msg := msg.(type) {
		case types.MsgAddDid:
			copy(pubKey[:], base58.Decode(msg.DidDoc.PubKey))
		case types.MsgAddCredential:
			did := msg.GetSignerDid()
			didDoc, _ := keeper.GetDidDoc(ctx, did)
			if didDoc == nil {
				return pubKey, sdk.ErrUnauthorized("Issuer did not found").Result()
			}
			copy(pubKey[:], base58.Decode(didDoc.GetPubKey()))
		default:
			return pubKey, sdk.ErrUnknownRequest("No match for message type.").Result()
		}
		return pubKey, sdk.Result{}
	}
}
func NewAnteHandler(didKeeper Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (_ sdk.Context, _ sdk.Result, abort bool) {

		ixoTx, ok := tx.(ixo.IxoTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be ixo.IxoTx").Result(), true
		}

		msg := ixoTx.GetMsgs()[0]
		didMsg := msg.(types.DidMsg)
		pubKey := [32]byte{}

		if didMsg.IsNewDid() {
			addDidMsg := didMsg.(types.MsgAddDid)
			copy(pubKey[:], base58.Decode(addDidMsg.DidDoc.PubKey))
		} else {
			did := ixo.Did(msg.GetSigners()[0])
			didDoc, _ := didKeeper.GetDidDoc(ctx, did)
			if didDoc == nil {
				return ctx,
					sdk.ErrUnauthorized("Issuer did not found").Result(),
					true
			}

			copy(pubKey[:], base58.Decode(didDoc.GetPubKey()))
		}

		var sigs = ixoTx.GetSignatures()
		if len(sigs) != 1 {
			return ctx,
				sdk.ErrUnauthorized("there can only be one signer").Result(),
				true
		}

		res := ixo.VerifySignature(msg, pubKey, sigs[0])

		if !res {
			return ctx, sdk.ErrInternal("Signature Verification failed").Result(), true
		}

		return ctx, sdk.Result{}, false // continue...
	}
}
