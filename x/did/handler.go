package did

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/dp-block/x/did/exported"

	"github.com/tokenchain/dp-block/x/did/internal/keeper"
	"github.com/tokenchain/dp-block/x/did/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgAddDid:
			return handleMsgAddDidDoc(ctx, k, msg)
		case types.MsgAddCredential:
			return handleMsgAddCredential(ctx, k, msg)
		default:
			return nil, exported.UnknownRequest("No match for message type.")
		}
	}
}

func handleMsgAddDidDoc(ctx sdk.Context, k keeper.Keeper, msg types.MsgAddDid) (*sdk.Result, error) {
	newDidDoc := msg.DidDoc

	if len(newDidDoc.Credentials) > 0 {
		return nil, exported.UnknownRequest("Cannot add a new DID with existing Credentials")
	}

	err := k.SetDidDoc(ctx, newDidDoc)
	if err != nil {
		return nil, err
	}
	fmt.Println("handleMsgAddDidDoc complete")
	return &sdk.Result{}, nil
}

func handleMsgAddCredential(ctx sdk.Context, k keeper.Keeper, msg types.MsgAddCredential) (*sdk.Result, error) {
	err := k.AddCredentials(ctx, msg.DidCredential.Claim.Id, msg.DidCredential)
	if err != nil {
		return nil, err
	}
	fmt.Println("handleMsgAddCredential complete")
	return &sdk.Result{}, nil
}
