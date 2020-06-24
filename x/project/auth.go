package project

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/project/internal/types"
)

func GetPubKeyGetter(keeper Keeper, didKeeper did.Keeper) ixo.PubKeyGetter {
	return func(ctx sdk.Context, msg ixo.IxoMsg) ([32]byte, sdk.Result) {
		// Message must be a ProjectMsg
		projectMsg := msg.(types.ProjectMsg)

		// Get signer PubKey
		var pubKey [32]byte
		if projectMsg.IsNewDid() {
			createProjectMsg := msg.(types.MsgCreateProject)
			copy(pubKey[:], base58.Decode(createProjectMsg.GetPubKey()))
		} else {
			if projectMsg.IsWithdrawal() {
				signerDid := msg.GetSignerDid()
				didDoc, _ := didKeeper.GetDidDoc(ctx, signerDid)
				if didDoc == nil {
					return pubKey, sdk.ErrUnauthorized("signer did not found").Result()
				}
				copy(pubKey[:], base58.Decode(didDoc.GetPubKey()))
			} else {
				projectDid := msg.GetSignerDid()
				projectDoc, err := keeper.GetProjectDoc(ctx, projectDid)
				if err != nil {
					return pubKey, sdk.ErrInternal("project did not found").Result()
				}
				copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
			}
		}
		return pubKey, sdk.Result{}
	}
}

func NewAnteHandler(projectKeeper Keeper, didKeeper did.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (_ sdk.Context, _ sdk.Result, abort bool) {

		ixoTx, ok := tx.(ixo.IxoTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be ixo.IxoTx").Result(), true
		}

		msg := ixoTx.GetMsgs()[0]
		projectMsg := msg.(types.ProjectMsg)
		pubKey := [32]byte{}

		if projectMsg.IsNewDid() {
			createProjectMsg := msg.(types.MsgCreateProject)
			copy(pubKey[:], base58.Decode(createProjectMsg.GetPubKey()))

		} else {
			if projectMsg.IsWithdrawal() {
				signerDid := ixo.Did(msg.GetSigners()[0])
				didDoc, _ := didKeeper.GetDidDoc(ctx, signerDid)
				if didDoc == nil {
					return ctx,
						sdk.ErrUnauthorized("Issuer did not found").Result(),
						true
				}

				copy(pubKey[:], base58.Decode(didDoc.GetPubKey()))
			} else {
				projectDid := ixo.Did(msg.GetSigners()[0])
				projectDoc, err := projectKeeper.GetProjectDoc(ctx, projectDid)
				if err != nil {
					return ctx, sdk.ErrInternal("project did not found").Result(), false
				}

				copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
			}
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
