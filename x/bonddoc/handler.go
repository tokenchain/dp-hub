package bonddoc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

type InternalAccountID = string

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateBond:
			return handleMsgCreateBond(ctx, k, msg)
		case MsgUpdateBondStatus:
			return handleMsgUpdateBondStatus(ctx, k, msg)
		default:
			return nil, x.UnknownRequest("No match for message type.")
		}
	}
}
func handleMsgCreateBond(ctx sdk.Context, k Keeper, msg MsgCreateBond) (*sdk.Result, error) {
	if k.BondDocExists(ctx, msg.GetBondDid()) {
		return nil, x.ErrorBondDocAlreadyExist()
	}
	k.SetBondDoc(ctx, &msg)
	return &sdk.Result{}, nil
}
func handleMsgUpdateBondStatus(ctx sdk.Context, k Keeper, msg MsgUpdateBondStatus) (*sdk.Result, error) {
	ExistingBondDoc, err := getBondDoc(ctx, k, msg.BondDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Bond")
	}
	newStatus := msg.Data.Status
	if !newStatus.IsValidProgressionFrom(ExistingBondDoc.GetStatus()) {
		return nil, x.UnknownRequest("Invalid Status Progression requested")
	}
	// TODO: actions depending on new status (refer to how projects module does this)
	ExistingBondDoc.SetStatus(newStatus)
	_, _ = k.UpdateBondDoc(ctx, ExistingBondDoc)
	return &sdk.Result{}, nil
}
func getBondDoc(ctx sdk.Context, k Keeper, bondDid types.Did) (StoredBondDoc, error) {
	ixoBondDoc, err := k.GetBondDoc(ctx, bondDid)
	if err != nil {
		return nil, err
	}
	return ixoBondDoc.(StoredBondDoc), nil
}
