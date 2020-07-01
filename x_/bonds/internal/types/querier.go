package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type QueryBonds []string

func (b QueryBonds) String() string {
	return strings.Join(b[:], "\n")
}

type QueryBuyPrice struct {
	AdjustedSupply sdk.Coin  `json:"adjusted_supply" yaml:"asdjusted_supply"`
	Prices         sdk.Coins `json:"prices" yaml:"prices"`
	TxFees         sdk.Coins `json:"tx_fees" yaml:"tx_fees"`
	TotalPrices    sdk.Coins `json:"total_prices" yaml:"total_prices"`
	TotalFees      sdk.Coins `json:"total_fees" yaml:"total_fees"`
}

type QuerySellReturn struct {
	AdjustedSupply sdk.Coin  `json:"adjusted_supply" yaml:"asdjusted_supply"`
	Returns        sdk.Coins `json:"returns" yaml:"returns"`
	TxFees         sdk.Coins `json:"tx_fees" yaml:"tx_fees"`
	ExitFees       sdk.Coins `json:"exit_fees" yaml:"exit_fees"`
	TotalReturns   sdk.Coins `json:"total_returns" yaml:"total_returns"`
	TotalFees      sdk.Coins `json:"total_fees" yaml:"total_fees"`
}

type QuerySwapReturn struct {
	TotalReturns sdk.Coins `json:"total_returns" yaml:"total_returns"`
	TotalFees    sdk.Coins `json:"total_fees" yaml:"total_fees"`
}
