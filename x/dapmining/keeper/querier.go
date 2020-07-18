package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"time"
)

const (
	QueryRewardHistory = "history"
)

type QueryRewardHistoryReturn struct {
	RewardTime           time.Time `json:"snapshot_date" yaml:"snapshot_date"`
	TotalMinedCoins      sdk.Coins `json:"total_returns" yaml:"total_returns"`
	TotalRewardsAccounts uint64    `json:"total_rewards" yaml:"total_rewards"`
}

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryRewardHistory:
			return qHistory(ctx, keeper)
		default:
			return nil, exported.UnknownRequest("unknown rewards query endpoint")
		}
	}
}
func qHistory(ctx sdk.Context, k Keeper) (res []byte, err error) {
	var result QueryRewardHistoryReturn
	bz, err2 := codec.MarshalJSONIndent(k.cdc, result)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
