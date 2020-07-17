package payments

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/payments/internal/keeper"
)

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) []abci.ValidatorUpdate {

	iterator := keeper.GetSubscriptionIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		subscription := keeper.MustGetSubscriptionByKey(ctx, iterator.Key())

		// Skip if should not effect
		if !subscription.ShouldEffect(ctx) {
			continue
		}

		// Effect subscription payment
		err := keeper.EffectSubscriptionPayment(ctx, subscription.Id)
		if err != nil {
			panic(err) // TODO: maybe shouldn't panic?
		}

		// Note: if payment can be re-effected immediately, this should be done
		// in the next block to prevent spending too much time effecting payments

		// Get updated subscription
		subscription, err = keeper.GetSubscription(ctx, subscription.Id)
		if err != nil {
			panic(err)
		}

		// Delete subscription if it has completed
		if subscription.IsComplete() {
			// TODO: delete subscription
		}

		// Note: no need to save the subscription, as it is being saved by the
		// functions operating on it, such as EffectSubscriptionPayment()
	}
	return []abci.ValidatorUpdate{}
}

func NewHandler(k Keeper, bk bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgSetPaymentContractAuthorisation:
			return handleMsgSetPaymentContractAuthorisation(ctx, k, msg)
		case MsgCreatePaymentTemplate:
			return handleMsgCreatePaymentTemplate(ctx, k, bk, msg)
		case MsgCreatePaymentContract:
			return handleMsgCreatePaymentContract(ctx, k, bk, msg)
		case MsgCreateSubscription:
			return handleMsgCreateSubscription(ctx, k, msg)
		case MsgGrantDiscount:
			return handleMsgGrantDiscount(ctx, k, msg)
		case MsgRevokeDiscount:
			return handleMsgRevokeDiscount(ctx, k, msg)
		case MsgEffectPayment:
			return handleMsgEffectPayment(ctx, k, bk, msg)
		default:
			return nil,exported.UnknownRequest("No match for message type.")
		}
	}
}

