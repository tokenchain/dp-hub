package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

func DefaultAnteHandler(ak auth.AccountKeeper, bk bank.Keeper, sk supply.Keeper, dk exported.DidKeeper) sdk.AnteHandler {
	//return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		NewDapConsumeGasForTxSizeDecorator(ak, dk),
		NewDapPubKeyDecorator(ak, dk), // SetPubKeyDecorator must be called before all signature verification decorators
		//ante.NewValidateSigCountDecorator(ak),
		NewDeductFeeDecorator(ak, sk, dk),
		NewConsumeVerSignGasDecorator(ak, dk),
		//ante.NewDeductFeeDecorator(ak, bk),
		//ante.NewSigGasConsumeDecorator(ak, sign),
		NewSigVerificationDecorator(ak, dk),
		//ante.NewIncrementSequenceDecorator(ak),
	)
}


func DidAnteHandler(ak auth.AccountKeeper, bk bank.Keeper, sk supply.Keeper, dk exported.DidKeeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		NewDapConsumeGasForTxSizeDecorator(ak, dk),
		NewDapPubKeyDecorator(ak, dk), // SetPubKeyDecorator must be called before all signature verification decorators
		NewDeductFeeDecorator(ak, sk, dk),
		NewConsumeVerSignGasDecorator(ak, dk),
		NewSigVerificationDecorator(ak, dk),
	)
}
