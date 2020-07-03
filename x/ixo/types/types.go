package types

import (
	"encoding/json"
	"fmt"
	err "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"gopkg.in/yaml.v2"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var (
	IxoDecimals  = sdk.NewDec(1000)
	maxGasWanted = uint64((1 << 63) - 1)
	// https://sovrin-foundation.github.io/sovrin/spec/did-method-spec-template.html
	// IsValidDid adapted from the above link but assumes no sub-namespaces
	// TODO: ValidDid needs to be updated once we no longer want to be able
	//   to consider project accounts as DIDs (especially in treasury module),
	//   possibly should just be `^did:(dxp:|sov:)([a-zA-Z0-9]){21,22}$`.
	_ sdk.Tx = (*DpTx)(nil)
)

const IxoNativeToken = "dap"

type (
	DpTx struct {
		Msgs       []sdk.Msg     `json:"payload" yaml:"payload"`
		Fee        auth.StdFee   `json:"fee" yaml:"fee"`
		Signatures []DpSignature `json:"signatures" yaml:"signatures"`
		Memo       string        `json:"memo" yaml:"memo"`
	}
	DpSignature struct {
		SignatureValue [64]byte  `json:"signatureValue" yaml:"signatureValue"`
		Created        time.Time `json:"created" yaml:"created"`
	}
	DpMsg interface {
		sdk.Msg
		GetSignerDid() did.Did
	}
)

func StringToAddr(str string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(str)))
}

// MarshalYAML returns the YAML representation of the signature.
func (is DpSignature) MarshalYAML() (interface{}, error) {
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

func NewSignature(created time.Time, signature [64]byte) DpSignature {
	return DpSignature{
		SignatureValue: signature,
		Created:        created,
	}
}

func NewIxoTx(msgs []sdk.Msg, fee auth.StdFee, sigs []DpSignature, memo string) DpTx {
	return DpTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

func NewIxoTxSingleMsg(msg sdk.Msg, fee auth.StdFee, signature DpSignature, memo string) DpTx {
	return NewIxoTx([]sdk.Msg{msg}, fee, []DpSignature{signature}, memo)
}

func (tx DpTx) GetMsgs() []sdk.Msg { return tx.Msgs }

func (tx DpTx) GetMemo() string { return "" }

func (tx DpTx) ValidateBasic() error {
	// Fee validation
	if tx.Fee.Gas > maxGasWanted {
		return err.Wrapf(x.ErrGasOverflow, "invalid gas supplied; %d > %d", tx.Fee.Gas, maxGasWanted)
		//return sdk.ErrGasOverflow(fmt.Sprintf("invalid gas supplied; %d > %d", tx.Fee.Gas, maxGasWanted))
	}
	if tx.Fee.Amount.IsAnyNegative() {
		//return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.Amount))
		return err.Wrapf(err.ErrInsufficientFee, "invalid fee %s amount provided", tx.Fee.Amount)
	}

	// Signatures validation
	var ixoSigs = tx.GetSignatures()
	if len(ixoSigs) == 0 {
		return err.Wrap(err.ErrNoSignatures, "no signers")
	}
	if len(ixoSigs) != 1 {
		return err.Wrap(err.ErrUnauthorized, "there can only be one signer")
	}
	// Messages validation
	if len(tx.Msgs) != 1 {
		return err.Wrap(err.ErrUnauthorized, "there can only be one message")
	}

	return nil
}

func (tx DpTx) GetSignatures() []DpSignature {
	return tx.Signatures
}

func (tx DpTx) String() string {
	output, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", string(output))
}

func (tx DpTx) GetSigner() sdk.AccAddress {
	return tx.GetMsgs()[0].GetSigners()[0]
}

func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {

		if len(txBytes) == 0 {
			return nil, err.Wrap(err.ErrTxDecode, "txBytes are empty")
		}

		if string(txBytes[0:1]) == "{" {
			var upTx map[string]interface{}
			er := json.Unmarshal(txBytes, &upTx)
			if er != nil {
				return nil, err.Wrap(err.ErrTxDecode, er.Error())
			}

			payloadArray := upTx["payload"].([]interface{})
			if len(payloadArray) != 1 {
				return nil, err.Wrap(err.ErrTxDecode, "Multiple messages not supported")
			}

			var tx DpTx
			er = cdc.UnmarshalJSON(txBytes, &tx)
			if er != nil {
				return nil, err.Wrap(err.ErrTxDecode, er.Error())
			}
			return tx, nil
		} else {
			var tx = auth.StdTx{}
			er := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
			if er != nil {
				return nil, err.Wrap(err.ErrTxDecode, er.Error())
			}
			return tx, nil
		}
	}
}
