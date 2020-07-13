package types

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var (
	maxGasWanted = uint64((1 << 63) - 1)
	//_            TxActor = (*IxoTx)(nil)
	_ sdk.Tx = (*IxoTx)(nil)
)

// GetSignBytes returns the signBytes of the tx for a given signer
func StringToAddr(str string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(str)))
}

func DidToAddr(did exported.Did) sdk.AccAddress {
	return StringToAddr(did)
}

type (
	IxoMsg interface {
		sdk.Msg
		GetSignerDid() exported.Did
	}
	IxoSignature struct {
		SignatureValue []byte    `json:"signatureValue" yaml:"signatureValue"`
		Created        time.Time `json:"created" yaml:"created"`
	}
	IxoTx struct {
		sdk.Tx
		Memo       string         `json:"memo" yaml:"memo"`
		Msgs       []sdk.Msg      `json:"payload" yaml:"payload"`
		Fee        auth.StdFee    `json:"fee" yaml:"fee"`
		Signatures []IxoSignature `json:"signatures" yaml:"signatures"`
	}
	TxActor interface {
		GetMsgs() []sdk.Msg
		ValidateBasic() error
		GetMemo() string
		String() string
		GetGas() uint64
		GetFee() sdk.Coins
		FeePayer() sdk.AccAddress
		GetSignBytes(ctx sdk.Context, acc authexported.Account) []byte
		GetSigner() sdk.AccAddress
		GetSignatures() []IxoSignature
		GetFirstSignature() []byte
	}
)

// MarshalYAML returns the YAML representation of the signature.
func (is IxoSignature) MarshalYAML() (interface{}, error) {
	var (
		bz  []byte
		err error
	)

	bz, err = yaml.Marshal(struct {
		SignatureValue string
		Created        string
	}{
		SignatureValue: fmt.Sprintf("%s", is.SignatureValue),
		Created:        is.Created.String(),
	})
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

func (is IxoSignature) String() string {
	return fmt.Sprintf("{%V}", is)
}

/*

func (ss StdSignature) MarshalYAML() (interface{}, error) {
	var (
		bz     []byte
		pubkey string
		err    error
	)

	if ss.PubKey != nil {
		pubkey, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, ss.PubKey)
		if err != nil {
			return nil, err
		}
	}

	bz, err = yaml.Marshal(struct {
		PubKey    string
		Signature string
	}{
		PubKey:    pubkey,
		Signature: fmt.Sprintf("%s", ss.Signature),
	})
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

*/
func NewSignature(created time.Time, signature []byte) IxoSignature {
	return IxoSignature{
		SignatureValue: signature,
		Created:        created,
	}
}

func NewIxoTx(msgs []sdk.Msg, fee auth.StdFee, sigs []IxoSignature, memo string) IxoTx {
	return IxoTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

func NewIxoTxSingleMsg(msg sdk.Msg, fee auth.StdFee, signature IxoSignature, memo string) IxoTx {
	return NewIxoTx([]sdk.Msg{msg}, fee, []IxoSignature{signature}, memo)
}

func (tx IxoTx) GetMsgs() []sdk.Msg { return tx.Msgs }

func (tx IxoTx) GetMemo() string { return tx.Memo }

func (tx IxoTx) ValidateBasic() error {
	// Fee validation
	if tx.Fee.Gas > maxGasWanted {
		return errors.Wrapf(x.ErrGasOverflow, "invalid gas supplied; %d > %d", tx.Fee.Gas, maxGasWanted)
		//return sdk.ErrGasOverflow(fmt.Sprintf("invalid gas supplied; %d > %d", tx.Fee.Gas, maxGasWanted))
	}
	if tx.Fee.Amount.IsAnyNegative() {
		//return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.Amount))
		return errors.Wrapf(errors.ErrInsufficientFee, "invalid fee %s amount provided", tx.Fee.Amount)
	}

	// Signatures validation
	var ixoSigs = tx.GetSignatures()
	if len(ixoSigs) == 0 {
		return errors.Wrap(errors.ErrNoSignatures, "no signers. dxp")
	}
	if len(ixoSigs) != 1 {
		return errors.Wrap(errors.ErrUnauthorized, "there can only be one signer. dxp")
	}
	// Messages validation
	if len(tx.Msgs) != 1 {
		return errors.Wrap(errors.ErrUnauthorized, "there can only be one message. dxp")
	}

	return nil
}
func (tx IxoTx) String() string {
	output, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", string(output))
}
func (tx IxoTx) GetGas() uint64 {
	return tx.Fee.Gas
}
func (tx IxoTx) GetFee() sdk.Coins {
	return tx.Fee.Amount
}
func (tx IxoTx) FeePayer() sdk.AccAddress {
	return tx.GetSigner()
}
func (tx IxoTx) GetSignBytes(ctx sdk.Context, acc authexported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}
	return auth.StdSignBytes(
		chainID, accNum, acc.GetSequence(), tx.Fee, tx.Msgs, tx.Memo,
	)
}
func (tx IxoTx) GetSigner() sdk.AccAddress {
	return tx.GetMsgs()[0].GetSigners()[0]
}
func (tx IxoTx) GetSignatures() []IxoSignature {
	return tx.Signatures
}
func (tx IxoTx) GetFirstSignature() []byte {
	return tx.Signatures[0].SignatureValue
}

/*
GetGas() uint64
GetFee() sdk.Coins
FeePayer() sdk.AccAddress*/

func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {

		if len(txBytes) == 0 {
			return nil, errors.Wrap(errors.ErrTxDecode, "txBytes are empty")
		}

		if string(txBytes[0:1]) == "{" {
			var upTx map[string]interface{}
			er := json.Unmarshal(txBytes, &upTx)
			if er != nil {
				return nil, errors.Wrap(errors.ErrTxDecode, er.Error())
			}

			payloadArray := upTx["payload"].([]interface{})
			if len(payloadArray) != 1 {
				return nil, errors.Wrap(errors.ErrTxDecode, "Multiple messages not supported")
			}
			var tx IxoTx
			er = cdc.UnmarshalJSON(txBytes, &tx)
			if er != nil {
				return nil, errors.Wrap(errors.ErrTxDecode, er.Error())
			}
			return tx, nil
		} else {
			var tx = auth.StdTx{}
			er := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
			if er != nil {
				return nil, errors.Wrap(errors.ErrTxDecode, er.Error())
			}
			return tx, nil
		}
	}
}

