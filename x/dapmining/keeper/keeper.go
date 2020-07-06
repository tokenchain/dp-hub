package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type (
	Keeper struct {
		bk           bank.Keeper
		SupplyKeeper supply.Keeper
		ak           auth.AccountKeeper
		sk           staking.Keeper
		mk           mint.Keeper
		storeKey     sdk.StoreKey

		cdc *codec.Codec
	}
)
