package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

const (
	QueryParams          = "queryParams"
	QueryPaymentTemplate = "queryPaymentTemplate"
	QueryPaymentContract = "queryPaymentContract"
	QuerySubscription    = "querySubscription"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, k)
		case QueryPaymentTemplate:
			return queryPaymentTemplate(ctx, path[1:], k)
		case QueryPaymentContract:
			return queryPaymentContract(ctx, path[1:], k)
		case QuerySubscription:
			return querySubscription(ctx, path[1:], k)
		default:
			return nil, exported.UnknownRequest("unknown payments query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil,  exported.ErrJsonMars(err.Error())
	}

	return res, nil
}

func queryPaymentTemplate(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	templateId := path[0]

	template, err := k.GetPaymentTemplate(ctx, templateId)
	if err != nil {
		return nil, exported.UnknownRequest(fmt.Sprintf("payment template '%s' does not exist", templateId))
	}

	res, err2 := codec.MarshalJSONIndent(k.cdc, template)
	if err2 != nil {
		return nil,  exported.ErrJsonMars(err2.Error())
	}

	return res, nil
}

func queryPaymentContract(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	contractId := path[0]

	contract, err := k.GetPaymentContract(ctx, contractId)
	if err != nil {
		return nil, exported.UnknownRequest(fmt.Sprintf("payment contract '%s' does not exist", contractId))
	}

	res, err2 := codec.MarshalJSONIndent(k.cdc, contract)
	if err2 != nil {
		return nil,  exported.ErrJsonMars(err2.Error())
	}

	return res, nil
}

func querySubscription(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	subscriptionId := path[0]

	subscription, err := k.GetSubscription(ctx, subscriptionId)
	if err != nil {
		return nil, exported.UnknownRequest(fmt.Sprintf("subscription '%s' does not exist", subscriptionId))
	}

	res, err2 := codec.MarshalJSONIndent(k.cdc, subscription)
	if err2 != nil {
		return nil,  exported.ErrJsonMars(err2.Error())
	}

	return res, nil
}
