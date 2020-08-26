package order

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenchain/ixo-blockchain/x/common/perf"
	"github.com/tokenchain/ixo-blockchain/x/order/keeper"
	"github.com/tokenchain/ixo-blockchain/x/order/match"
	"github.com/tokenchain/ixo-blockchain/x/order/types"
)

// EndBlocker called every block
// 1. execute matching engine
// 2. flush cache
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	seq := perf.GetPerf().OnEndBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnEndBlockExit(ctx, types.ModuleName, seq)
	match.GetEngine().Run(ctx, keeper)
	// flush cache at the end
	keeper.Cache2Disk(ctx)
	keeper.SetMetric()
	ret := keeper.GetOperationMetric()
	msg := fmt.Sprintf(
		"fullFilled<%d>, pending<%d>, canceled<%d>, expired<%d>, partialFilled<%d>",
		ret.FullFillNum,
		ret.OpenNum,
		ret.CancelNum,
		ret.ExpireNum,
		ret.PartialFillNum)

	perf.GetPerf().EnqueueMsg(msg)
}
