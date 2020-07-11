package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
)

type DeductFeeDecorator struct {
	ak           keeper.AccountKeeper
	supplyKeeper bank.Keeper
}

func NewDeductFeeDecorator(ak keeper.AccountKeeper, sk bank.Keeper) DeductFeeDecorator {
	return DeductFeeDecorator{
		ak:           ak,
		supplyKeeper: sk,
	}
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
	ak        keeper.AccountKeeper
	pgetter   PubKeyGetter
	publicKey crypto.PubKey
}

func NewDapConsumeGasForTxSizeDecorator(ak keeper.AccountKeeper, p PubKeyGetter) ConsumeTxSizeGasDecorator {
	return ConsumeTxSizeGasDecorator{
		ak:      ak,
		pgetter: p,
	}
}
func (cgts ConsumeTxSizeGasDecorator) getTempPubkey(ctx sdk.Context, tx sdk.Tx) error {

	sigTx, ok := tx.(IxoTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	// all messages must be of type IxoMsg
	msg, ok := sigTx.GetMsgs()[0].(IxoMsg)
	if !ok {
		//gInfo = sdk.GasInfo{}
		return IntErr("msg must be ixo.IxoMsg. dxp")
	}

	signer := sigTx.GetSigner()
	acc := cgts.ak.GetAccount(ctx, signer)
	if acc != nil {
		p := acc.GetPubKey()
		if p != nil {
			cgts.publicKey = p
			return nil
		}
	}

	// Get pubKey
	pubKey, err := cgts.pgetter(ctx, msg)
	if err != nil {
		return err
	}
	cgts.publicKey = pubKey
	return nil
}
func (cgts ConsumeTxSizeGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(IxoTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}
	params := cgts.ak.GetParams(ctx)
	ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(ctx.TxBytes())), "txSize")

	// simulate gas cost for signatures in simulate mode
	if simulate {
		if e := cgts.getTempPubkey(ctx, sigTx); e != nil {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrTxDecode, "get darkpool pubkey failure %s! ", e.Error())
		}

		// use stdsignature to mock the size of a full signature
		simSig := types.StdSignature{
			Signature: simSecp256k1Sig[:],
			PubKey:    cgts.publicKey,
		}
		sigBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(simSig)
		cost := sdk.Gas(len(sigBz) + 6)

		// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
		// number of signers.
		if _, ok := cgts.publicKey.(multisig.PubKeyMultisigThreshold); ok {
			cost *= params.TxSigLimit
		}

		ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
		//}
	}

	return next(ctx, tx, simulate)
}
