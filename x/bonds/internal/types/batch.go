package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
)

type Batch struct {
	BondDid         types.Did    `json:"bond_did" yaml:"bond_did"`
	BlocksRemaining sdk.Uint     `json:"blocks_remaining" yaml:"blocks_remaining"`
	TotalBuyAmount  sdk.Coin     `json:"total_buy_amount" yaml:"total_buy_amount"`
	TotalSellAmount sdk.Coin     `json:"total_sell_amount" yaml:"total_sell_amount"`
	BuyPrices       sdk.DecCoins `json:"buy_prices" yaml:"buy_prices"`
	SellPrices      sdk.DecCoins `json:"sell_prices" yaml:"sell_prices"`
	Bids            []BuyOrder   `json:"buys" yaml:"buys"`
	Asks            []SellOrder  `json:"sells" yaml:"sells"`
	Swaps           []SwapOrder  `json:"swaps" yaml:"swaps"`
}

func (b Batch) MoreBuysThanSells() bool { return b.TotalSellAmount.IsLT(b.TotalBuyAmount) }
func (b Batch) MoreSellsThanBuys() bool { return b.TotalBuyAmount.IsLT(b.TotalSellAmount) }
func (b Batch) EqualBuysAndSells() bool { return b.TotalBuyAmount.IsEqual(b.TotalSellAmount) }

func NewBatch(bondDid types.Did, token string, blocks sdk.Uint) Batch {
	return Batch{
		BondDid:         bondDid,
		BlocksRemaining: blocks,
		TotalBuyAmount:  sdk.NewInt64Coin(token, 0),
		TotalSellAmount: sdk.NewInt64Coin(token, 0),
	}
}

type BaseOrder struct {
	AccountDid   types.Did `json:"sender_did" yaml:"sender_did"`
	Amount       sdk.Coin  `json:"amount" yaml:"amount"`
	Cancelled    string    `json:"cancelled" yaml:"cancelled"`
	CancelReason string    `json:"cancel_reason" yaml:"cancel_reason"`
}

func NewBaseOrder(accountDid types.Did, amount sdk.Coin) BaseOrder {
	return BaseOrder{
		AccountDid:   accountDid,
		Amount:       amount,
		Cancelled:    FALSE,
		CancelReason: "",
	}
}

func (bo BaseOrder) IsCancelled() bool {
	return bo.Cancelled == TRUE
}

type BuyOrder struct {
	BaseOrder
	MaxPrices sdk.Coins `json:"max_prices" yaml:"max_prices"`
}

func NewBuyOrder(buyerDid types.Did, amount sdk.Coin, maxPrices sdk.Coins) BuyOrder {
	return BuyOrder{
		BaseOrder: NewBaseOrder(buyerDid, amount),
		MaxPrices: maxPrices,
	}
}

type SellOrder struct {
	BaseOrder
}

func NewSellOrder(sellerDid types.Did, amount sdk.Coin) SellOrder {
	return SellOrder{
		BaseOrder: NewBaseOrder(sellerDid, amount),
	}
}

type SwapOrder struct {
	BaseOrder
	ToToken string `json:"to_token" yaml:"to_token"`
}

func NewSwapOrder(swapperDid types.Did, from sdk.Coin, toToken string) SwapOrder {
	return SwapOrder{
		BaseOrder: NewBaseOrder(swapperDid, from),
		ToToken:   toToken,
	}
}
