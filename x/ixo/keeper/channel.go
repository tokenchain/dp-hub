package keeper
/*
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/dp-hub/app"
)

func (k Keeper) SetChannel(ctx sdk.Context, chainName string, port string, channel string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChannelStoreKey(chainName, port), []byte(channel))
}

func (k Keeper) GetChannel(ctx sdk.Context, chainName string, port string) (string, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.HasChannel(ctx, chainName, port) {
		return "", app.UnknownRequest("channel not found")
	}
	return string(store.Get(types.ChannelStoreKey(chainName, port))), nil
}

func (k Keeper) HasChannel(ctx sdk.Context, chainName string, port string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ChannelStoreKey(chainName, port))
}
*/