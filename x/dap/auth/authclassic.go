package auth

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	sdkexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/spf13/viper"
	"github.com/tendermint/ed25519"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tokenchain/dp-block/x/dap/types"
	"github.com/tokenchain/dp-block/x/did/ante"
	"github.com/tokenchain/dp-block/x/did/exported"

	"os"
	"time"
)

const (
	Ed25519SignatureLen = 64
)

var (
	expectedMinGasPrices       = "0.025" + types.NativeToken
	approximationGasAdjustment = float64(1.5)
	// TODO: parameterise (or remove) hard-coded gas prices and adjustments

	// simulation signature values used to estimate gas consumption
	simEd25519Pubkey   ed25519tm.PubKeyEd25519
	simEd25519Sig      [Ed25519SignatureLen]byte
	simSecp256k1Pubkey secp256k1.PubKeySecp256k1
	simSecp256k1Sig    [Ed25519SignatureLen]byte
)




/*



used by modules:
did
project
bonds



*/

/*




store := ctx.KVStore(capKey)
txTest := tx.(txTest)

if txTest.FailOnAnte {
return newCtx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "ante handler failure")
}

_, err = incrementingCounter(t, store, storeKey, txTest.Counter)
if err != nil {
return newCtx, err
}

return newCtx, nil



*/
func consumeSimSigGas(gasmeter sdk.GasMeter, pubkey crypto.PubKey, sig ante.IxoSignature, params auth.Params) {
	simSig := ante.IxoSignature{}
	if binary.Size(sig.SignatureValue) == 0 {
		simSig.SignatureValue = simEd25519Sig[:]
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
func ProcessSig(ctx sdk.Context, acc sdkexported.Account, sig ante.IxoSignature, signBytes []byte, simulate bool, p auth.Params) (updatedAcc sdkexported.Account, err error) {
	pubKey, res := proccessPublicKey(acc, sig, simulate)
	if res != nil {
		return nil, res
	}

	err = acc.SetPubKey(pubKey)
	if err != nil {
		//return nil, sdk.ErrInternal("setting PubKey on signer's account").Result()
		return nil, ante.InvalidPubKey("setting PubKey on signer's account")
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
		return nil, ante.Unauthorized("Signature Verification failed. dxp")
	}

	if err = acc.SetSequence(acc.GetSequence() + 1); err != nil {
		panic(err)
	}

	return acc, err
}

/*
func NewDefaultAnteHandler(ak auth.AccountKeeper, sk supply.Keeper, pubKeyGetter PubKeyGetter) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {

		if addr := sk.GetModuleAddress(auth.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", auth.FeeCollectorName))
		}

		// all transactions must be of type ixo.IxoTx
		ixoTx, ok := tx.(IxoTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = auth.SetGasMeter(simulate, ctx, 0)
			return newCtx, sdk.ErrInternal("tx must be ixo.IxoTx").Result(), true
		}

		params := ak.GetParams(ctx)

		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			res := auth.EnsureSufficientMempoolFees(ctx, ixoTx.Fee)
			if !res.IsOK() {
				return newCtx, res, true
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
					log := fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, ixoTx.Fee.Gas, newCtx.GasMeter().GasConsumed(),
					)
					res = sdk.ErrOutOfGas(log).Result()

					res.GasWanted = ixoTx.Fee.Gas
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()

		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err.Result(), true
		}

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")

		if res := auth.ValidateMemo(auth.StdTx{Memo: ixoTx.Memo}, params); !res.IsOK() {
			return newCtx, res, true
		}

		// all messages must be of type IxoMsg
		msg, ok := ixoTx.GetMsgs()[0].(IxoMsg)
		if !ok {
			return newCtx, sdk.ErrInternal("msg must be ixo.IxoMsg").Result(), true
		}

		// Get pubKey
		pubKey, res := pubKeyGetter(ctx, msg)
		if !res.IsOK() {
			return newCtx, res, true
		}

		// fetch first (and only) signer, who's going to pay the fees
		signerAddr := sdk.AccAddress(pubKey.Address())
		signerAcc, res := auth.GetSignerAcc(newCtx, ak, signerAddr)
		if !res.IsOK() {
			return newCtx, res, true
		}

		// deduct the fees
		if !ixoTx.Fee.Amount.IsZero() {
			res = auth.DeductFees(sk, newCtx, signerAcc, ixoTx.Fee.Amount)
			if !res.IsOK() {
				return newCtx, res, true
			}

			// reload the account as fees have been deducted
			signerAcc = ak.GetAccount(newCtx, signerAcc.GetAddress())
		}

		// check signature, return account with incremented nonce
		ixoSig := auth.StdSignature{PubKey: pubKey, Signature: ixoTx.GetSignatures()[0].SignatureValue[:]}
		isGenesis := ctx.BlockHeight() == 0
		signBytes := getSignBytes(newCtx.ChainID(), ixoTx, signerAcc, isGenesis)
		signerAcc, res = ProcessSig(newCtx, signerAcc, ixoSig, signBytes, simulate, params)
		if !res.IsOK() {
			return newCtx, res, true
		}
		ak.SetAccount(newCtx, signerAcc)
		return newCtx, sdk.Result{GasWanted: ixoTx.Fee.Gas}, false // continue...
	}
}

*/

func proccessPublicKey(acc sdkexported.Account, sig ante.IxoSignature, simulate bool) (crypto.PubKey, error) {
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
			return nil, ante.InvalidPubKey("PubKey not found")
		}

		if !bytes.Equal(pubKey.Address(), acc.GetAddress()) {
			return nil, ante.Unauthorizedf("PubKey does not match Signer address %s", acc.GetAddress())
		}
	}

	return pubKey, nil
}

