package keeper

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	er "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/params"

	exportedDid "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/payments"

	"github.com/tokenchain/ixo-blockchain/x/project/internal/types"
)

type Keeper struct {
	cdc            *codec.Codec
	storeKey       sdk.StoreKey
	paramSpace     params.Subspace
	AccountKeeper  auth.AccountKeeper
	paymentsKeeper payments.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	accountKeeper auth.AccountKeeper, paymentsKeeper payments.Keeper) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       key,
		paramSpace:     paramSpace.WithKeyTable(types.ParamKeyTable()),
		AccountKeeper:  accountKeeper,
		paymentsKeeper: paymentsKeeper,
	}
}

// GetParams returns the total set of project parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of project parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) GetProjectDocIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ProjectKey)
}

func (k Keeper) MustGetProjectDocByKey(ctx sdk.Context, key []byte) types.StoredProjectDoc {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		panic("project doc not found")
	}

	bz := store.Get(key)
	var projectDoc types.MsgCreateProject
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &projectDoc)

	return &projectDoc
}

func (k Keeper) ProjectDocExists(ctx sdk.Context, projectDid exportedDid.Did) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetProjectPrefixKey(projectDid))
}

func (k Keeper) GetProjectDoc(ctx sdk.Context, projectDid exportedDid.Did) (types.StoredProjectDoc, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetProjectPrefixKey(projectDid)

	bz := store.Get(key)
	if bz == nil {
		return nil, exportedDid.Invalid("Invalid ProjectDid Address")
	}

	var projectDoc types.MsgCreateProject
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &projectDoc)

	return &projectDoc, nil
}

func (k Keeper) SetProjectDoc(ctx sdk.Context, projectDoc types.StoredProjectDoc) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetProjectPrefixKey(projectDoc.GetProjectDid())
	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(projectDoc))
}

func (k Keeper) UpdateProjectDoc(ctx sdk.Context, newProjectDoc types.StoredProjectDoc) (types.StoredProjectDoc, error) {
	existedDoc, _ := k.GetProjectDoc(ctx, newProjectDoc.GetProjectDid())
	if existedDoc == nil {

		return nil, exportedDid.Invalid("ProjectDoc details are not exist")
	} else {

		existedDoc.SetStatus(newProjectDoc.GetStatus())
		k.SetProjectDoc(ctx, newProjectDoc)

		return newProjectDoc, nil
	}
}

func (k Keeper) SetAccountMap(ctx sdk.Context, projectDid exportedDid.Did, accountMap types.AccountMap) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(accountMap)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetAccountPrefixKey(projectDid), bz)
}

func (k Keeper) GetAccountMap(ctx sdk.Context, projectDid exportedDid.Did) types.AccountMap {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAccountPrefixKey(projectDid)

	bz := store.Get(key)
	if bz == nil {
		return make(types.AccountMap)
	} else {
		var accountMap types.AccountMap
		if err := json.Unmarshal(bz, &accountMap); err != nil {
			panic(err)
		}

		return accountMap
	}
}

func (k Keeper) AddAccountToProjectAccounts(ctx sdk.Context, projectDid exportedDid.Did,
	accountId types.InternalAccountID, account exported.Account) {
	accountMap := k.GetAccountMap(ctx, projectDid)
	_, found := accountMap[accountId]
	if found {
		return
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetAccountPrefixKey(projectDid)
	accountMap[accountId] = account.GetAddress()

	bz, err := json.Marshal(accountMap)
	if err != nil {
		panic(err)
	}

	store.Set(key, bz)
}

func (k Keeper) CreateNewAccount(ctx sdk.Context, projectDid exportedDid.Did,
	accountId types.InternalAccountID) (exported.Account, error) {
	address := exportedDid.StringToAddr(accountId.ToAddressKey(projectDid))

	if k.AccountKeeper.GetAccount(ctx, address) != nil {
		return nil, er.Wrap(er.ErrInvalidAddress, "Generate account already exists")
	}

	account := k.AccountKeeper.NewAccountWithAddress(ctx, address)
	k.AccountKeeper.SetAccount(ctx, account)

	return account, nil
}

func (k Keeper) SetProjectWithdrawalTransactions(ctx sdk.Context, projectDid exportedDid.Did, txs []types.WithdrawalInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(txs)
	store.Set(types.GetWithdrawalPrefixKey(projectDid), bz)
}

func (k Keeper) GetProjectWithdrawalTransactions(ctx sdk.Context, projectDid exportedDid.Did) ([]types.WithdrawalInfo, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetWithdrawalPrefixKey(projectDid)

	bz := store.Get(key)
	if bz == nil {
		return []types.WithdrawalInfo{}, exportedDid.Invalid("ProjectDoc doesn't exist")
	} else {
		var txs []types.WithdrawalInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &txs)

		return txs, nil
	}
}

func (k Keeper) AddProjectWithdrawalTransaction(ctx sdk.Context, projectDid exportedDid.Did, info types.WithdrawalInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetWithdrawalPrefixKey(projectDid)

	txs, _ := k.GetProjectWithdrawalTransactions(ctx, projectDid)
	txs = append(txs, info)

	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(txs))
}
