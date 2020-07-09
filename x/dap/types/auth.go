package types

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	sdkexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/spf13/viper"
	"github.com/tendermint/ed25519"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"os"
	"runtime/debug"
	"time"
)

var (
	expectedMinGasPrices       = "0.025" + NativeToken
	approximationGasAdjustment = float64(1.5)
	// TODO: parameterise (or remove) hard-coded gas prices and adjustments

	// simulation signature values used to estimate gas consumption
	simEd25519Pubkey   ed25519tm.PubKeyEd25519
	simEd25519Sig      [ed25519SignatureLen]byte
	simSecp256k1Pubkey secp256k1.PubKeySecp256k1
	simSecp256k1Sig    [ed25519SignatureLen]byte
)

func init() {
	// This decodes a valid hex string into a ed25519Pubkey for use in transaction simulation
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E14")
	copy(simEd25519Pubkey[:], bz)
	//copy(simSecp256k1Pubkey[:], bz)
}

type PubKeyGetter func(ctx sdk.Context, msg IxoMsg) (crypto.PubKey, error)

func consumeSimSigGas(gasmeter sdk.GasMeter, pubkey crypto.PubKey, sig IxoSignature, params auth.Params) {
	simSig := IxoSignature{}
	if binary.Size(sig.SignatureValue) == 0 {
		simSig.SignatureValue = simEd25519Sig
		//simSig.SignatureValue = simEd25519Sig[:]
	}
	simSig.Created = simSig.Created.Add(1) // maximizes signature length
	ModuleCdc := codec.New()
	sigBz := ModuleCdc.MustMarshalBinaryLengthPrefixed(simSig)
	cost := sdk.Gas(len(sigBz) + 6)

	// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
	// number of signers.
	if _, ok := pubkey.(multisig.PubKeyMultisigThreshold); ok {
		cost *= params.TxSigLimit
	}

	gasmeter.ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
}

// verify the signature and increment the sequence. If the account doesn't have
// a pubkey, set it.
func ProcessSig(ctx sdk.Context, acc sdkexported.Account, sig IxoSignature, signBytes []byte, simulate bool, p auth.Params) (updatedAcc sdkexported.Account, err error) {
	pubKey, res := proccessPublicKey(acc, sig, simulate)
	if res != nil {
		return nil, res
	}

	err = acc.SetPubKey(pubKey)
	if err != nil {
		//return nil, sdk.ErrInternal("setting PubKey on signer's account").Result()
		return nil, InvalidPubKey("setting PubKey on signer's account")
	}

	if simulate {
		// Simulated txs should not contain a signature and are not required to
		// contain a pubkey, so we must account for tx size of including an
		// IxoSignature and simulate gas consumption (assuming an ED25519 key).
		consumeSimSigGas(ctx.GasMeter(), pubKey, sig, p)
	}

	// Consume signature gas
	ctx.GasMeter().ConsumeGas(p.SigVerifyCostED25519, "ante verify: ed25519")
	// Verify signature
	if !simulate && !pubKey.VerifyBytes(signBytes, sig.SignatureValue[:]) {
		return nil, Unauthorized("Signature Verification failed. dxp")
	}

	if err = acc.SetSequence(acc.GetSequence() + 1); err != nil {
		panic(err)
	}

	return acc, err
}

func getSignBytes(chainID string, ixoTx IxoTx, acc sdkexported.Account, genesis bool) []byte {
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	return auth.StdSignBytes(
		chainID, accNum, acc.GetSequence(), ixoTx.Fee, ixoTx.Msgs, ixoTx.Memo,
	)
}

