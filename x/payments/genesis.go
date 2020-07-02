package payments

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis new payments genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {

	// Init payment templates
	for _, pt := range data.PaymentTemplates {
		keeper.SetPaymentTemplate(ctx, pt)
	}
	//print("===3===")
	// Init payment contracts
	for _, pc := range data.PaymentContracts {
		keeper.SetPaymentContract(ctx, pc)
	}
	//print("===2===")
	// Init subscriptions
	for _, s := range data.Subscriptions {
		keeper.SetSubscription(ctx, s)
	}
	// Init params
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	// Export params
	params := keeper.GetParams(ctx)
	// Export payment templates
	var templates []PaymentTemplate
	iterator := keeper.GetPaymentTemplateIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		template := keeper.MustGetPaymentTemplateByKey(ctx, iterator.Key())
		templates = append(templates, template)
	}
	// Export payment contracts
	var contracts []PaymentContract
	iterator = keeper.GetPaymentContractIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		contract := keeper.MustGetPaymentContractByKey(ctx, iterator.Key())
		contracts = append(contracts, contract)
	}
	// Export subscriptions
	var subscriptions []Subscription
	iterator = keeper.GetSubscriptionIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		subscription := keeper.MustGetSubscriptionByKey(ctx, iterator.Key())
		subscriptions = append(subscriptions, subscription)
	}
	return NewGenesisState(params, templates, contracts, subscriptions)
}
