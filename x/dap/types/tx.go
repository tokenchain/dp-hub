package types

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"gopkg.in/yaml.v2"
	"time"
)

var (
	maxGasWanted = uint64((1 << 63) - 1)
)

func StringToAddr(str string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(str)))
}

func DidToAddr(did exported.Did) sdk.AccAddress {
	return StringToAddr(did)
}

type IxoTx struct {
	sdk.Tx
	Msgs       []sdk.Msg      `json:"payload" yaml:"payload"`
	Fee        auth.StdFee    `json:"fee" yaml:"fee"`
	Signatures []IxoSignature `json:"signatures" yaml:"signatures"`
	Memo       string         `json:"memo" yaml:"memo"`
}

//var _ sdk.Tx = IxoTx{}

type IxoSignature struct {
	SignatureValue [ed25519SignatureLen]byte `json:"signatureValue" yaml:"signatureValue"`
	Created        time.Time                 `json:"created" yaml:"created"`
}

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

type IxoMsg interface {
	sdk.Msg
	GetSignerDid() exported.Did
}

func NewSignature(created time.Time, signature [ed25519SignatureLen]byte) IxoSignature {
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
		return errors.Wrap(errors.ErrNoSignatures, "no signers")
	}
	if len(ixoSigs) != 1 {
		return errors.Wrap(errors.ErrUnauthorized, "there can only be one signer")
	}
	// Messages validation
	if len(tx.Msgs) != 1 {
		return errors.Wrap(errors.ErrUnauthorized, "there can only be one message")
	}

	return nil
}
func (tx IxoTx) GetSignatures() []IxoSignature {
	return tx.Signatures
}

func (tx IxoTx) String() string {
	output, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", string(output))
}

func (tx IxoTx) GetSigner() sdk.AccAddress {
	return tx.GetMsgs()[0].GetSigners()[0]
}

var _ sdk.Tx = (*IxoTx)(nil)

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
