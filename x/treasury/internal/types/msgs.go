package types

import (
	"encoding/json"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSend           = "send"
	TypeMsgOracleTransfer = "oracle-transfer"
	TypeMsgOracleMint     = "oracle-mint"
	TypeMsgOracleBurn     = "oracle-burn"
)

var (
	_ types.DpMsg = MsgSend{}
	_ types.DpMsg = MsgOracleTransfer{}
	_ types.DpMsg = MsgOracleMint{}
	_ types.DpMsg = MsgOracleBurn{}
)

type MsgSend struct {
	PubKey  string     `json:"pub_key" yaml:"pub_key"`
	FromDid did.Did `json:"from_did" yaml:"from_did"`
	ToDid   did.Did `json:"to_did" yaml:"to_did"`
	Amount  sdk.Coins  `json:"amount" yaml:"amount"`
}

func (msg MsgSend) Type() string  { return TypeMsgSend }
func (msg MsgSend) Route() string { return RouterKey }
func (msg MsgSend) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.FromDid, "FromDid"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.ToDid, "ToDid"); !valid {
		return err
	}

	// Check that DIDs valid
	if !did.IsValidDid(msg.FromDid) {
		return x.ErrInvalidDid("from did is invalid")
	} else if !did.IsValidDid(msg.ToDid) {
		return x.ErrInvalidDid("to did is invalid")
	}

	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	return nil
}

func (msg MsgSend) GetSignerDid() did.Did { return msg.FromDid }
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{did.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgSend) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgSend) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgOracleTransfer struct {
	PubKey    string     `json:"pub_key" yaml:"pub_key"`
	OracleDid did.Did `json:"oracle_did" yaml:"oracle_did"`
	FromDid   did.Did `json:"from_did" yaml:"from_did"`
	ToDid     did.Did `json:"to_did" yaml:"to_did"`
	Amount    sdk.Coins  `json:"amount" yaml:"amount"`
	Proof     string     `json:"proof" yaml:"proof"`
}

func (msg MsgOracleTransfer) Type() string  { return TypeMsgOracleTransfer }
func (msg MsgOracleTransfer) Route() string { return RouterKey }
func (msg MsgOracleTransfer) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.OracleDid, "OracleDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.FromDid, "FromDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.ToDid, "ToDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Proof, "Proof"); !valid {
		return err
	}

	// Check that DIDs valid
	if !did.IsValidDid(msg.OracleDid) {
		return x.ErrInvalidDid("oracle did is invalid")
	} else if !did.IsValidDid(msg.FromDid) {
		return x.ErrInvalidDid("from did is invalid")
	} else if !did.IsValidDid(msg.ToDid) {
		return x.ErrInvalidDid("to did is invalid")
	}

	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	return nil
}

func (msg MsgOracleTransfer) GetSignerDid() did.Did { return msg.OracleDid }
func (msg MsgOracleTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{did.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgOracleTransfer) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgOracleTransfer) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgOracleMint struct {
	PubKey    string     `json:"pub_key" yaml:"pub_key"`
	OracleDid did.Did `json:"oracle_did" yaml:"oracle_did"`
	ToDid     did.Did `json:"to_did" yaml:"to_did"`
	Amount    sdk.Coins  `json:"amount" yaml:"amount"`
	Proof     string     `json:"proof" yaml:"proof"`
}

func (msg MsgOracleMint) Type() string  { return TypeMsgOracleMint }
func (msg MsgOracleMint) Route() string { return RouterKey }
func (msg MsgOracleMint) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.OracleDid, "OracleDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.ToDid, "ToDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Proof, "Proof"); !valid {
		return err
	}

	// Check that DIDs valid
	if !did.IsValidDid(msg.OracleDid) {
		return x.ErrInvalidDid("oracle did is invalid")
	} else if !did.IsValidDid(msg.ToDid) {
		return x.ErrInvalidDid("to did is invalid")
	}

	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	return nil
}

func (msg MsgOracleMint) GetSignerDid() did.Did { return msg.OracleDid }
func (msg MsgOracleMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{did.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgOracleMint) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgOracleMint) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgOracleBurn struct {
	PubKey    string     `json:"pub_key" yaml:"pub_key"`
	OracleDid did.Did `json:"oracle_did" yaml:"oracle_did"`
	FromDid   did.Did `json:"from_did" yaml:"from_did"`
	Amount    sdk.Coins  `json:"amount" yaml:"amount"`
	Proof     string     `json:"proof" yaml:"proof"`
}

func (msg MsgOracleBurn) Type() string  { return TypeMsgOracleBurn }
func (msg MsgOracleBurn) Route() string { return RouterKey }
func (msg MsgOracleBurn) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.OracleDid, "OracleDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.FromDid, "FromDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Proof, "Proof"); !valid {
		return err
	}

	// Check that DIDs valid
	if !did.IsValidDid(msg.OracleDid) {
		return x.ErrInvalidDid("oracle did is invalid")
	} else if !did.IsValidDid(msg.FromDid) {
		return x.ErrInvalidDid("from did is invalid")
	}
	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	return nil
}

func (msg MsgOracleBurn) GetSignerDid() did.Did { return msg.OracleDid }
func (msg MsgOracleBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{did.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgOracleBurn) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgOracleBurn) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}
