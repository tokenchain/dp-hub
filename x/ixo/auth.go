package ixo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	sdk "github.com/cosmos/cosmos-sdk/types"
	erro "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	params "github.com/cosmos/cosmos-sdk/x/auth/types"
	std "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/spf13/viper"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"os"
)

var (
	expectedMinGasPrices       = "0.025" + types.IxoNativeToken
	approximationGasAdjustment = float64(1.5)
	// TODO: parameterise (or remove) hard-coded gas prices and adjustments
)

type PubKeyGetter func(ctx sdk.Context, msg types.IxoMsg) ([32]byte, error)

// EnsureSufficientMempoolFees verifies that the given transaction has supplied
// enough fees to cover a proposer's minimum fees. A result object is returned
// indicating success or failure.
//
// Contract: This should only be called during CheckTx as it cannot be part of
// consensus.
func EnsureSufficientMempoolFees(ctx sdk.Context, stdFee std.StdFee) error {
	minGasPrices := ctx.MinGasPrices()
	if !minGasPrices.IsZero() {
		requiredFees := make(sdk.Coins, len(minGasPrices))

		// Determine the required fees by multiplying each required minimum gas
		// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
		glDec := sdk.NewDec(int64(stdFee.Gas))
		for i, gp := range minGasPrices {
			fee := gp.Amount.Mul(glDec)
			requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		if !stdFee.Amount.IsAnyGTE(requiredFees) {
			return erro.Wrapf(erro.ErrInsufficientFee, "insufficient fees; got: %q required: %q", stdFee.Amount, requiredFees)
		}
	}

	return nil
}
func ProcessSig(ctx sdk.Context, acc exported.Account, signBytes []byte, pubKey [32]byte,
	sig types.IxoSignature, simulate bool, params auth.Params) (updatedAcc exported.Account, res error) {

	if simulate {
		// Simulated txs should not contain a signature and are not required to
		// contain a pubkey, so we must account for tx size of including an
		// IxoSignature and simulate gas consumption (assuming ED25519 key).
		//consumeSimSigGas(ctx.GasMeter(), sig, params)

		// NOTE: this is not the case in the ixo blockchain. The IxoSignature
		// will be blank but still count towards the transaction size given
		// that it uses a fixed length byte array [64]byte as the sig value.
	}

	// Consume signature gas
	ctx.GasMeter().ConsumeGas(params.SigVerifyCostED25519, "ante verify: ed25519")

	// Verify signature
	if !simulate && !VerifySignature(signBytes, pubKey, sig) {
		return nil, erro.Wrap(erro.ErrUnauthorized, "Signature Verification failed")
		//return nil, sdk.ErrUnauthorized("Signature Verification failed").Result()
	}

	if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
		panic(err)
	}

	return acc, res
}

func getSignBytes(chainID string, ixoTx types.IxoTx, acc exported.Account, genesis bool) []byte {
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	return auth.StdSignBytes(
		chainID, accNum, acc.GetSequence(), ixoTx.Fee, ixoTx.Msgs, ixoTx.Memo,
	)
}

func NewAnteHandler(ak auth.AccountKeeper, sk supply.Keeper, pubKeyGetter PubKeyGetter) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {

		if addr := sk.GetModuleAddress(auth.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", auth.FeeCollectorName))
		}

		// all transactions must be of type ixo.IxoTx
		ixoTx, ok := tx.(types.IxoTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = auth.SetGasMeter(simulate, ctx, 0)
			return newCtx, err
		}

		params := ak.GetParams(ctx)

		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			err := EnsureSufficientMempoolFees(ctx, ixoTx.Fee)
			if err != nil {
				return newCtx, err
			}
		}

		newCtx = auth.SetGasMeter(simulate, ctx, ixoTx.Fee.Gas)

		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.
		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					/*
						      log := fmt.Sprintf(
									"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
									rType.Descriptor, ixoTx.Fee.Gas, newCtx.GasMeter().GasConsumed(),
								)
					*/
					res := sdk.GasInfo{}
					// res = sdk.ErrOutOfGas(log).Result()
					err = erro.Wrapf(erro.ErrOutOfGas, "out of gas in location: %v; gasWanted: %d, gasUsed: %d", rType.Descriptor, ixoTx.Fee.Gas, newCtx.GasMeter().GasConsumed())
					res.GasWanted = ixoTx.Fee.Gas
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					//sdk.Result{}
					//abort = true
				default:
					panic(r)
				}
			}
		}()

		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err
		}

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")

		if res := ValidateMemo(types.IxoTx{Memo: ixoTx.Memo}, params); res != nil {
			return newCtx, res
		}

		// fetch first (and only) signer, who's going to pay the fees
		signerAddr := ixoTx.GetSigner()
		signerAcc, res := auth.GetSignerAcc(newCtx, ak, signerAddr)
		if res != nil {
			return newCtx, res
		}

		// deduct the fees
		if !ixoTx.Fee.Amount.IsZero() {
			res = auth.DeductFees(sk, newCtx, signerAcc, ixoTx.Fee.Amount)
			if res != nil {
				return newCtx, res
			}

			// reload the account as fees have been deducted
			signerAcc = ak.GetAccount(newCtx, signerAcc.GetAddress())
		}

		// all messages must be of type IxoMsg
		msg, ok := ixoTx.GetMsgs()[0].(types.IxoMsg)
		if !ok {
			return newCtx, x.IntErr("msg must be ixo.IxoMsg")
		}

		// Get pubKey
		pubKey, res := pubKeyGetter(ctx, msg)
		if res != nil {
			return newCtx, res
		}

		// check signature, return account with incremented nonce
		ixoSig := ixoTx.GetSignatures()[0]
		isGenesis := ctx.BlockHeight() == 0
		signBytes := getSignBytes(newCtx.ChainID(), ixoTx, signerAcc, isGenesis)
		signerAcc, res = ProcessSig(newCtx, signerAcc, signBytes, pubKey, ixoSig, simulate, params)
		if res != nil {
			return newCtx, res
		}

		ak.SetAccount(newCtx, signerAcc)
		//newCtx.WithMinGasPrices(ixoTx.Fee.Gas)
		return newCtx, nil // continue...
	}
}

