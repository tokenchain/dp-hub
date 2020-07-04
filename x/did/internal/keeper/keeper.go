package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	er "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	fn "github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

func (k Keeper) GetDidDoc(ctx sdk.Context, didc did.Did) (did.DidDoc, error) {
	store := ctx.KVStore(k.storeKey)
	key := fn.GetDidPrefixKey(didc)
	bz := store.Get(key)
	if bz == nil {
		return nil, er.Wrap(x.ErrorInvalidDidE, "Invalid Did Address")
	}

	var didDoc did.BaseDidDoc
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &didDoc)

	return didDoc, nil
}

func (k Keeper) SetDidDoc(ctx sdk.Context, did did.DidDoc) (err error) {
	existedDidDoc, err := k.GetDidDoc(ctx, did.GetDid())
	if existedDidDoc != nil {
		return er.Wrap(x.ErrorInvalidDidE, "Did already exists")
	}

	k.AddDidDoc(ctx, did)
	return nil
}

func (k Keeper) AddDidDoc(ctx sdk.Context, did did.DidDoc) {
	store := ctx.KVStore(k.storeKey)
	key := fn.GetDidPrefixKey(did.GetDid())
	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(did))
}

func (k Keeper) AddCredentials(ctx sdk.Context, didc did.Did, credential did.DidCredential) (err error) {
	existedDid, err := k.GetDidDoc(ctx, didc)
	if err != nil {
		return err
	}

	baseDidDoc := existedDid.(did.BaseDidDoc)
	credentials := baseDidDoc.GetCredentials()

	for _, data := range credentials {
		if data.Issuer == credential.Issuer && data.CredType[0] == credential.CredType[0] && data.CredType[1] == credential.CredType[1] && data.Claim.KYCValidated == credential.Claim.KYCValidated {
			return er.Wrap(x.ErrorInvalidCredentials, "credentials already exist")
		}
	}

	baseDidDoc.AddCredential(credential)
	k.AddDidDoc(ctx, baseDidDoc)

	return nil
}

func (k Keeper) GetAllDidDocs(ctx sdk.Context) (didDocs []did.DidDoc) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, did.DidKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var didDoc did.BaseDidDoc
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &didDoc)
		didDocs = append(didDocs, &didDoc)
	}

	return didDocs
}

func (k Keeper) GetAddDids(ctx sdk.Context) (dids []did.Did) {
	didDocs := k.GetAllDidDocs(ctx)
	for _, did := range didDocs {
		dids = append(dids, did.GetDid())
	}

	return dids
}
