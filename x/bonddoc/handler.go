package bonddoc

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

type InternalAccountID = string

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateBond:
			return handleMsgCreateBond(ctx, k, msg)
		case MsgUpdateBondStatus:
			return handleMsgUpdateBondStatus(ctx, k, msg)
		default:
			return sdk.ErrUnknownRequest("No match for message type.").Result()
		}
	}
}

func handleMsgCreateBond(ctx sdk.Context, k Keeper, msg MsgCreateBond) sdk.Result {

	if k.BondDocExists(ctx, msg.GetBondDid()) {
		return did.ErrorInvalidDid(types.DefaultCodespace, fmt.Sprintf("Bond doc already exists")).Result()
	}
	k.SetBondDoc(ctx, &msg)

	return sdk.Result{}
}

func handleMsgUpdateBondStatus(ctx sdk.Context, k Keeper, msg MsgUpdateBondStatus) sdk.Result {

	ExistingBondDoc, err := getBondDoc(ctx, k, msg.BondDid)
	if err != nil {
		return sdk.ErrUnknownRequest("Could not find Bond").Result()
	}

	newStatus := msg.Data.Status
	if !newStatus.IsValidProgressionFrom(ExistingBondDoc.GetStatus()) {
		return sdk.ErrUnknownRequest("Invalid Status Progression requested").Result()
	}

	// TODO: actions depending on new status (refer to how projects module does this)

	ExistingBondDoc.SetStatus(newStatus)
	_, _ = k.UpdateBondDoc(ctx, ExistingBondDoc)

	return sdk.Result{}
}

func getBondDoc(ctx sdk.Context, k Keeper, bondDid ixo.Did) (StoredBondDoc, sdk.Error) {
	ixoBondDoc, err := k.GetBondDoc(ctx, bondDid)
	if err != nil {
		return nil, err
	}

	return ixoBondDoc.(StoredBondDoc), nil
}
