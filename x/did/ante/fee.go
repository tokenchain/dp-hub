package ante

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
)

type DeductFeeDecorator struct {
	supplyKeeper types.SupplyKeeper
	SigVerification
}
type OnlyOneMsgDecorator struct {
	supplyKeeper types.SupplyKeeper
	SigVerification
}

func NewOnlyOneMsgDecorator() OnlyOneMsgDecorator {
	return OnlyOneMsgDecorator{}
}

// Number of messages in the tx must be 1

func (g OnlyOneMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) != 1 {
		return ctx, IntErr("number of messages must be 1")
	}
	fmt.Println("✅  only one message pass ....")
	return next(ctx, tx, simulate)
}
func NewDeductFeeDecorator(ak keeper.AccountKeeper, sk types.SupplyKeeper, p PubKeyGetter) DeductFeeDecorator {
	return DeductFeeDecorator{
		SigVerification: NewSigVerification(ak, p),
		supplyKeeper:    sk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	fmt.Println("--- DeductFeeDecorator .1")
	sv, _, e := dfd.RetrievePubkey(ctx, tx, simulate)
	fmt.Println("--- DeductFeeDecorator .2")
	if e != nil {
		return ctx, InvalidTxDecodePubkeyNotFound(e)
	}

	fmt.Println("--- DeductFeeDecorator .3")
	if addr := dfd.supplyKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
	}

	fmt.Println("--- DeductFeeDecorator .4")
	// deduct the fees
	if !sv.dap_tx.GetFee().IsZero() {

		fmt.Println("--- DeductFeeDecorator .4.1")
		if err = DeductFees(dfd.supplyKeeper, ctx, dfd.GetSignerAccount(ctx), sv.dap_tx.GetFee()); err != nil {
			return ctx, err
		}

		fmt.Println("--- DeductFeeDecorator .5")
		// reload the account as fees have been deducted
		if err = dfd.RefreshAccount(ctx); err != nil {
			return ctx, err
		}
	}

	fmt.Println("✅  fee deduction pass ....")
	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the BankKeeper (in addition to the AccountKeeper, because
// the BankKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(supplyKeeper types.SupplyKeeper, ctx sdk.Context, acc exported.Account, fees sdk.Coins) error {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	fmt.Println("--- DeductFeeDecorator .4.2")
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	fmt.Println("--- DeductFeeDecorator .4.3")
	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	fmt.Println("--- DeductFeeDecorator .4.4")
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	fmt.Println("--- DeductFeeDecorator .4.5")
	spendableCoins := acc.SpendableCoins(blockTime)

	fmt.Println("--- DeductFeeDecorator .4.6")
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	fmt.Println("--- DeductFeeDecorator .4.7")
	err := supplyKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)

	fmt.Println("--- DeductFeeDecorator .4.8")
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}

// ConsumeTxSizeGasDecorator will take in parameters and consume gas proportional
// to the size of tx before calling next AnteHandler. Note, the gas costs will be
// slightly over estimated due to the fact that any given signing account may need
// to be retrieved from state.
//
// CONTRACT: If simulate=true, then signatures must either be completely filled
// in or empty.
// CONTRACT: To use this decorator, signatures of transaction must be represented
// as types.StdSignature otherwise simulate mode will incorrectly estimate gas cost.
type ConsumeTxSizeGasDecorator struct {
	SigVerification
}

func NewDapConsumeGasForTxSizeDecorator(ak keeper.AccountKeeper, p PubKeyGetter) ConsumeTxSizeGasDecorator {
	return ConsumeTxSizeGasDecorator{
		NewSigVerification(ak, p),
	}
}
func (cgts ConsumeTxSizeGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	params := cgts.ak.GetParams(ctx)

	ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(ctx.TxBytes())), "txSize")
	// simulate gas cost for signatures in simulate mode
	if simulate {

		nsv, pubkeytm, e := cgts.RetrievePubkey(ctx, tx, simulate)
		if e != nil {
			return ctx, InvalidTxDecodePubkeyNotFound(e)
		}

		// use stdsignature to mock the size of a full signature

		sigBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(nsv.stdSignature)
		cost := sdk.Gas(len(sigBz) + 6)

		// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
		// number of signers.
		if _, ok := pubkeytm.(multisig.PubKeyMultisigThreshold); ok {
			cost *= params.TxSigLimit
		}

		ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
		//}
	}
	fmt.Println("✅  ConsumeTxSizeGasDecorator pass ....")
	return next(ctx, tx, simulate)
}

type ConsumeVerSignGasDecorator struct {
	SigVerification
}

func NewConsumeVerSignGasDecorator(ak keeper.AccountKeeper, p PubKeyGetter) ConsumeVerSignGasDecorator {
	return ConsumeVerSignGasDecorator{
		NewSigVerification(ak, p),
	}
}
func (svc ConsumeVerSignGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	params := svc.ak.GetParams(ctx)
	if simulate {
		nsv, pk, e := svc.RetrievePubkey(ctx, tx, simulate)
		if e != nil {
			return ctx, InvalidTxDecodePubkeyNotFound(e)
		}
		svc.consumeSimSigGas(nsv, ctx.GasMeter(), nsv.ak.GetParams(ctx), pk)
	} else {
		ctx.GasMeter().ConsumeGas(params.SigVerifyCostED25519, "ante verify: ed25519")
	}
	fmt.Println("✅  ConsumeVerSignGasDecorator pass ....")
	return next(ctx, tx, simulate)
}

// Simulated txs should not contain a signature and are not required to
// contain a pubkey, so we must account for tx size of including an
// IxoSignature and simulate gas consumption (assuming an ED25519 key).
func (svc ConsumeVerSignGasDecorator) consumeSimSigGas(signctx SigVerification, gasmeter sdk.GasMeter, params auth.Params, pubKey crypto.PubKey) {
	simSig := IxoSignature{}
	if len(signctx.stdSignature.Signature) == 0 {
		simSig.SignatureValue = simEd25519Sig[:]
	}
	simSig.Created = simSig.Created.Add(1) // maximizes signature length

	//sigBz := ModuleCdc.MustMarshalBinaryLengthPrefixed(simSig)
	cost := sdk.Gas(29999 + 6)
	//cost := sdk.Gas(len(sigBz) + 6)

	// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
	// number of signers.

	if _, ok := pubKey.(multisig.PubKeyMultisigThreshold); ok {
		cost *= params.TxSigLimit
	}

	gasmeter.ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
}