/**
used by modules:
did
project
bonds
*/
func NewDefaultAnteHandler(ak auth.AccountKeeper, sk supply.Keeper, pubKeyGetter PubKeyGetter) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
		// determined by the GasMeter. We need access to the context to get the gas
		// meter so we initialize upfront.
		var gasWanted uint64
		var result *sdk.Result
		var gInfo sdk.GasInfo

		if addr := sk.GetModuleAddress(auth.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", auth.FeeCollectorName))
		}

		// all transactions must be of type ixo.IxoTx
		ixoTx, ok := tx.(IxoTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = auth.SetGasMeter(simulate, ctx, 0)
			return newCtx, IntErr("tx must be dxp dedicated type IxoTx")
		}

		params := ak.GetParams(ctx)

		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			// Set high gas price so standard test fee fails
			/*gasPrice := sdk.NewDecCoinFromDec(NativeToken, sdk.NewDec(2000).Quo(IxoDecimals))
			highGasPrice := []sdk.DecCoin{gasPrice}*/
			ctx = ctx.WithMinGasPrices(ctx.MinGasPrices())
			mfd := ante.NewMempoolFeeDecorator()
			antehandler := sdk.ChainAnteDecorators(mfd)
			_, err := antehandler(ctx, tx, false)
			if err != nil {
				return newCtx, IntErr("Decorator should have errored on too low fee for local gasPrice")
			}
		}
		//setting gas to high up
		newCtx = auth.SetGasMeter(simulate, ctx, ixoTx.Fee.Gas)

		// ser= sdkerrors.ResponseDeliverTx(err, gInfo.GasWanted, gInfo.GasUsed)
		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					err = errors.Wrap(
						errors.ErrOutOfGas, fmt.Sprintf(
							"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
							rType.Descriptor, gasWanted, ctx.GasMeter().GasConsumed(),
						),
					)

				default:
					err = errors.Wrap(
						errors.ErrPanic, fmt.Sprintf(
							"recovered: %v\nstack:\n%v", r, string(debug.Stack()),
						),
					)
				}

				result = nil
			}

			gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
		}()

		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err
		}
		estimatedGas := params.TxSizeCostPerByte * sdk.Gas(len(newCtx.TxBytes()))
		/*	ctx.BlockGasMeter().ConsumeGas(
			ctx.GasMeter().GasConsumedToLimit(), "block gas meter",
		)*/
		newCtx.GasMeter().ConsumeGas(estimatedGas, "txSize")
		/*
			if res := auth.ValidateMemo(auth.StdTx{Memo: ixoTx.Memo}, params); res != nil {
				return newCtx, res
			}*/

		// all messages must be of type IxoMsg
		msg, ok := ixoTx.GetMsgs()[0].(IxoMsg)
		if !ok {
			gInfo = sdk.GasInfo{}
			return newCtx, IntErr("msg must be ixo.IxoMsg. dxp")
		}

		// Get pubKey
		pubKey, err := pubKeyGetter(ctx, msg)
		if err != nil {
			return newCtx, err
		}

		// fetch first (and only) signer, who's going to pay the fees
		signerAddr := sdk.AccAddress(pubKey.Address())
		println(fmt.Sprintf("--- the keys %s", string(signerAddr)))
		println(fmt.Sprintf("--- the pub key got %s", string(pubKey.Address().String())))
		signerAcc, err := auth.GetSignerAcc(newCtx, ak, signerAddr)
		if err != nil {
			return newCtx, err
		}
		println(fmt.Sprintf("found account and getting this pubkey %s", signerAcc.GetPubKey().Address().String()))

		// deduct the fees
		if !ixoTx.Fee.Amount.IsZero() {
			err = auth.DeductFees(sk, newCtx, signerAcc, ixoTx.Fee.Amount)
			if err != nil {
				return newCtx, err
			}

			// reload the account as fees have been deducted
			signerAcc = ak.GetAccount(newCtx, signerAcc.GetAddress())
		}

		// check signature, return account with incremented nonce
		//ixoSig := auth.StdSignature{PubKey: pubKey, Signature: ixoTx.GetSignatures()[0].SignatureValue[:]}
		ixoSig := NewSignature(time.Now(), ixoTx.GetSignatures()[0].SignatureValue)
		isGenesis := ctx.BlockHeight() == 0
		signBytes := getSignBytes(newCtx.ChainID(), ixoTx, signerAcc, isGenesis)
		signerAcc, err = ProcessSig(newCtx, signerAcc, ixoSig, signBytes, simulate, params)
		if err != nil {
			return newCtx, err
		}

		ak.SetAccount(newCtx, signerAcc)
		//res = sdk.Result{}

		return newCtx, nil // continue...
	}
}

func proccessPublicKey(acc sdkexported.Account, sig IxoSignature, simulate bool) (crypto.PubKey, error) {
	// If pubkey is not known for account, set it from the types.StdSignature.
	pubKey := acc.GetPubKey()
	if simulate {
		// In simulate mode the transaction comes with no signatures, thus if the
		// account's pubkey is nil, both signature verification and gasKVStore.Set()
		// shall consume the largest amount, i.e. it takes more gas to verify
		// secp256k1 keys than ed25519 ones.
		if pubKey == nil {
			return simSecp256k1Pubkey, nil
		}
		return pubKey, nil
	}

	if pubKey == nil {
		pubKey = acc.GetPubKey()
		if pubKey == nil {
			return nil, InvalidPubKey("PubKey not found")
		}

		if !bytes.Equal(pubKey.Address(), acc.GetAddress()) {
			return nil, Unauthorizedf("PubKey does not match Signer address %s", acc.GetAddress())
		}
	}

	return pubKey, nil
}

