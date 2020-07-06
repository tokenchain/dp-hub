package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	r "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/types"
	types2 "github.com/tokenchain/ixo-blockchain/x/dap/types"
)

type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

func (k Keeper) GetBondDocIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.BondKey)
}

func (k Keeper) MustGetBondDocByKey(ctx sdk.Context, key []byte) types.StoredBondDoc {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		panic("bond doc not found")
	}

	bz := store.Get(key)
	var bondDoc types.MsgCreateBond
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &bondDoc)

	return &bondDoc
}

func (k Keeper) BondDocExists(ctx sdk.Context, bondDid types2.Did) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBondPrefixKey(bondDid))
}

func (k Keeper) GetBondDoc(ctx sdk.Context, bondDid types2.Did) (types.StoredBondDoc, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBondPrefixKey(bondDid)

	bz := store.Get(key)
	if bz == nil {
		return nil,
			r.Wrap(x.ErrorInvalidDidE, "Invalid BondDid Address")
	}

	var bondDoc types.MsgCreateBond
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &bondDoc)

	return &bondDoc, nil
}

func (k Keeper) SetBondDoc(ctx sdk.Context, bondDoc types.StoredBondDoc) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBondPrefixKey(bondDoc.GetBondDid())
	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(bondDoc))
}

func (k Keeper) UpdateBondDoc(ctx sdk.Context, newBondDoc types.StoredBondDoc) (types.StoredBondDoc, error) {
	existedDoc, _ := k.GetBondDoc(ctx, newBondDoc.GetBondDid())
	if existedDoc == nil {
		return nil,
			r.Wrap(x.ErrorInvalidDidE, "BondDoc details are not exist")
	} else {

		existedDoc.SetStatus(newBondDoc.GetStatus())
		k.SetBondDoc(ctx, newBondDoc)

		return newBondDoc, nil
	}
}
