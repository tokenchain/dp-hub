package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tokenchain/dp-block/x/did"
	"github.com/tokenchain/dp-block/x/did/ante"
	"github.com/tokenchain/dp-block/x/did/exported"
)

const (
	TypeMsgCreatePaymentTemplate           = "create-payment-template"
	TypeMsgCreatePaymentContract           = "create-payment-contract"
	TypeMsgCreateSubscription              = "create-subscription"
	TypeMsgSetPaymentContractAuthorisation = "set-payment-contract-authorisation"
	TypeMsgGrantDiscount                   = "grant-discount"
	TypeMsgRevokeDiscount                  = "revoke-discount"
	TypeMsgEffectPayment                   = "effect-payment"
)

var (
	_ ante.IxoMsg = MsgCreatePaymentTemplate{}
	_ ante.IxoMsg = MsgCreatePaymentContract{}
	_ ante.IxoMsg = MsgCreateSubscription{}
	_ ante.IxoMsg = MsgSetPaymentContractAuthorisation{}
	_ ante.IxoMsg = MsgGrantDiscount{}
	_ ante.IxoMsg = MsgRevokeDiscount{}
	_ ante.IxoMsg = MsgEffectPayment{}
)

type MsgCreatePaymentTemplate struct {
	CreatorDid      exported.Did    `json:"creator_did" yaml:"creator_did"`
	PaymentTemplate PaymentTemplate `json:"payment_template" yaml:"payment_template"`
}

func (msg MsgCreatePaymentTemplate) Type() string  { return TypeMsgCreatePaymentTemplate }
func (msg MsgCreatePaymentTemplate) Route() string { return RouterKey }
func (msg MsgCreatePaymentTemplate) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.CreatorDid, "CreatorDid"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.CreatorDid) {
		return errors.Wrap(exported.ErrorInvalidDidE, "creator did is invalid")
	}

	// Validate PaymentTemplate
	if err := msg.PaymentTemplate.Validate(); err != nil {
		return err
	}

	return nil
}

func (msg MsgCreatePaymentTemplate) GetSignerDid() exported.Did { return msg.CreatorDid }
func (msg MsgCreatePaymentTemplate) GetSigners() []sdk.AccAddress {
	//return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
	panic("tried to use unimplemented GetSigners function")
}

func (msg MsgCreatePaymentTemplate) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgCreatePaymentTemplate) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgCreatePaymentContract struct {
	CreatorDid        exported.Did   `json:"creator_did" yaml:"creator_did"`
	PaymentTemplateId string         `json:"payment_template_id" yaml:"payment_template_id"`
	PaymentContractId string         `json:"payment_contract_id" yaml:"payment_contract_id"`
	Payer             sdk.AccAddress `json:"payer" yaml:"payer"`
	CanDeauthorise    bool           `json:"can_deauthorise" yaml:"can_deauthorise"`
	DiscountId        sdk.Uint       `json:"discount_id" yaml:"discount_id"`
}

func (msg MsgCreatePaymentContract) Type() string  { return TypeMsgCreatePaymentContract }
func (msg MsgCreatePaymentContract) Route() string { return RouterKey }
func (msg MsgCreatePaymentContract) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.CreatorDid, "CreatorDid"); !valid {
		return err
	} else if msg.Payer.Empty() {
		return exported.ErrInvalidAddress("payer address is empty")
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.CreatorDid) {
		return errors.Wrap(did.ErrorInvalidDid, "creator did is invalid")
	}

	// Check that IDs valid
	if !IsValidPaymentTemplateId(msg.PaymentTemplateId) {
		return ErrInvalidId("payment template id invalid")
	} else if !IsValidPaymentContractId(msg.PaymentContractId) {
		return ErrInvalidId("payment contract id invalid")
	}

	return nil
}

func (msg MsgCreatePaymentContract) GetSignerDid() exported.Did { return msg.CreatorDid }
func (msg MsgCreatePaymentContract) GetSigners() []sdk.AccAddress {
	//return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
	panic("tried to use unimplemented GetSigners function")
}

func (msg MsgCreatePaymentContract) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgCreatePaymentContract) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgCreateSubscription struct {
	PubKey            string       `json:"pub_key" yaml:"pub_key"`
	CreatorDid        exported.Did `json:"creator_did" yaml:"creator_did"`
	SubscriptionId    string       `json:"subscription_id" yaml:"subscription_id"`
	PaymentContractId string       `json:"payment_contract_id" yaml:"payment_contract_id"`
	MaxPeriods        sdk.Uint     `json:"max_periods" yaml:"max_periods"`
	Period            Period       `json:"period" yaml:"period"`
}

func (msg MsgCreateSubscription) Type() string  { return TypeMsgCreateSubscription }
func (msg MsgCreateSubscription) Route() string { return RouterKey }
func (msg MsgCreateSubscription) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.CreatorDid, "CreatorDid"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.CreatorDid) {
		return errors.Wrap(did.ErrorInvalidDid, "creator did is invalid")
	}

	// Check that IDs valid
	if !IsValidSubscriptionId(msg.SubscriptionId) {
		return ErrInvalidId("payment template id invalid")
	}

	// Validate Period
	if err := msg.Period.Validate(); err != nil {
		return err
	}

	return nil
}

