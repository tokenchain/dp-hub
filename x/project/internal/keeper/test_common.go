package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmDB "github.com/tendermint/tm-db"

	"github.com/tokenchain/ixo-blockchain/x/payments"
	"github.com/tokenchain/ixo-blockchain/x/project/internal/types"
)

func CreateTestInput() (sdk.Context, Keeper, *codec.Codec,
	payments.Keeper, bank.Keeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	actStoreKey := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey("subspace")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")
	keyFees := sdk.NewKVStoreKey(payments.StoreKey)

	db := tmDB.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(actStoreKey, sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(keyFees, sdk.StoreTypeIAVL, nil)
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, true, log.NewNopLogger())
	cdc := MakeTestCodec()

	pk1 := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(
		cdc, actStoreKey, pk1.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)

	paymentsSubspace := pk1.Subspace(payments.DefaultParamspace)
	projectSubspace := pk1.Subspace(types.DefaultParamspace)

	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk1.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)
	paymentsKeeper := payments.NewKeeper(cdc, keyFees, paymentsSubspace, bankKeeper, nil)
	keeper := NewKeeper(cdc, storeKey, projectSubspace, accountKeeper, paymentsKeeper)

	paymentsKeeper.SetParams(ctx, payments.DefaultParams())

	return ctx, keeper, cdc, paymentsKeeper, bankKeeper
}

func MakeTestCodec() *codec.Codec {
	return codec.New()
}
