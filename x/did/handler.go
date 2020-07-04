package did

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgAddDid:
			return handleMsgAddDidDoc(ctx, k, msg)
		case MsgAddCredential:
			return handleMsgAddCredential(ctx, k, msg)
		default:
			return nil, x.UnknownRequest("No match for message type.")
		}
	}
}

func handleMsgAddDidDoc(ctx sdk.Context, k Keeper, msg MsgAddDid) (*sdk.Result, error) {
	newDidDoc := msg.DidDoc

	if len(newDidDoc.Credentials) > 0 {
		return nil, x.UnknownRequest("Cannot add a new DID with existing Credentials")
	}

	err := k.SetDidDoc(ctx, newDidDoc)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleMsgAddCredential(ctx sdk.Context, k Keeper, msg MsgAddCredential) (*sdk.Result, error) {
	err := k.AddCredentials(ctx, msg.DidCredential.Claim.Id, msg.DidCredential)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}
