package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	exported "github.com/tokenchain/dp-hub/x/did/exported"
	"strings"
)

func NewMsgSetPaymentContractAuthorisation(contractId string, authorised bool,
	payerDid exported.Did) MsgSetPaymentContractAuthorisation {
	return MsgSetPaymentContractAuthorisation{
		PayerDid:          payerDid,
		PaymentContractId: contractId,
		Authorised:        authorised,
	}
}

func NewMsgCreatePaymentTemplate(template PaymentTemplate,
	creatorDid exported.Did) MsgCreatePaymentTemplate {
	return MsgCreatePaymentTemplate{
		CreatorDid:      creatorDid,
		PaymentTemplate: template,
	}
}

func NewMsgCreatePaymentContract(templateId, contractId string,
	payer sdk.AccAddress, canDeauthorise bool, discountId sdk.Uint,
	creatorDid exported.Did) MsgCreatePaymentContract {
	return MsgCreatePaymentContract{
		CreatorDid:        creatorDid,
		PaymentTemplateId: templateId,
		PaymentContractId: contractId,
		Payer:             payer,
		CanDeauthorise:    canDeauthorise,
		DiscountId:        discountId,
	}
}

func NewMsgCreateSubscription(subscriptionId, contractId string, maxPeriods sdk.Uint,
	period Period, creatorDid exported.Did) MsgCreateSubscription {
	return MsgCreateSubscription{
		CreatorDid:        creatorDid,
		SubscriptionId:    subscriptionId,
		PaymentContractId: contractId,
		MaxPeriods:        maxPeriods,
		Period:            period,
	}
}

func NewMsgGrantDiscount(contractId string, discountId sdk.Uint,
	recipient sdk.AccAddress, creatorDid exported.Did) MsgGrantDiscount {
	return MsgGrantDiscount{
		SenderDid:         creatorDid,
		PaymentContractId: contractId,
		DiscountId:        discountId,
		Recipient:         recipient,
	}
}

func NewMsgRevokeDiscount(contractId string, holder sdk.AccAddress,
	creatorDid exported.Did) MsgRevokeDiscount {
	return MsgRevokeDiscount{
		SenderDid:         creatorDid,
		PaymentContractId: contractId,
		Holder:            holder,
	}
}

func NewMsgEffectPayment(contractId string, creatorDid exported.Did) MsgEffectPayment {
	return MsgEffectPayment{
		SenderDid:         creatorDid,
		PaymentContractId: contractId,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, exported.UnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
