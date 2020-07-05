package bonds

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/bonds/internal/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Initialise bonds
	println("Bonds ===")
	if data.Bonds != nil {
		println("Bonds ok=")
		for _, b := range data.Bonds {
			keeper.SetBond(ctx, b.BondDid, b)
			keeper.SetBondDid(ctx, b.Token, b.BondDid)
		}
	}
	println("Batches ===")
	if data.Batches != nil {
		println("Batches ok=")
		// Initialise batches
		for _, b := range data.Batches {
			keeper.SetBatch(ctx, b.BondDid, b)
		}
	}
}




func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// Export bonds and batches
	var bonds []types.Bond
	var batches []types.Batch

	iterator := k.GetBondIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		bond := k.MustGetBondByKey(ctx, iterator.Key())
		batch := k.MustGetBatch(ctx, bond.BondDid)

		bonds = append(bonds, bond)
		batches = append(batches, batch)
	}

	return GenesisState{
		Bonds:   bonds,
		Batches: batches,
	}
}
