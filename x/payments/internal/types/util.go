package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
	"strings"
)

func NewMsgSetPaymentContractAuthorisation(contractId string, authorised bool,
	payerDid types.SovrinDid) MsgSetPaymentContractAuthorisation {
	return MsgSetPaymentContractAuthorisation{
		PubKey:            payerDid.VerifyKey,
		PayerDid:          payerDid.Did,
		PaymentContractId: contractId,
		Authorised:        authorised,
	}
}

func NewMsgCreatePaymentTemplate(template PaymentTemplate,
	creatorDid types.SovrinDid) MsgCreatePaymentTemplate {
	return MsgCreatePaymentTemplate{
		PubKey:          creatorDid.VerifyKey,
		CreatorDid:      creatorDid.Did,
		PaymentTemplate: template,
	}
}

func NewMsgCreatePaymentContract(templateId, contractId string,
	payer sdk.AccAddress, canDeauthorise bool, discountId sdk.Uint,
	creatorDid types.SovrinDid) MsgCreatePaymentContract {
	return MsgCreatePaymentContract{
		PubKey:            creatorDid.VerifyKey,
		CreatorDid:        creatorDid.Did,
		PaymentTemplateId: templateId,
		PaymentContractId: contractId,
		Payer:             payer,
		CanDeauthorise:    canDeauthorise,
		DiscountId:        discountId,
	}
}

func NewMsgCreateSubscription(subscriptionId, contractId string, maxPeriods sdk.Uint,
	period Period, creatorDid types.SovrinDid) MsgCreateSubscription {
	return MsgCreateSubscription{
		PubKey:            creatorDid.VerifyKey,
		CreatorDid:        creatorDid.Did,
		SubscriptionId:    subscriptionId,
		PaymentContractId: contractId,
		MaxPeriods:        maxPeriods,
		Period:            period,
	}
}

func NewMsgGrantDiscount(contractId string, discountId sdk.Uint,
	recipient sdk.AccAddress, creatorDid types.SovrinDid) MsgGrantDiscount {
	return MsgGrantDiscount{
		PubKey:            creatorDid.VerifyKey,
		SenderDid:         creatorDid.Did,
		PaymentContractId: contractId,
		DiscountId:        discountId,
		Recipient:         recipient,
	}
}

func NewMsgRevokeDiscount(contractId string, holder sdk.AccAddress,
	creatorDid types.SovrinDid) MsgRevokeDiscount {
	return MsgRevokeDiscount{
		PubKey:            creatorDid.VerifyKey,
		SenderDid:         creatorDid.Did,
		PaymentContractId: contractId,
		Holder:            holder,
	}
}

func NewMsgEffectPayment(contractId string, creatorDid types.SovrinDid) MsgEffectPayment {
	return MsgEffectPayment{
		PubKey:            creatorDid.VerifyKey,
		SenderDid:         creatorDid.Did,
		PaymentContractId: contractId,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, x.UnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
