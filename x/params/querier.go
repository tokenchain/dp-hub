package params

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"

	"github.com/tokenchain/ixo-blockchain/x/params/types"
)

// NewQuerier returns all query handlers
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, did.UnknownRequest("unknown params query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, did.IntErr(fmt.Sprintf("could not marshal result to JSON %s", err.Error()))
	}
	return bz, nil
}
