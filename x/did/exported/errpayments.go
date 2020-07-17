package exported

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)


func ErrNegativeSharePercentage() error {
	errMsg := fmt.Sprintf("payment distribution share percentage must be positive")
	return errors.Wrap(EInvalidShare, errMsg)
}

func ErrDistributionPercentagesNot100(total sdk.Dec) error {
	errMsg := fmt.Sprintf("payment distribution percentages should add up to 100, not %s", total.String())
	return errors.Wrap(ErrInvalidDistribution, errMsg)
}

func ErrInvalidPeriod(errMsg string) error {
	errMsg = fmt.Sprintf("period is invalid: %s", errMsg)
	return errors.Wrap(EInvalidPeriod, errMsg)
}

func ErrPaymentContractCannotBeDeauthorised() error {
	errMsg := fmt.Sprintf("payment contract cannot be deauthorised")
	return errors.Wrap(EInvalidPaymentCA, errMsg)
}

func ErrDiscountIDsBeSequentialFrom1() error {
	errMsg := fmt.Sprintf("discount IDs must be sequential starting with 1")
	return errors.Wrap(EInvalidDiscount, errMsg)
}

func ErrNegativeDiscountPercentage() error {
	errMsg := fmt.Sprintf("discount percentage must be positive")
	return errors.Wrap(EInvalidDiscount, errMsg)
}

func ErrDiscountPercentageGreaterThan100() error {
	errMsg := fmt.Sprintf("discount percentage cannot exceed 100%%")
	return errors.Wrap(EInvalidDiscount, errMsg)
}

func ErrDiscountIdIsNotInTemplate() error {
	errMsg := fmt.Sprintf("discount ID specified is not one of the template's discounts")
	return errors.Wrap(EInvalidDiscountReq, errMsg)
}

func ErrInvalidPaymentTemplate(errMsg string) error {
	errMsg = fmt.Sprintf("payment template invalid; %s", errMsg)
	return errors.Wrap(EInvalidPaymentTemplate, errMsg)
}

func ErrTriedToEffectSubscriptionPaymentWhenShouldnt() error {
	errMsg := fmt.Sprintf("tried to effect subscription payment when shouldn't")
	return errors.Wrap(EInvalidSubAction, errMsg)
}

func ErrInvalidId(errMsg string) error {
	return errors.Wrap(EInvalidId, errMsg)
}

func ErrInvalidArgument(errMsg string) error {
	return errors.Wrap(EInvalidArgs, errMsg)
}

func ErrAlreadyExists(errMsg string) error {
	return errors.Wrap(EAlreadyExists, errMsg)
}
