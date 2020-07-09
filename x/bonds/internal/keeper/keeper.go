package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tokenchain/ixo-blockchain/x/bonds/internal/types"
)

type (
	Keeper struct {
		BankKeeper    bank.Keeper
		SupplyKeeper  supply.Keeper
		accountKeeper auth.AccountKeeper
		StakingKeeper staking.Keeper
		storeKey      sdk.StoreKey
		cdc           *codec.Codec
		paramSpace    params.Subspace
	}
)

func NewKeeper(bankKeeper bank.Keeper, supplyKeeper supply.Keeper,
	accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper,
	storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {

	// ensure batches module account is set
	if addr := supplyKeeper.GetModuleAddress(types.BatchesIntermediaryAccount); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BatchesIntermediaryAccount))
	}

	return Keeper{
		BankKeeper:    bankKeeper,
		SupplyKeeper:  supplyKeeper,
		accountKeeper: accountKeeper,
		StakingKeeper: stakingKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
	}
}
func (k Keeper) GetCodec() *codec.Codec {
	return k.cdc
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

/*
func (k Keeper) Iterator(ctx sdk.Context, cb IteratorCB) {
	kv := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(kv, []byte(valKey))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		aB := iter.Value()
		var asset types.Asset
		k.cdc.MustUnmarshalBinaryBare(aB, &asset)

		if !cb(asset) {
			break
		}
	}
}
*/
