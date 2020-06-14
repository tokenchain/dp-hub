package types

import (
	"encoding/json"
	"github.com/tokenchain/ixo-blockchain/x/did"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

type MsgCreateBond struct {
	SignBytes string  `json:"signBytes" yaml:"signBytes"`
	TxHash    string  `json:"txHash" yaml:"txHash"`
	SenderDid ixo.Did `json:"senderDid" yaml:"senderDid"`
	BondDid   ixo.Did `json:"bondDid" yaml:"bondDid"`
	PubKey    string  `json:"pubKey" yaml:"pubKey"`
	Data      BondDoc `json:"data" yaml:"data"`
}

var _ sdk.Msg = MsgCreateBond{}

func (msg MsgCreateBond) Type() string  { return "create-bond" }
func (msg MsgCreateBond) Route() string { return RouterKey }
func (msg MsgCreateBond) ValidateBasic() sdk.Error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.BondDid, "BondDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Data.CreatedBy, "CreatedBy"); !valid {
		return err
	}

	// Check that DIDs valid
	if !ixo.IsValidDid(msg.BondDid) {
		return did.ErrorInvalidDid(DefaultCodespace, "bond did is invalid")
	} else if !ixo.IsValidDid(msg.SenderDid) {
		return did.ErrorInvalidDid(DefaultCodespace, "sender did is invalid")
	}

	// No need for extra checks on Data since a blank status is valid

	return nil
}

func (msg MsgCreateBond) GetBondDid() ixo.Did   { return msg.BondDid }
func (msg MsgCreateBond) GetSenderDid() ixo.Did { return msg.SenderDid }
func (msg MsgCreateBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.GetBondDid())}
}

func (msg MsgCreateBond) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgCreateBond) GetPubKey() string     { return msg.PubKey }
func (msg MsgCreateBond) GetStatus() BondStatus { return msg.Data.Status }
func (msg *MsgCreateBond) SetStatus(status BondStatus) {
	msg.Data.Status = status
}

func (msg MsgCreateBond) GetSignBytes() []byte {
	return []byte(msg.SignBytes)
}

func (msg MsgCreateBond) IsNewDid() bool { return true }

var _ StoredBondDoc = (*MsgCreateBond)(nil)

type MsgUpdateBondStatus struct {
	SignBytes string              `json:"signBytes" yaml:"signBytes"`
	SenderDid ixo.Did             `json:"senderDid" yaml:"senderDid"`
	BondDid   ixo.Did             `json:"bondDid" yaml:"bondDid"`
	Data      UpdateBondStatusDoc `json:"data" yaml:"data"`
}

func (msg MsgUpdateBondStatus) Type() string  { return "update-bond-status" }
func (msg MsgUpdateBondStatus) Route() string { return RouterKey }

func (msg MsgUpdateBondStatus) ValidateBasic() sdk.Error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.BondDid, "BondDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.SenderDid, "SenderDid"); !valid {
		return err
	}

	// Check that DIDs valid
	if !ixo.IsValidDid(msg.BondDid) {
		return did.ErrorInvalidDid(DefaultCodespace, "bond did is invalid")
	} else if !ixo.IsValidDid(msg.SenderDid) {
		return did.ErrorInvalidDid(DefaultCodespace, "sender did is invalid")
	}

	// No need for extra checks on Data since a blank status is valid
	// IsValidProgressionFrom checked by the handler

	return nil
}

func (msg MsgUpdateBondStatus) GetSignBytes() []byte {
	return []byte(msg.SignBytes)
}

func (msg MsgUpdateBondStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.GetBondDid())}
}

func (msg MsgUpdateBondStatus) GetBondDid() ixo.Did {
	return msg.BondDid
}

func (msg MsgUpdateBondStatus) GetStatus() BondStatus {
	return msg.Data.Status
}

func (msg MsgUpdateBondStatus) IsNewDid() bool     { return false }
func (msg MsgUpdateBondStatus) IsWithdrawal() bool { return false }
