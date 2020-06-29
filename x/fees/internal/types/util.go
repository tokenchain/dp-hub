package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"strings"
)

func NewMsgSetFeeContractAuthorisation(feeContractId string, authorised bool,
	payerDid types.SovrinDid) MsgSetFeeContractAuthorisation {
	return MsgSetFeeContractAuthorisation{
		PubKey:        payerDid.VerifyKey,
		PayerDid:      payerDid.Did,
		FeeContractId: feeContractId,
		Authorised:    authorised,
	}
}

func NewMsgCreateFee(fee Fee, creatorDid types.SovrinDid) MsgCreateFee {
	return MsgCreateFee{
		PubKey:     creatorDid.VerifyKey,
		CreatorDid: creatorDid.Did,
		Fee:        fee,
	}
}

func NewMsgCreateFeeContract(feeId, feeContractId string, payer sdk.AccAddress,
	canDeauthorise bool, discountId sdk.Uint, creatorDid types.SovrinDid) MsgCreateFeeContract {
	return MsgCreateFeeContract{
		PubKey:         creatorDid.VerifyKey,
		CreatorDid:     creatorDid.Did,
		FeeId:          feeId,
		FeeContractId:  feeContractId,
		Payer:          payer,
		CanDeauthorise: canDeauthorise,
		DiscountId:     discountId,
	}
}

func NewMsgCreateSubscription(subscriptionId, feeContractId string, maxPeriods sdk.Uint,
	period Period, creatorDid types.SovrinDid) MsgCreateSubscription {
	return MsgCreateSubscription{
		PubKey:         creatorDid.VerifyKey,
		CreatorDid:     creatorDid.Did,
		SubscriptionId: subscriptionId,
		FeeContractId:  feeContractId,
		MaxPeriods:     maxPeriods,
		Period:         period,
	}
}

func NewMsgGrantFeeDiscount(feeContractId string, discountId sdk.Uint,
	recipient sdk.AccAddress, creatorDid types.SovrinDid) MsgGrantFeeDiscount {
	return MsgGrantFeeDiscount{
		PubKey:        creatorDid.VerifyKey,
		SenderDid:     creatorDid.Did,
		FeeContractId: feeContractId,
		DiscountId:    discountId,
		Recipient:     recipient,
	}
}

func NewMsgRevokeFeeDiscount(feeContractId string, holder sdk.AccAddress,
	creatorDid types.SovrinDid) MsgRevokeFeeDiscount {
	return MsgRevokeFeeDiscount{
		PubKey:        creatorDid.VerifyKey,
		SenderDid:     creatorDid.Did,
		FeeContractId: feeContractId,
		Holder:        holder,
	}
}

func NewMsgChargeFee(feeContractId string, creatorDid types.SovrinDid) MsgChargeFee {
	return MsgChargeFee{
		PubKey:        creatorDid.VerifyKey,
		SenderDid:     creatorDid.Did,
		FeeContractId: feeContractId,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err sdk.Error) {
	if strings.TrimSpace(value) == "" {
		return false, sdk.ErrUnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}