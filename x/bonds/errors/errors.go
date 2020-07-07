package errors

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

var (
	ErrArgument                             = errors.Register(ModuleName, CodeArgumentInvalid, "Cannot be empty")
	ErrArgumentMissingOrIncorrectType       = errors.Register(ModuleName, CodeArgumentMissingOrIncorrectType, "Missing or Incorrect Type")
	ErrCodeIncorrectNumberOfValues          = errors.Register(ModuleName, CodeIncorrectNumberOfValues, "Incorrect code number of value")
	ErrUnrecognizedFunctionType             = errors.Register(ModuleName, CodeUnrecognizedFunctionType, "Unrecognized function type")
	ErrCodeInvalidFuncParam                 = errors.Register(ModuleName, CodeInvalidFunctionParameter, "Invalid Function Parameter")
	ErrFunctionRequiresNonZeroCurrentSupply = errors.Register(ModuleName, CodeFunctionRequiresNonZeroCurrentSupply, "Function requires the current supply to be non zero")
	ErrTokenIsNotAValidReserveTokenCode     = errors.Register(ModuleName, CodeReserveTokenInvalid, "Function requires the current supply to be non zero")
	ErrMaxSupplyDenomDoesNotMatchTokenDenom = errors.Register(ModuleName, CodeMaxSupplyDenomInvalid, "Max supply denom does not match token denom")
	ErrBondInvalidToken                     = errors.Register(ModuleName, CodeBondTokenInvalid, "bond token is invalid")
	ErrReserveDenomsMismatchE               = errors.Register(ModuleName, CodeReserveDenomsMismatch, "reserve denom mismatch")
	ErroInvalidCoinDenomination             = errors.Register(ModuleName, CodeInvalidCoinDenomination, "wrong coin denomination")
	EInvalidResultantSupply                 = errors.Register(ModuleName, CodeInvalidResultantSupply, "Invalid resultant supply")
	EPriceExceed                            = errors.Register(ModuleName, CodeMaxPriceExceeded, "price exceeded")
	ESwapAmountInvalid                      = errors.Register(ModuleName, CodeSwapAmountInvalid, "invalid amount in swap")
	ErrOrderQuantityLimitExceeded           = errors.Register(ModuleName, CodeOrderLimitExceeded, "Order quantity limits exceeded")
	ErrValuesViolateSanityRate              = errors.Register(ModuleName, CodeSanityRateViolated, "Values violate sanity rate")
	ErrFeesCannotBeOrExceed100Percent       = errors.Register(ModuleName, CodeFeeTooLarge, "Sum of fees is or exceeds 100 percent")
	ErrCodeBondDoesNotExist                 = errors.Register(ModuleName, CodeBondDoesNotExist, "Code bond does not exist")
	ErrCodeBondAlreadyExists                = errors.Register(ModuleName, CodeBondAlreadyExists, "Code bond already exist")
	ErrCodeBondDoesNotAllowSelling          = errors.Register(ModuleName, CodeBondDoesNotAllowSelling, "Code bond does not allow selling")
	ErrCodeDidNotEditAnything               = errors.Register(ModuleName, CodeDidNotEditAnything, "Did not edit anything from the bond.")
	ErrFromAndToCannotBeTheSameToken_E      = errors.Register(ModuleName, CodeInvalidSwapper, "From and To tokens cannot be the same token.")
	ErrDuplicateReserveToken                = errors.Register(ModuleName, CodeInvalidBond, "Cannot have duplicate tokens in reserve tokens.")
	ErrFunctionNotAvailableForFunctionType  = errors.Register(ModuleName, CodeFunctionNotAvailableForFunctionType, "Function is not available for the function type")
)

func ErrBondDoesNotExist(bondDid string) error {
	return errors.Wrapf(ErrCodeBondDoesNotExist, "Bond '%s' does not exist", bondDid)
}
func ErrBondAlreadyExists(bonddid string) error {
	return errors.Wrapf(ErrCodeBondAlreadyExists, "Bond '%s' already exists", bonddid)
}

