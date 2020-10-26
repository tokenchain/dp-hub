package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/dp-block/x/did/exported"
)

const (
	QueryDidDoc     = "queryDidDoc"
	QueryAllDids    = "queryAllDids"
	QueryAllDidDocs = "queryAllDidDocs"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryDidDoc:
			return queryDidDoc(ctx, path[1:], k)
		case QueryAllDids:
			return queryAllDids(ctx, k)
		case QueryAllDidDocs:
			return queryAllDidDocs(ctx, k)
		default:
			return nil, exported.UnknownRequest("Unknown did query endpoint")
		}
	}
}

func queryDidDoc(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	didDoc, err := k.GetDidDoc(ctx, path[0])
	if err != nil {
		return nil, err
	}

	res, errRes := codec.MarshalJSONIndent(k.cdc, didDoc)
	if errRes != nil {
		return nil, exported.IntErr(fmt.Sprintf("failed to marshal data %s", errRes))
	}

	return res, nil
}

func queryAllDids(ctx sdk.Context, k Keeper) ([]byte, error) {
	allDids := k.GetAddDids(ctx)

	res, errRes := json.Marshal(allDids)
	if errRes != nil {
		return nil, exported.IntErr(fmt.Sprintf("failed to marshal data %s", errRes.Error()))
	}

	return res, nil
}

func queryAllDidDocs(ctx sdk.Context, k Keeper) ([]byte, error) {
	var didDocs []exported.DidDoc
	didDocs = k.GetAllDidDocs(ctx)

	res, errRes := json.Marshal(didDocs)
	if errRes != nil {
		return nil, exported.IntErr(fmt.Sprintf("failed to marshal data %s", errRes.Error()))
	}

	return res, nil
}
