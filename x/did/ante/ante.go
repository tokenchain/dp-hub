package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
)
//use in bond
func DefaultAnteHandler(ak auth.AccountKeeper, bk bank.Keeper, sk supply.Keeper, pubKeyGetter PubKeyGetter) sdk.AnteHandler {
	//return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		NewDapConsumeGasForTxSizeDecorator(ak, pubKeyGetter),
		NewDapPubKeyDecorator(ak, pubKeyGetter), // SetPubKeyDecorator must be called before all signature verification decorators
		//ante.NewValidateSigCountDecorator(ak),
		NewDeductFeeDecorator(ak, sk, pubKeyGetter),
		NewConsumeVerSignGasDecorator(ak, pubKeyGetter),
		NewSigVerificationAndIncrementSequenceDecorator(ak, pubKeyGetter),
	)
}


func DidAnteHandler(ak auth.AccountKeeper, bk bank.Keeper, sk supply.Keeper, pubKeyGetter PubKeyGetter) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		NewDapConsumeGasForTxSizeDecorator(ak, pubKeyGetter),
		NewDapPubKeyDecorator(ak, pubKeyGetter), // SetPubKeyDecorator must be called before all signature verification decorators
		NewDeductFeeDecorator(ak, sk, pubKeyGetter),
		NewConsumeVerSignGasDecorator(ak, pubKeyGetter),
		NewSigVerificationAndIncrementSequenceDecorator(ak, pubKeyGetter),
	)
}