func handleMsgSetPaymentContractAuthorisation(ctx sdk.Context, k Keeper, msg MsgSetPaymentContractAuthorisation) (*sdk.Result, error) {

	// Get payment contract
	contract, err := k.GetPaymentContract(ctx, msg.PaymentContractId)
	if err != nil {
		return nil, err
	}

	// Confirm that signer is actually the payer in the payment contract
	payerAddr := exported.DidToAddr(msg.PayerDid)
	if !payerAddr.Equals(contract.Payer) {
		return nil, exported.InvalidAddress("signer must be payment contract payer")
	}

	// Set authorised status
	err = k.SetPaymentContractAuthorised(ctx, msg.PaymentContractId, msg.Authorised)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleMsgCreatePaymentTemplate(ctx sdk.Context, k Keeper, bk bank.Keeper, msg MsgCreatePaymentTemplate) (*sdk.Result, error) {

	// Ensure that payment template doesn't already exist
	if k.PaymentTemplateExists(ctx, msg.PaymentTemplate.Id) {
		return nil, exported.ErrAlreadyExists(fmt.Sprintf("payment template '%s' already exists", msg.PaymentTemplate.Id))
	}

	// Ensure that payment template ID is not reserved
	if k.PaymentTemplateIdReserved(msg.PaymentTemplate.Id) {
		return nil, exported.Unauthorized(fmt.Sprintf("%s is not allowed as it is "+
			"using a reserved prefix", msg.PaymentTemplate.Id))
	}

	// Create and validate payment template
	if err := msg.PaymentTemplate.Validate(); err != nil {
		return nil, err
	}

	// Ensure no blacklisted address in wallet distribution
	for _, share := range msg.PaymentTemplate.WalletDistribution {
		if bk.BlacklistedAddr(share.Address) {
			return nil, exported.Unauthorized(fmt.Sprintf("%s is not allowed "+
				"to receive transactions", share.Address))
		}
	}

	// Submit payment template
	k.SetPaymentTemplate(ctx, msg.PaymentTemplate)

	return &sdk.Result{}, nil
}

func handleMsgCreatePaymentContract(ctx sdk.Context, k Keeper, bk bank.Keeper,
	msg MsgCreatePaymentContract) (*sdk.Result, error) {

	// Ensure that payment contract doesn't already exist
	if k.PaymentContractExists(ctx, msg.PaymentContractId) {
		return nil, exported.ErrAlreadyExists(fmt.Sprintf(
			"payment contract '%s' already exists", msg.PaymentContractId))
	}

	// Ensure that payment contract ID is not reserved
	if k.PaymentContractIdReserved(msg.PaymentContractId) {
		return nil, exported.Unauthorized(fmt.Sprintf("%s is not allowed as it is "+
			"using a reserved prefix", msg.PaymentContractId))
	}

	// Ensure payer is not a blacklisted address
	if bk.BlacklistedAddr(msg.Payer) {
		return nil, exported.Unauthorized(fmt.Sprintf("%s is not allowed "+
			"to receive transactions", msg.Payer))
	}

	// Confirm that payment template exists
	if !k.PaymentTemplateExists(ctx, msg.PaymentTemplateId) {
		return nil, exported.IntErr("invalid payment template")
	}

	// Create payment contract and validate
	creatorAddr := exported.DidToAddr(msg.CreatorDid)
	contract := NewPaymentContract(msg.PaymentContractId, msg.PaymentTemplateId,
		creatorAddr, msg.Payer, msg.CanDeauthorise, false, msg.DiscountId)
	if err := contract.Validate(); err != nil {
		return nil, err
	}

	// Submit payment contract
	k.SetPaymentContract(ctx, contract)

	return &sdk.Result{}, nil
}

func handleMsgCreateSubscription(ctx sdk.Context, k Keeper,
	msg MsgCreateSubscription) (*sdk.Result, error) {

	// Ensure that subscription doesn't already exist
	if k.SubscriptionExists(ctx, msg.SubscriptionId) {
		return nil, exported.ErrAlreadyExists(fmt.Sprintf(
			"subscription '%s' already exists", msg.SubscriptionId))
	}

	// Ensure that subscription ID is not reserved
	if k.SubscriptionIdReserved(msg.SubscriptionId) {
		return nil, exported.Unauthorized(fmt.Sprintf("%s is not allowed as it is "+
			"using a reserved prefix", msg.SubscriptionId))
	}

	// Get payment contract
	contract, err := k.GetPaymentContract(ctx, msg.PaymentContractId)
	if err != nil {
		return nil, err
	}

	// Confirm that signer is actually the creator of the payment contract
	creatorAddr := exported.DidToAddr(msg.CreatorDid)
	if !creatorAddr.Equals(contract.Creator) {
		return nil, exported.InvalidAddress("signer must be payment contract creator")
	}

	// Create subscription and validate
	subscription := NewSubscription(msg.SubscriptionId,
		msg.PaymentContractId, msg.MaxPeriods, msg.Period)
	if err := subscription.Validate(); err != nil {
		return nil, err
	}

	// Submit subscription
	k.SetSubscription(ctx, subscription)

	return &sdk.Result{}, nil
}

func handleMsgGrantDiscount(ctx sdk.Context, k Keeper, msg MsgGrantDiscount) (*sdk.Result, error) {

	// Get PaymentContract
	contract, err := k.GetPaymentContract(ctx, msg.PaymentContractId)
	if err != nil {
		return nil, err
	}

	// Confirm that signer is actually the creator of the payment contract
	creatorAddr := exported.DidToAddr(msg.SenderDid)
	if !creatorAddr.Equals(contract.Creator) {
		return nil, exported.InvalidAddress("signer must be payment contract creator")
	}

	// Confirm that discount ID is in the template (to avoid invalid discount IDs)
	found, err := k.DiscountIdExists(ctx, contract.PaymentTemplateId, msg.DiscountId)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, exported.ErrInvalidId(
			"discount ID not in payment template's discount list")
	}

	// Grant the discount
	err = k.GrantDiscount(ctx, contract.Id, msg.DiscountId)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleMsgRevokeDiscount(ctx sdk.Context, k Keeper, msg MsgRevokeDiscount) (*sdk.Result, error) {

	// Get PaymentContract
	contract, err := k.GetPaymentContract(ctx, msg.PaymentContractId)
	if err != nil {
		return nil, err
	}

	// Confirm that signer is actually the creator of the payment contract
	creatorAddr := exported.DidToAddr(msg.SenderDid)
	if !creatorAddr.Equals(contract.Creator) {
		return nil, exported.InvalidAddress("signer must be payment contract creator")
	}

	// Revoke the discount
	err = k.RevokeDiscount(ctx, contract.Id)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleMsgEffectPayment(ctx sdk.Context, k Keeper, bk bank.Keeper, msg MsgEffectPayment) (*sdk.Result, error) {

	// Get payment contract
	contract, err := k.GetPaymentContract(ctx, msg.PaymentContractId)
	if err != nil {
		return nil, err
	}

	// Confirm that signer is actually the creator of the payment contract
	creatorAddr := exported.DidToAddr(msg.SenderDid)
	if !creatorAddr.Equals(contract.Creator) {
		return nil, exported.InvalidAddress("signer must be payment contract creator")
	}

	// Effect payment
	effected, err := k.EffectPayment(ctx, bk, msg.PaymentContractId)
	if err != nil {
		return nil, err
	}

	// Payment not effected but no error, meaning that payment should have been effected
	if !effected {
		return nil, exported.IntErr("payment not effected due to unknown reason")
	}

	return &sdk.Result{}, nil
}