func ErrFromAndToCannotBeTheSameToken() error {
	return errors.Wrap(ErrFromAndToCannotBeTheSameToken_E, "From and To tokens cannot be the same token.")
}
func ErrBondTokenIsTaken(bondToken string) error {
	return errors.Wrapf(ErrCodeBondAlreadyExists, "Bond token '%s' is taken", bondToken)
}
func ErrBondDoesNotAllowSelling() error {
	return errors.Wrap(ErrCodeBondDoesNotAllowSelling, "Bond does not allow selling.")
}
func ErrDidNotEditAnything() error {
	return errors.Wrap(ErrCodeDidNotEditAnything, "Did not edit anything from the bond.")
}
func FunctionNotAvailableForFunctionType() error {
	return errors.Wrap(ErrFunctionNotAvailableForFunctionType, "Function is not available for the function type.")
}
func FunctionRequiresNonZeroCurrentSupply() error {
	return errors.Wrap(ErrFunctionRequiresNonZeroCurrentSupply, "Function requires the current supply to be non zero. ")
}
func ErrIncorrectNumberOfReserveTokens(expected int) error {
	return errors.Wrapf(ErrCodeIncorrectNumberOfValues, "Incorrect number of reserve tokens; expected: %d", expected)
}
func UnrecognizedFunctionType() error {
	return errors.Wrapf(ErrUnrecognizedFunctionType, "Cannot recognize the function type")
}
func BondTokenCannotAlsoBeReserveToken() error {
	errMsg := "Token cannot also be a reserve token"
	return errors.Wrap(ErrBondInvalidToken, errMsg)
}
func InvalidCoinDenomination(denom string) error {
	return errors.Wrapf(ErroInvalidCoinDenomination, "Sum of fees is or exceeds 100 percent")
}
func DuplicateReserveToken() error {
	errMsg := "Token Reserve is duplicated"
	return errors.Wrap(ErrDuplicateReserveToken, errMsg)
}
func ArgumentCannotBeNegative(arg string) error {
	return errors.Wrapf(ErrArgument, "%s argument cannot be negative", arg)
}
func ArgumentMustBePositive(arg string) error {
	return errors.Wrapf(ErrArgument, "%s argument must be a positive value", arg)
}
func FeesCannotBeOrExceed100Percent() error {
	return errors.Wrapf(ErrFeesCannotBeOrExceed100Percent, "")
}

func FunctionParameterMissingOrNonInteger(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s parameter is missing or is not an integer", arg)
}
func IncorrectNumberOfFunctionParameters(expected int) error {
	return errors.Wrapf(ErrCodeIncorrectNumberOfValues, "Incorrect number of function parameters; expected: %d", expected)
}
func ArgumentCannotBeEmpty(argument string) error {
	return errors.Wrapf(ErrArgument, "%s argument cannot be empty", argument)
}
func MaxSupplyDenomDoesNotMatchTokenDenom() error {
	return errors.Wrap(ErrMaxSupplyDenomDoesNotMatchTokenDenom, "Max supply denom does not match token denom")
}
func ArgumentMissingOrNonFloat(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s argument is missing or is not a float", arg)
}
func ArgumentMissingOrNonInteger(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s argument is missing or is not an integer", arg)
}
func ArgumentMissingOrNonUInteger(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s argument is missing or is not an unsigned integer", arg)
}
func ArgumentMissingOrNonBoolean(arg string) error {
	return errors.Wrapf(ErrCodeIncorrectNumberOfValues, "%s argument is missing or is not true or false", arg)
}
func InvalidFunctionParameter(parameter string) error {
	return errors.Wrapf(ErrCodeInvalidFuncParam, "Invalid function parameter '%s'", parameter)
}
func TokenIsNotAValidReserveToken(denom string) error {
	return errors.Wrapf(ErrTokenIsNotAValidReserveTokenCode, "Token '%s' is not a valid reserve token", denom)
}
func BondTokenCannotBeStakingToken() error {
	errMsg := "Bond token cannot be staking token"
	return errors.Wrap(ErrBondInvalidToken, errMsg)
}
func BondTokenDoesNotMatchBond() error {
	errMsg := "Bond token does not match bond"
	return errors.Wrap(ErrBondInvalidToken, errMsg)
}
func ReserveDenomsMismatch(inputDenoms string, actualDenoms []string) error {
	return errors.Wrapf(ErrReserveDenomsMismatchE, "Denoms in %s do not match reserve denoms; expected: %s", inputDenoms, strings.Join(actualDenoms, ","))
}
func CannotMintMoreThanMaxSupply() error {
	errMsg := "Cannot mint more tokens than the max supply"
	return errors.Wrap(EInvalidResultantSupply, errMsg)
}
func CannotBurnMoreThanSupply() error {
	errMsg := "Cannot burn more tokens than the current supply"
	return errors.Wrap(EInvalidResultantSupply, errMsg)
}
func MaxPriceExceeded(totalPrice, maxPrice sdk.Coins) error {
	return errors.Wrapf(EPriceExceed, "Actual prices %s exceed max prices %s", totalPrice.String(), maxPrice.String())
}
func SwapAmountTooSmallToGiveAnyReturn(fromToken, toToken string) error {
	return errors.Wrapf(ESwapAmountInvalid, "%s swap amount too small to give any %s return", fromToken, toToken)
}
func SwapAmountCausesReserveDepletion(fromToken, toToken string) error {
	return errors.Wrapf(ESwapAmountInvalid, "%s swap amount too large and causes %s reserve to be depleted", fromToken, toToken)
}

func OrderQuantityLimitExceeded() error {
	return errors.Wrap(ErrOrderQuantityLimitExceeded, "Order quantity limits exceeded")
}

func ValuesViolateSanityRate() error {
	return errors.Wrap(ErrValuesViolateSanityRate, "liquidity violates sanity rate")
}

func InternalErr(m string) error {
	return errors.Wrap(errors.ErrPanic, m)
}
func Unauthorizedf(statement string, args ...interface{}) error {
	return errors.Wrap(errors.ErrUnauthorized, fmt.Sprintf(statement, args...))
}
