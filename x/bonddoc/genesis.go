package bonddoc

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	println(fmt.Sprintf("===BondDocs=== %v", data))
	// Initialise bond docs
	println("BondDocs ===")
	if data.BondDocs != nil {
		println("BondDocs ok!")
		for _, b := range data.BondDocs {
			keeper.SetBondDoc(ctx, &b)
		}
	}
	println(fmt.Sprintf("===BondDocs done=== %v", data))
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// Export bond docs
	var bondDocs []MsgCreateBond

	iterator := k.GetBondDocIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		bondDoc := k.MustGetBondDocByKey(ctx, iterator.Key())

		bondDocs = append(bondDocs, *bondDoc.(*MsgCreateBond))
	}

	return GenesisState{
		BondDocs: bondDocs,
	}
}