type SignTxPack struct {
	ctxCli    context.CLIContext
	msg       sdk.Msg
	did       exported.IxoDid
	signature IxoSignature
	txBldr    auth.TxBuilder
}

func NewDidTxBuild(ctx context.CLIContext, msg sdk.Msg, ixoDid exported.IxoDid) SignTxPack {
	instance := SignTxPack{
		ctxCli: ctx,
		msg:    msg,
		did:    ixoDid,
	}
	instance.signature = instance.getSignature()
	instance.txBldr = auth.NewTxBuilderFromCLI(ctx.Input)
	return instance
}

func (tb SignTxPack) collectMsgs() []sdk.Msg {
	return []sdk.Msg{tb.msg}
}

func (tb SignTxPack) collectSignatures() []IxoSignature {
	return []IxoSignature{tb.signature}
}

func (tb SignTxPack) getSignature() IxoSignature {
	privateKey := tb.did.GetPriKeyByte()
	//signatureBytes := ed25519.Sign(&privateKey, tb.msg.GetSignBytes())
	ixt, ok := tb.msg.(IxoMsg)
	if !ok {
		fmt.Println("the msg type is not IxoMsg")
	}
	signatureBytes := ed25519.Sign(privateKey[:], ixt.GetSignBytes())
	//return NewSignature(time.Now(), signatureBytes[:])
	return NewSignature(time.Now(), signatureBytes)
}

func (tb SignTxPack) SignAndGenerateMessage(t auth.StdSignMsg) IxoTx {
	//sign message
	sign_signature := tb.collectSignatures()
	//collection of messages
	messages := tb.collectMsgs()
	return NewIxoTx(messages, t.Fee, sign_signature, t.Memo)
}

func (tb SignTxPack) printUnsignedStdTx(stdSignMsg auth.StdSignMsg) error {
	if tb.txBldr.SimulateAndExecute() {
		if err := tb.doSimulate(); err != nil {
			return err
		}
	}

	var json []byte
	var err error

	if viper.GetBool(flags.FlagIndentResponse) {
		json, err = tb.ctxCli.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
	} else {
		json, err = tb.ctxCli.Codec.MarshalJSON(stdSignMsg)
	}

	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(tb.ctxCli.Output, "%s\n", json)

	return nil
}
func (tb SignTxPack) doSimulate() error {
	if tb.ctxCli.Simulate {
		txBldr, err := utils.EnrichWithGas(tb.txBldr, tb.ctxCli, tb.collectMsgs())
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}
	return nil
}
func (tb SignTxPack) CompleteAndBroadcastTxCLI() error {
	txBldr, err := utils.PrepareTxBuilder(tb.txBldr, tb.ctxCli)

	if err != nil {
		return err
	}

	stdSignMsg, err := txBldr.BuildSignMsg([]sdk.Msg{tb.msg})
	if err != nil {
		return err
	}

	if err := tb.printUnsignedStdTx(stdSignMsg); err != nil {
		return err
	}

	if tb.ctxCli.Simulate {
		return nil
	}

	fmt.Println("=============== public key ==============")
	fmt.Println(tb.did.GetPubKey())
	fmt.Println("=============== private key ==============")
	fmt.Println(tb.did.GetPriKeyByte())
	fmt.Println("=============== signature equals to==============")
	fmt.Println(tb.signature.SignatureValue)

	if !tb.ctxCli.SkipConfirm {

		var json []byte

		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = tb.ctxCli.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
		} else {
			json, err = tb.ctxCli.Codec.MarshalJSON(stdSignMsg)
		}

		if err != nil {
			panic(err)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)

		ok, err := input.GetConfirmation("Confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	signed_tx_msg := tb.SignAndGenerateMessage(stdSignMsg)
	fmt.Println("=============== pre-tx-signature ==============")
	fmt.Println(signed_tx_msg.GetFirstSignature())

	bz, err := tb.ctxCli.Codec.MarshalJSON(signed_tx_msg)
	if err != nil {
		return fmt.Errorf("Could not marshall tx to binary. Error: %s! ", err.Error())
	}

	res, err := tb.ctxCli.BroadcastTx(bz)
	if err != nil {
		return fmt.Errorf("Could not broadcast tx. Error: %s! ", err.Error())
	}

	fmt.Println(res.String())
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.TxHash)
	return nil
}
