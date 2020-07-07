package types

import (
	"encoding/json"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSend           = "send"
	TypeMsgOracleTransfer = "oracle-transfer"
	TypeMsgOracleMint     = "oracle-mint"
	TypeMsgOracleBurn     = "oracle-burn"
)

var (
	_ types.IxoMsg = MsgSend{}
	_ types.IxoMsg = MsgOracleTransfer{}
	_ types.IxoMsg = MsgOracleMint{}
	_ types.IxoMsg = MsgOracleBurn{}
)

type MsgSend struct {
	FromDid     exported.Did `json:"from_did" yaml:"from_did"`
	ToDidOrAddr exported.Did `json:"to_did" yaml:"to_did"`
	Amount      sdk.Coins    `json:"amount" yaml:"amount"`
}

func (msg MsgSend) Type() string  { return TypeMsgSend }
func (msg MsgSend) Route() string { return RouterKey }
func (msg MsgSend) ValidateBasic() error {
	// Check that not empty

	if valid, err := CheckNotEmpty(msg.FromDid, "FromDid"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.ToDidOrAddr, "ToDidOrAddr"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.FromDid) {
		return x.ErrInvalidDid("from did is invalid")
	}

	_, err := sdk.AccAddressFromBech32(msg.ToDidOrAddr)
	if err != nil && !exported.IsValidDid(msg.ToDidOrAddr) {
		return x.InvalidAddress("recipient is neither a did nor an address")
	}

	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	return nil
}

func (msg MsgSend) GetSignerDid() exported.Did { return msg.FromDid }
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
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
	OracleDid   exported.Did `json:"oracle_did" yaml:"oracle_did"`
	FromDid     exported.Did `json:"from_did" yaml:"from_did"`
	ToDidOrAddr exported.Did `json:"to_did" yaml:"to_did"`
	Amount      sdk.Coins    `json:"amount" yaml:"amount"`
	Proof       string       `json:"proof" yaml:"proof"`
}

func (msg MsgOracleTransfer) Type() string  { return TypeMsgOracleTransfer }
func (msg MsgOracleTransfer) Route() string { return RouterKey }
func (msg MsgOracleTransfer) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.OracleDid, "OracleDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.FromDid, "FromDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.ToDidOrAddr, "ToDidOrAddr"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Proof, "Proof"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.OracleDid) {
		return x.ErrInvalidDid("oracle did is invalid")
	} else if !exported.IsValidDid(msg.FromDid) {
		return x.ErrInvalidDid("from did is invalid")
	}

	_, err := sdk.AccAddressFromBech32(msg.ToDidOrAddr)
	if err != nil && !exported.IsValidDid(msg.ToDidOrAddr) {
		return x.InvalidAddress("recipient is neither a did nor an address")
	}
	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	return nil
}

func (msg MsgOracleTransfer) GetSignerDid() exported.Did { return msg.OracleDid }
func (msg MsgOracleTransfer) GetSigners() []sdk.AccAddress {
	panic("tried to use unimplemented GetSigners function")
	//	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
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
	OracleDid   exported.Did `json:"oracle_did" yaml:"oracle_did"`
	ToDidOrAddr exported.Did `json:"to_did" yaml:"to_did"`
	Amount      sdk.Coins    `json:"amount" yaml:"amount"`
	Proof       string       `json:"proof" yaml:"proof"`
}

func (msg MsgOracleMint) Type() string  { return TypeMsgOracleMint }
func (msg MsgOracleMint) Route() string { return RouterKey }
func (msg MsgOracleMint) ValidateBasic() error {
	// Check that not empty

	if valid, err := CheckNotEmpty(msg.OracleDid, "OracleDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.ToDidOrAddr, "ToDidOrAddr"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Proof, "Proof"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.OracleDid) {
		return x.ErrInvalidDid("oracle did is invalid")
	}

	_, err := sdk.AccAddressFromBech32(msg.ToDidOrAddr)
	if err != nil && !exported.IsValidDid(msg.ToDidOrAddr) {
		return x.InvalidAddress("recipient is neither a did nor an address")
	}

	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	return nil
}

func (msg MsgOracleMint) GetSignerDid() exported.Did { return msg.OracleDid }
func (msg MsgOracleMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
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
	OracleDid exported.Did `json:"oracle_did" yaml:"oracle_did"`
	FromDid   exported.Did `json:"from_did" yaml:"from_did"`
	Amount    sdk.Coins    `json:"amount" yaml:"amount"`
	Proof     string       `json:"proof" yaml:"proof"`
}

func (msg MsgOracleBurn) Type() string  { return TypeMsgOracleBurn }
func (msg MsgOracleBurn) Route() string { return RouterKey }
func (msg MsgOracleBurn) ValidateBasic() error {
	// Check that not empty

	if valid, err := CheckNotEmpty(msg.OracleDid, "OracleDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.FromDid, "FromDid"); !valid {
		return err
	} else if valid, err := CheckNotEmpty(msg.Proof, "Proof"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.OracleDid) {
		return x.ErrInvalidDid("oracle did is invalid")
	} else if !exported.IsValidDid(msg.FromDid) {
		return x.ErrInvalidDid("from did is invalid")
	}
	// Check amount (note: validity also checks that coins are positive)
	if !msg.Amount.IsValid() {
		return x.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	return nil
}

func (msg MsgOracleBurn) GetSignerDid() exported.Did { return msg.OracleDid }
func (msg MsgOracleBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
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
