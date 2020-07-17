package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

// --------------------------------------------- Discounts

type Discounts []Discount

func NewDiscounts(discounts ...Discount) Discounts {
	return Discounts(discounts)
}

func (ds Discounts) Validate() error {
	if len(ds) == 0 {
		return nil
	}

	// Check that discount IDs are sequential, STRICTLY starting with 1,
	// since in a payment contract zero indicates the lack of discount
	id := sdk.OneUint()
	for _, d := range ds {
		if !d.Id.Equal(id) {
			return exported.ErrDiscountIDsBeSequentialFrom1()
		}
		id = id.Add(sdk.OneUint())
	}

	// Validate list of discounts
	for _, d := range ds {
		if err := d.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Discount struct {
	Id      sdk.Uint `json:"id" yaml:"id"`
	Percent sdk.Dec  `json:"percent" yaml:"percent"`
}

func NewDiscount(id sdk.Uint, percent sdk.Dec) Discount {
	return Discount{
		Id:      id,
		Percent: percent,
	}
}

func (d Discount) Validate() error {
	if !d.Percent.IsPositive() {
		return exported.ErrNegativeDiscountPercentage()
	} else if d.Percent.GT(sdk.NewDec(100)) {
		return exported.ErrDiscountPercentageGreaterThan100()
	}

	return nil
}
