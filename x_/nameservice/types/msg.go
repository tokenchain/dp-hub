package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	r "github.com/cosmos/cosmos-sdk/types/errors"

)

//type _ MsgSetName = sdk.Msg
var _ sdk.Msg = MsgSetName{}
var _ sdk.Msg = MsgDeleteName{}
var _ sdk.Msg = MsgBuyName{}

// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string         `json:"name"`
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// MsgDeleteName defines a DeleteName message
type MsgDeleteName struct {
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// MsgBuyName defines the BuyName message
type MsgBuyName struct {
	Name  string         `json:"name"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Route should return the name of the module
func (osg MsgSetName) Route() string { return RouterKey }

// Type should return the action
func (osg MsgSetName) Type() string { return "set_name" }

// ValidateBasic runs stateless checks on the message
func (osg MsgSetName) ValidateBasic() error {
	if osg.Owner.Empty() {
		return r.Wrap(r.ErrInvalidAddress, osg.Owner.String())
	}
	if len(osg.Name) == 0 || len(osg.Value) == 0 {
		return r.Wrap(r.ErrUnknownAddress, "Name and/or Value cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (osg MsgSetName) GetSignBytes() []byte {
	//	return sdk.MustSortJSON()
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(osg))
}

// GetSigners defines whose signature is required
func (osg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{osg.Owner}
}

// NewMsgBuyName is the constructor function for MsgBuyName
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (osg MsgBuyName) Route() string { return RouterKey }

// Type should return the action
func (osg MsgBuyName) Type() string { return "buy_name" }

// ValidateBasic runs stateless checks on the message
func (osg MsgBuyName) ValidateBasic() error {
	if osg.Buyer.Empty() {
		return r.Wrap(r.ErrInvalidAddress, osg.Buyer.String())
	}
	if len(osg.Name) == 0 {
		return r.Wrap(r.ErrUnknownAddress,"Name cannot be empty")
	}
	if !osg.Bid.IsAllPositive() {
		return r.Wrap(r.ErrInsufficientFunds,"")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (osg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(osg))
}

// GetSigners defines whose signature is required
func (osg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{osg.Buyer}
}

//var _ sdk.Msg = (*MsgDeleteName)(nil)

// NewMsgDeleteName is a constructor function for MsgDeleteName
func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		Name:  name,
		Owner: owner,
	}
}

// Route should return the name of the module
func (osg MsgDeleteName) Route() string { return RouterKey }

// Type should return the action
func (osg MsgDeleteName) Type() string { return "delete_name" }

// ValidateBasic runs stateless checks on the message
func (osg MsgDeleteName) ValidateBasic() error {
	if osg.Owner.Empty() {
		return r.Wrap(r.ErrInvalidAddress, osg.Owner.String())
	}
	if len(osg.Name) == 0 {
		return r.Wrap(r.ErrUnknownRequest, "name cannot be empty in validation")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (osg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(osg))
}

// GetSigners defines whose signature is required
func (osg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{osg.Owner}
}