func signAndBroadcast(ctx context.CLIContext, msg auth.StdSignMsg,
	sovrinDid types.SovrinDid) (sdk.TxResponse, error) {
	if len(msg.Msgs) != 1 {
		panic("expected one message")
	}

	privKey := [64]byte{}
	copy(privKey[:], base58.Decode(sovrinDid.Secret.SignKey))
	copy(privKey[32:], base58.Decode(sovrinDid.VerifyKey))

	signature := SignIxoMessage(msg.Bytes(), privKey)
	tx := types.NewIxoTxSingleMsg(msg.Msgs[0], msg.Fee, signature, msg.Memo)

	bz, err := ctx.Codec.MarshalJSON(tx)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("Could not marshall tx to binary. Error: %s", err.Error())
	}

	res, err := ctx.BroadcastTx(bz)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("Could not broadcast tx. Error: %s", err.Error())
	}

	return res, nil
}

func simulateMsgs(txBldr auth.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (estimated, adjusted uint64, err error) {
	// Build the transaction
	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return
	}

	// Signature set to a blank signature
	signature := types.IxoSignature{}
	signature.Created = signature.Created.Add(1) // maximizes signature length
	tx := types.NewIxoTxSingleMsg(
		stdSignMsg.Msgs[0], stdSignMsg.Fee, signature, stdSignMsg.Memo)

	bz, err := cliCtx.Codec.MarshalJSON(tx)
	if err != nil {
		err = fmt.Errorf("Could not marshall tx to binary. Error: %s", err.Error())
		return
	}

	estimated, adjusted, err = utils.CalculateGas(
		cliCtx.QueryWithData, cliCtx.Codec, bz, txBldr.GasAdjustment())
	return
}

func enrichWithGas(txBldr auth.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (auth.TxBuilder, error) {
	_, adjusted, err := simulateMsgs(txBldr, cliCtx, msgs)
	if err != nil {
		return txBldr, err
	}

	return txBldr.WithGas(adjusted), nil
}

func ApproximateFeeForTx(cliCtx context.CLIContext, tx types.IxoTx, chainId string) (auth.StdFee, error) {

	// Set up a transaction builder
	cdc := cliCtx.Codec
	txEncoder := auth.DefaultTxEncoder
	gasAdjustment := approximationGasAdjustment
	fees := sdk.NewCoins(sdk.NewCoin(types.IxoNativeToken, sdk.OneInt()))
	txBldr := auth.NewTxBuilder(txEncoder(cdc), 0, 0, 0, gasAdjustment, true, chainId, tx.Memo, fees, nil)

	// Approximate gas consumption
	txBldr, err := enrichWithGas(txBldr, cliCtx, tx.Msgs)
	if err != nil {
		return auth.StdFee{}, err
	}

	// Clear fees and set gas-prices to deduce updated fee = (gas * gas-prices)
	signMsg, err := txBldr.WithFees("").WithGasPrices(expectedMinGasPrices).BuildSignMsg(tx.Msgs)
	if err != nil {
		return auth.StdFee{}, err
	}

	return signMsg.Fee, nil
}

func SignAndBroadcastTxCli(cliCtx context.CLIContext, msg sdk.Msg, sovrinDid types.SovrinDid) error {
	txBldr, err := utils.PrepareTxBuilder(auth.NewTxBuilderFromCLI(cliCtx.Input), cliCtx)
	if err != nil {
		return err
	}

	msgs := []sdk.Msg{msg}

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		var err error // important so that enrichWithGas overwrites txBldr
		txBldr, err = enrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	// Build the transaction
	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return err
	}

	// Sign and broadcast
	res, err := signAndBroadcast(cliCtx, stdSignMsg, sovrinDid)
	if err != nil {
		return err
	}

	fmt.Println(res.String())
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.TxHash)
	return nil
}

func SignAndBroadcastTxRest(cliCtx context.CLIContext, msg sdk.Msg, sovrinDid types.SovrinDid) ([]byte, error) {

	// TODO: implement using txBldr or just remove function completely (ref: #123)

	// Construct dummy tx and approximate and set fee
	tx := types.NewIxoTxSingleMsg(msg, auth.StdFee{}, types.IxoSignature{}, "")
	chainId := viper.GetString(flags.FlagChainID)
	fee, err := ApproximateFeeForTx(cliCtx, tx, chainId)
	if err != nil {
		return nil, err
	}

	// Construct sign message
	stdSignMsg := auth.StdSignMsg{
		Fee:  fee,
		Msgs: []sdk.Msg{msg},
		Memo: "",
	}

	// Sign and broadcast
	res, err := signAndBroadcast(cliCtx, stdSignMsg, sovrinDid)
	if err != nil {
		return nil, err
	}

	output, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}

func SignAndBroadcastTxFromStdSignMsg(cliCtx context.CLIContext, msg auth.StdSignMsg, sovrinDid types.SovrinDid) (sdk.TxResponse, error) {
	return signAndBroadcast(cliCtx, msg, sovrinDid)
}

// ValidateMemo validates the memo size.
func ValidateMemo(stdTx types.IxoTx, params params.Params) error {
	memoLength := len(stdTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return erro.Wrapf(erro.ErrMemoTooLarge,
			"maximum number of characters is %d but received %d characters",
			params.MaxMemoCharacters, memoLength, )
	}

	return nil
}