func signAndBroadcast(ctx context.CLIContext, msg auth.StdSignMsg, ixoDid exported.IxoDid) (sdk.TxResponse, error) {
	if len(msg.Msgs) != 1 {
		panic("expected one message")
	}

	var privKey ed25519tm.PrivKeyEd25519
	copy(privKey[:], base58.Decode(ixoDid.Secret.SignKey))
	copy(privKey[32:], base58.Decode(ixoDid.VerifyKey))

	signature := SignIxoMessage(msg.Bytes(), privKey)
	tx := NewIxoTxSingleMsg(msg.Msgs[0], msg.Fee, signature, msg.Memo)

	bz, err := ctx.Codec.MarshalJSON(tx)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("Could not marshall tx to binary. Error: %s! ", err.Error())
	}

	res, err := ctx.BroadcastTx(bz)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("Could not broadcast tx. Error: %s! ", err.Error())
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
	signature := IxoSignature{}
	tx := NewIxoTxSingleMsg(
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

func ApproximateFeeForTx(cliCtx context.CLIContext, tx IxoTx, chainId string) (auth.StdFee, error) {

	// Set up a transaction builder
	cdc := cliCtx.Codec
	txEncoder := auth.DefaultTxEncoder
	gasAdjustment := approximationGasAdjustment
	fees := sdk.NewCoins(sdk.NewCoin(NativeToken, sdk.OneInt()))
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

func GenerateOrBroadcastMsgs(cliCtx context.CLIContext, msg sdk.Msg, ixoDid exported.IxoDid) error {
	msgs := []sdk.Msg{msg}
	txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input)

	if cliCtx.GenerateOnly {
		return utils.PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	}

	return CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs, ixoDid)
}

func CompleteAndBroadcastTxCLI(txBldr auth.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg, ixoDid exported.IxoDid) error {
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	//fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
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

	//passphrase, err := keys.GetPassphrase(fromName)
	//if err != nil {
	//	return err
	//}

	// Build the transaction
	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return err
	}

	// Sign and broadcast to a Tendermint node
	res, err := signAndBroadcast(cliCtx, stdSignMsg, ixoDid)
	if err != nil {
		return err
	}

	fmt.Println(res.String())
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.TxHash)
	return nil
}

func CompleteAndBroadcastTxRest(cliCtx context.CLIContext, msg sdk.Msg, ixoDid exported.IxoDid) ([]byte, error) {

	// TODO: implement using txBldr or just remove function completely (ref: #123)

	// Construct dummy tx and approximate and set fee
	tx := NewIxoTxSingleMsg(msg, auth.StdFee{}, IxoSignature{}, "")
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

	// Sign and broadcast to a Tendermint node
	res, err := signAndBroadcast(cliCtx, stdSignMsg, ixoDid)
	if err != nil {
		return nil, err
	}

	output, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}

func SignAndBroadcastTxFromStdSignMsg(cliCtx context.CLIContext,
	msg auth.StdSignMsg, ixoDid exported.IxoDid) (sdk.TxResponse, error) {
	return signAndBroadcast(cliCtx, msg, ixoDid)
}

func SignIxoMessage(signBytes []byte, privKey [ed25519.PrivateKeySize]byte) IxoSignature {
	signatureBytes := ed25519.Sign(&privKey, signBytes)
	return NewSignature(time.Now(), *signatureBytes)
}

func SignAndBroadcastTxCli(cliCtx context.CLIContext, msg sdk.Msg, sovrinDid exported.IxoDid) error {

	bldr := auth.NewTxBuilderFromCLI(cliCtx.Input).
		WithTxEncoder(utils.GetTxEncoder(cliCtx.Codec)).
		WithKeybase(cliCtx.Keybase)
	println("sign 1 ---- ")
	txBldr, err := utils.PrepareTxBuilder(bldr, cliCtx)
	if err != nil {
		return err
	}
	println("sign 2 ---- ")
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
	println("sign 3 ---- ")
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
		println("sign 4 ---- ")
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

func SignAndBroadcastTxRest(cliCtx context.CLIContext, msg sdk.Msg, sovrinDid exported.IxoDid) ([]byte, error) {

	// TODO: implement using txBldr or just remove function completely (ref: #123)

	// Construct dummy tx and approximate and set fee
	tx := NewIxoTxSingleMsg(msg, auth.StdFee{}, IxoSignature{}, "")
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

// ValidateMemo validates the memo size.
/*
func ValidateMemo(stdTx IxoTx, params params.Params) error {
	memoLength := len(stdTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return erro.Wrapf(erro.ErrMemoTooLarge,
			"maximum number of characters is %d but received %d characters",
			params.MaxMemoCharacters, memoLength, )
	}

	return nil
}
*/

func UnknownRequest(m string) error {
	return errors.Wrap(errors.ErrUnknownRequest, m)
}
func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}
func Unauthorizedf(format string, a ...interface{}) error {
	return errors.Wrap(errors.ErrUnauthorized, fmt.Sprintf(format, a...))
}
func IntErr(m string) error {
	return errors.Wrap(errors.ErrPanic, m)
}
func ErrJsonMars(m string) error {
	return errors.Wrapf(errors.ErrJSONMarshal, "Json marshall error %s", m)
}

func InvalidPubKey(m string) error {
	return errors.Wrapf(errors.ErrInvalidPubKey, "PubKey error %s", m)
}
