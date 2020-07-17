package ante

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	anteDid "github.com/tokenchain/ixo-blockchain/x/did/ante"
	didexported "github.com/tokenchain/ixo-blockchain/x/did/exported"
)

type DeductFeeDecorator struct {
	supplyKeeper types.SupplyKeeper
	anteDid.SigVerification
}

func NewDeductFeeDecorator(ak keeper.AccountKeeper, sk types.SupplyKeeper, p didexported.DidKeeper) DeductFeeDecorator {
	return DeductFeeDecorator{
		SigVerification: anteDid.NewSigVerification(ak, p),
		supplyKeeper:    sk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sv, _, e := dfd.RetrievePubkey(ctx, tx, simulate)
	if e != nil {
		return ctx, anteDid.InvalidTxDecodePubkeyNotFound(e)
	}

	if addr := dfd.supplyKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
	}

	// deduct the fees
	if !sv.GetDapTx().GetFee().IsZero() {
		if err = anteDid.DeductFees(dfd.supplyKeeper, ctx, dfd.GetSignerAccount(ctx), sv.GetDapTx().GetFee()); err != nil {
			return ctx, err
		}
		// reload the account as fees have been deducted
		if err = dfd.RefreshAccount(ctx); err != nil {
			return ctx, err
		}
	}
	fmt.Println("✅  fee deduction pass ....")
	return next(ctx, tx, simulate)
}
