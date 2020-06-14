package treasury

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/treasury/internal/types"
)

func NewAnteHandler(didKeeper did.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (_ sdk.Context, _ sdk.Result, abort bool) {

		ixoTx, ok := tx.(ixo.IxoTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be ixo.IxoTx").Result(), true
		}

		msg := ixoTx.GetMsgs()[0]
		treasuryMsg := msg.(types.TreasuryMessage)
		pubKey := [32]byte{}
		copy(pubKey[:], base58.Decode(treasuryMsg.GetPubKey()))

		// Check that sender's DID is ledgered
		senderDidDoc, _ := didKeeper.GetDidDoc(ctx, treasuryMsg.GetSenderDid())
		if senderDidDoc == nil {
			return ctx,
				sdk.ErrUnauthorized("Sender did not found").Result(),
				true
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