func (msg MsgCreateSubscription) GetSignerDid() exported.Did { return msg.CreatorDid }
func (msg MsgCreateSubscription) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{ante.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgCreateSubscription) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgCreateSubscription) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgSetPaymentContractAuthorisation struct {
	PayerDid          exported.Did `json:"payer_did" yaml:"payer_did"`
	PaymentContractId string       `json:"payment_contract_id" yaml:"payment_contract_id"`
	Authorised        bool         `json:"authorised" yaml:"authorised"`
}

func (msg MsgSetPaymentContractAuthorisation) Type() string {
	return TypeMsgSetPaymentContractAuthorisation
}
func (msg MsgSetPaymentContractAuthorisation) Route() string { return RouterKey }
func (msg MsgSetPaymentContractAuthorisation) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PayerDid, "PayerDid"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.PayerDid) {
		return errors.Wrap(did.ErrorInvalidDid, "payer did is invalid")
	}

	// Check that IDs valid
	if !IsValidPaymentContractId(msg.PaymentContractId) {
		return ErrInvalidId("payment contract id invalid")
	}

	return nil
}

func (msg MsgSetPaymentContractAuthorisation) GetSignerDid() exported.Did { return msg.PayerDid }
func (msg MsgSetPaymentContractAuthorisation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{ante.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgSetPaymentContractAuthorisation) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgSetPaymentContractAuthorisation) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgGrantDiscount struct {
	PubKey            string         `json:"pub_key" yaml:"pub_key"`
	SenderDid         exported.Did   `json:"sender_did" yaml:"sender_did"`
	PaymentContractId string         `json:"payment_contract_id" yaml:"payment_contract_id"`
	DiscountId        sdk.Uint       `json:"discount_id" yaml:"discount_id"`
	Recipient         sdk.AccAddress `json:"recipient" yaml:"recipient"`
}

func (msg MsgGrantDiscount) Type() string  { return TypeMsgGrantDiscount }
func (msg MsgGrantDiscount) Route() string { return RouterKey }
func (msg MsgGrantDiscount) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.SenderDid, "SenderDid"); !valid {
		return err
	} else if msg.Recipient.Empty() {
		return exported.ErrInvalidAddress("recipient address is empty")
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.SenderDid) {
		return errors.Wrap(did.ErrorInvalidDid, "sender did is invalid")
	}

	// Check that IDs valid
	if !IsValidPaymentContractId(msg.PaymentContractId) {
		return ErrInvalidId("payment contract id invalid")
	}

	return nil
}

func (msg MsgGrantDiscount) GetSignerDid() exported.Did { return msg.SenderDid }
func (msg MsgGrantDiscount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{ante.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgGrantDiscount) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgGrantDiscount) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgRevokeDiscount struct {
	PubKey            string         `json:"pub_key" yaml:"pub_key"`
	SenderDid         exported.Did   `json:"sender_did" yaml:"sender_did"`
	PaymentContractId string         `json:"payment_contract_id" yaml:"payment_contract_id"`
	Holder            sdk.AccAddress `json:"holder" yaml:"holder"`
}

func (msg MsgRevokeDiscount) Type() string  { return TypeMsgRevokeDiscount }
func (msg MsgRevokeDiscount) Route() string { return RouterKey }
func (msg MsgRevokeDiscount) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.SenderDid, "SenderDid"); !valid {
		return err
	} else if msg.Holder.Empty() {
		return exported.ErrInvalidAddress("holder address is empty")
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.SenderDid) {
		return errors.Wrap(did.ErrorInvalidDid, "sender did is invalid")
	}

	// Check that IDs valid
	if !IsValidPaymentContractId(msg.PaymentContractId) {
		return ErrInvalidId("payment contract id invalid")
	}

	return nil
}

func (msg MsgRevokeDiscount) GetSignerDid() exported.Did { return msg.SenderDid }
func (msg MsgRevokeDiscount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{ante.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgRevokeDiscount) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgRevokeDiscount) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

type MsgEffectPayment struct {
	PubKey            string       `json:"pub_key" yaml:"pub_key"`
	SenderDid         exported.Did `json:"sender_did" yaml:"sender_did"`
	PaymentContractId string       `json:"payment_contract_id" yaml:"payment_contract_id"`
}

func (msg MsgEffectPayment) Type() string  { return TypeMsgEffectPayment }
func (msg MsgEffectPayment) Route() string { return RouterKey }
func (msg MsgEffectPayment) ValidateBasic() error {
	// Check that not empty
	if valid, err := CheckNotEmpty(msg.PubKey, "PubKey"); !valid {
		return err
	} else if valid, err = CheckNotEmpty(msg.SenderDid, "SenderDid"); !valid {
		return err
	}

	// Check that DIDs valid
	if !exported.IsValidDid(msg.SenderDid) {
		return errors.Wrap(did.ErrorInvalidDid, "sender did is invalid")
	}

	// Check that IDs valid
	if !IsValidPaymentContractId(msg.PaymentContractId) {
		return ErrInvalidId("payment contract id invalid")
	}

	return nil
}

func (msg MsgEffectPayment) GetSignerDid() exported.Did { return msg.SenderDid }
func (msg MsgEffectPayment) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{ante.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgEffectPayment) String() string {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg MsgEffectPayment) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}
