package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var oneHundred = sdk.NewDec(100)

type Distribution []DistributionShare

func NewDistribution(shares ...DistributionShare) Distribution {
	return Distribution(shares)
}

func (d Distribution) Validate() sdk.Error {
	// Shares must add up to 100% (no shares means 0%)
	if len(d) == 0 {
		return ErrDistributionPercentagesNot100(DefaultCodespace, sdk.ZeroDec())
	}

	// Validate shares and calculate total
	total := sdk.ZeroDec()
	for _, share := range d {
		total = total.Add(share.Percentage)
		if err := share.Validate(); err != nil {
			return err
		}
	}

	// Shares must add up to 100%
	if !total.Equal(sdk.NewDec(100)) {
		return ErrDistributionPercentagesNot100(DefaultCodespace, total)
	}

	return nil
}

func (d Distribution) GetDistributionsFor(amount sdk.Coins) []sdk.DecCoins {
	decAmount := sdk.NewDecCoins(amount)
	distributions := make([]sdk.DecCoins, len(d))

	// Calculate distribution amount for each share of the distribution
	var distributed sdk.DecCoins
	for i, share := range d {
		distributions[i] = share.GetShareOf(decAmount)
		distributed = distributed.Add(distributions[i])
	}

	// Distributed amount should equal original amount
	if !distributed.IsEqual(decAmount) {
		panic("distributing more or less than original amount")
	}

	return distributions
}

type DistributionShare struct {
	Address    sdk.AccAddress `json:"address" yaml:"address"`
	Percentage sdk.Dec        `json:"percentage" yaml:"percentage"`
}

func NewDistributionShare(address sdk.AccAddress, percentage sdk.Dec) DistributionShare {
	return DistributionShare{
		Address:    address,
		Percentage: percentage,
	}
}

func (d DistributionShare) Validate() sdk.Error {
	if !d.Percentage.IsPositive() {
		return ErrNegativeSharePercentage(DefaultCodespace)
	} else if d.Address.Empty() {
		return sdk.ErrInvalidAddress("empty distribution share address")
	}

	return nil
}

func (d DistributionShare) GetShareOf(amount sdk.DecCoins) sdk.DecCoins {
	return amount.MulDec(d.Percentage.Quo(oneHundred))
}