// sign transactions new and make publishing
func signAndBroadcast(ctx context.CLIContext, msg auth.StdSignMsg, ixoDid exported.IxoDid) (sdk.TxResponse, error) {
	if len(msg.Msgs) != 1 {
		panic("expected one message")
	}
	privKey := exported.RecoverDidEd25519ToPrivateKey(ixoDid)
	signature := SignIxoMessageEd25519(msg.Bytes(), privKey)

	tx := ante.NewIxoTxSingleMsg(msg.Msgs[0], msg.Fee, signature, msg.Memo)
	fmt.Println(msg)
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

func ApproximateFeeForTxDap(cliCtx context.CLIContext, tx ante.IxoTx, chainId string) (auth.StdFee, error) {

	// Set up a transaction builder
	cdc := cliCtx.Codec
	txEncoder := auth.DefaultTxEncoder
	gasAdjustment := approximationGasAdjustment
	fees := sdk.NewCoins(sdk.NewCoin(types.NativeToken, sdk.OneInt()))
	txBldr := auth.NewTxBuilder(txEncoder(cdc), 0, 0, 0, gasAdjustment, true, chainId, tx.Memo, fees, nil)

	// Approximate gas consumption
	txBldr, err := utils.EnrichWithGas(txBldr, cliCtx, tx.Msgs)
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

func ApproximateFeeForTx(cliCtx context.CLIContext, tx auth.StdTx, chainId string) (auth.StdFee, error) {

	// Set up a transaction builder
	cdc := cliCtx.Codec
	txEncoder := auth.DefaultTxEncoder
	gasAdjustment := approximationGasAdjustment
	fees := sdk.NewCoins(sdk.NewCoin(types.NativeToken, sdk.OneInt()))
	txBldr := auth.NewTxBuilder(txEncoder(cdc), 0, 0, 0, gasAdjustment, true, chainId, tx.Memo, fees, nil)

	// Approximate gas consumption
	txBldr, err := utils.EnrichWithGas(txBldr, cliCtx, tx.Msgs)
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

func CompleteAndBroadcastTxRest(cliCtx context.CLIContext, msg sdk.Msg, ixoDid exported.IxoDid) ([]byte, error) {

	// TODO: implement using txBldr or just remove function completely (ref: #123)

	// Construct dummy tx and approximate and set fee
	tx := ante.NewIxoTxSingleMsg(msg, auth.StdFee{}, ante.IxoSignature{}, "")
	chainId := viper.GetString(flags.FlagChainID)
	fee, err := ApproximateFeeForTxDap(cliCtx, tx, chainId)
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

func SignAndBroadcastTxFromStdSignMsg(cliCtx context.CLIContext, msg auth.StdSignMsg, ixoDid exported.IxoDid) (sdk.TxResponse, error) {
	return signAndBroadcast(cliCtx, msg, ixoDid)
}

func SignIxoMessageEd25519(signBytes []byte, privKey [ed25519.PrivateKeySize]byte) ante.IxoSignature {
	signatureBytes := ed25519.Sign(&privKey, signBytes)
	return ante.NewSignature(time.Now(), signatureBytes[:])
}

/*
func SignIxoMessageSecp256k1(signBytes []byte, privKey [32]byte) IxoSignature {
	signatureBytes := ed25519.Sign(&privKey, signBytes)
	return NewSignature(time.Now(), *signatureBytes)
}
*/
func SignAndBroadcastTxCli(cliCtx context.CLIContext, msg sdk.Msg, sovrinDid exported.IxoDid) error {

	bldr := auth.NewTxBuilderFromCLI(cliCtx.Input).
		WithTxEncoder(utils.GetTxEncoder(cliCtx.Codec)).
		WithKeybase(cliCtx.Keybase)

	txBldr, err := utils.PrepareTxBuilder(bldr, cliCtx)
	if err != nil {
		return err
	}

	msgs := []sdk.Msg{msg}

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		var err error // important so that enrichWithGas overwrites txBldr
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
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
		//println("sign 4 ---- ")
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("📝  Confirm transaction before signing and broadcasting", buf)
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
	fmt.Println("========================🗳 BROADCAST COMPLETE")
	fmt.Println(res.String())
	fmt.Println("========================🏁 BLOCK COMMITTED")
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.TxHash)
	return nil
}

func SignAndBroadcastTxRest(cliCtx context.CLIContext, msg sdk.Msg, sovrinDid exported.IxoDid) ([]byte, error) {
	// TODO: implement using txBldr or just remove function completely (ref: #123)
	// Construct dummy tx and approximate and set fee
	tx := ante.NewIxoTxSingleMsg(msg, auth.StdFee{}, ante.IxoSignature{}, "")
	chainId := viper.GetString(flags.FlagChainID)
	fee, err := ApproximateFeeForTxDap(cliCtx, tx, chainId)
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
