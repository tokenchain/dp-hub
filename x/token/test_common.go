package token

import (
	"github.com/tokenchain/ixo-blockchain/x/token/keeper"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/tokenchain/ixo-blockchain/x/params"
	"github.com/tokenchain/ixo-blockchain/x/token/types"
)

// CreateParam create okchain parm for test
func CreateParam(t *testing.T, isCheckTx bool) (sdk.Context, keeper.Keeper, *sdk.KVStoreKey, []byte) {
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	keyToken := sdk.NewKVStoreKey("token")
	keyLock := sdk.NewKVStoreKey("lock")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLock, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, nil)

	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Keeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)
	//feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	blacklistedAddrs := make(map[string]bool)
	//blacklistedAddrs[feeCollectorAcc.String()] = true

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Keeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)
	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.ModuleName:      nil,
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)
	tk := keeper.NewKeeper(bk,
		pk.Keeper.Subspace(DefaultParamspace),
		auth.FeeCollectorName,
		supplyKeeper,
		keyToken,
		keyLock,
		cdc,
		true)
	tk.SetParams(ctx, types.DefaultParams())

	return ctx, tk, keyParams, []byte("testToken")
}
