package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/supply"
	didAnte "github.com/tokenchain/ixo-blockchain/x/did/ante"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

func ProjectAnteHandle(ak auth.AccountKeeper, sk supply.Keeper, dk exported.DidKeeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		didAnte.NewDapConsumeGasForTxSizeDecorator(ak, dk),
		NewDeductFeeDecorator(ak, sk, dk),
		didAnte.NewConsumeVerSignGasDecorator(ak, dk),
		didAnte.NewSigVerificationDecorator(ak, dk),
	)
}
