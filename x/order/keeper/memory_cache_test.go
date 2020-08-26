package keeper

import (
	"testing"
	"time"

	"github.com/tokenchain/ixo-blockchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestCache_GetCancelNum(t *testing.T) {
	cache := NewCache()
	cache.addUpdatedOrderID("ID0000000010-1")
	require.EqualValues(t, 1, len(cache.updatedOrderIDs))

	cache.addUpdatedOrderID("ID0000000010-2")
	require.EqualValues(t, 2, len(cache.getUpdatedOrderIDs()))

	cache.IncreaseCancelNum()
	cache.IncreaseExpireNum()
	cache.IncreaseFullFillNum()
	cache.IncreasePartialFillNum()

	cache.DecreaseCancelNum()
	cache.DecreaseFullFillNum()
	cache.DecreasePartialFillNum()

	require.EqualValues(t, 0, cache.GetCancelNum())
	require.EqualValues(t, 1, cache.GetExpireNum())
	require.EqualValues(t, 0, cache.GetFullFillNum())
	require.EqualValues(t, 0, cache.GetPartialFillNum())

	feeParams := types.DefaultTestParams()
	cache.SetParams(&feeParams)

	require.EqualValues(t, types.DefaultOrderExpireBlocks, cache.params.OrderExpireBlocks)

	res := types.BlockMatchResult{
		BlockHeight: 0,
		ResultMap:   nil,
		TimeStamp:   time.Now().Unix(),
	}
	cache.setBlockMatchResult(&res)

	require.NotEqual(t, 0, cache.getBlockMatchResult().TimeStamp)

	cache.reset()
	require.Nil(t, cache.GetParams())

}
