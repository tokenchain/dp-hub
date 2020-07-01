package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/x"
)

const (
	QueryBondDoc = "queryBondDoc"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req types.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryBondDoc:
			return queryBondDoc(ctx, path[1:], k)
		default:
			return nil, x.UnknownRequest("Unknown bond query endpoint")
		}
	}
}

func queryBondDoc(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	storedDoc, err := k.GetBondDoc(ctx, path[0])
	if err != nil {
		return nil, err
	}

	res, errRes := codec.MarshalJSONIndent(k.cdc, storedDoc)
	if errRes != nil {

		return nil,
			x.ErrJsonMars(errRes.Error())
	}

	return res, nil
}
