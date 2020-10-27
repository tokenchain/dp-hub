package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/tokenchain/dp-hub/x/dapmining/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "mine-accounts", MineAccounts(k))
}

func MineAccounts(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		// Get supply of coins held in accounts (includes stake token)
		supplyInAccounts := sdk.Coins{}
		holderAccounts := make([]sdk.AccAddress, 0)
		k.ak.IterateAccounts(ctx, func(acc exported.Account) bool {
			if acc.GetAddress().Empty() {
				return false
			}
			if acc.GetCoins().Len() == 0 {
				return false
			}
			if acc.GetCoins().Sort().AmountOf("dap").GT(sdk.ZeroInt()) || acc.GetCoins().Sort().AmountOf("dollar").GT(sdk.ZeroInt()) {
				supplyInAccounts = supplyInAccounts.Add(acc.GetCoins()...)
				holderAccounts = append(holderAccounts, acc.GetAddress())
			}
			return false
		})



		broken := count != 0
		return sdk.FormatInvariant(types.ModuleName, "mine-account", fmt.Sprintf(
			"%d Bonds supply invariants broken\n%s", count, msg)), broken
	}
}
