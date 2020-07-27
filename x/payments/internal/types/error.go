package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

func ErrNegativeSharePercentage() error {
	errMsg := fmt.Sprintf("payment distribution share percentage must be positive")
	return errors.Wrap(exported.EInvalidShare, errMsg)
}

func ErrDistributionPercentagesNot100(total sdk.Dec) error {
	errMsg := fmt.Sprintf("payment distribution percentages should add up to 100, not %s", total.String())
	return errors.Wrap(exported.ErrInvalidDistribution, errMsg)
}

func ErrInvalidPeriod(errMsg string) error {
	errMsg = fmt.Sprintf("period is invalid: %s", errMsg)
	return errors.Wrap(exported.EInvalidPeriod, errMsg)
}

func ErrPaymentContractCannotBeDeauthorised() error {
	errMsg := fmt.Sprintf("payment contract cannot be deauthorised")
	return errors.Wrap(exported.EInvalidPaymentCA, errMsg)
}

func ErrDiscountIDsBeSequentialFrom1() error {
	errMsg := fmt.Sprintf("discount IDs must be sequential starting with 1")
	return errors.Wrap(exported.EInvalidDiscount, errMsg)
}

func ErrNegativeDiscountPercentage() error {
	errMsg := fmt.Sprintf("discount percentage must be positive")
	return errors.Wrap(exported.EInvalidDiscount, errMsg)
}

func ErrDiscountPercentageGreaterThan100() error {
	errMsg := fmt.Sprintf("discount percentage cannot exceed 100%%")
	return errors.Wrap(exported.EInvalidDiscount, errMsg)
}

func ErrDiscountIdIsNotInTemplate() error {
	errMsg := fmt.Sprintf("discount ID specified is not one of the template's discounts")
	return errors.Wrap(exported.EInvalidDiscountReq, errMsg)
}

func ErrInvalidPaymentTemplate(errMsg string) error {
	errMsg = fmt.Sprintf("payment template invalid; %s", errMsg)
	return errors.Wrap(exported.EInvalidPaymentTemplate, errMsg)
}

func ErrTriedToEffectSubscriptionPaymentWhenShouldnt() error {
	errMsg := fmt.Sprintf("tried to effect subscription payment when shouldn't")
	return errors.Wrap(exported.EInvalidSubAction, errMsg)
}

func ErrInvalidId(errMsg string) error {
	return errors.Wrap(exported.EInvalidId, errMsg)
}

func ErrInvalidArgument(errMsg string) error {
	return errors.Wrap(exported.EInvalidArgs, errMsg)
}

func ErrAlreadyExists(errMsg string) error {
	return errors.Wrap(exported.EAlreadyExists, errMsg)
}
