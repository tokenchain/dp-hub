package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const DefaultTestFeeAmountPerBlock = "0.000001" // dap

var DefaultTestFeePerBlock = sdk.NewDecCoinFromDec(DefaultFeeDenomPerBlock, sdk.MustNewDecFromStr(DefaultTestFeeAmountPerBlock))

func DefaultTestParams() Params {
	return Params{
		OrderExpireBlocks: DefaultOrderExpireBlocks,
		MaxDealsPerBlock:  DefaultMaxDealsPerBlock,
		FeePerBlock:       DefaultTestFeePerBlock,
		TradeFeeRate:      sdk.MustNewDecFromStr(DefaultFeeRateTrade),
	}
}

// nolint
func MockOrder(orderID, product, side, price, quantity string) *Order {
	order := &Order{
		OrderID:           orderID,
		Product:           product,
		Side:              side,
		Price:             sdk.MustNewDecFromStr(price),
		FilledAvgPrice:    sdk.ZeroDec(),
		Quantity:          sdk.MustNewDecFromStr(quantity),
		RemainQuantity:    sdk.MustNewDecFromStr(quantity),
		Status:            OrderStatusOpen,
		OrderExpireBlocks: DefaultOrderExpireBlocks,
		FeePerBlock:       DefaultTestFeePerBlock,
	}
	if side == BuyOrder {
		order.RemainLocked = order.Price.Mul(order.Quantity)
	} else {
		order.RemainLocked = order.Quantity
	}
	return order
}
