package x

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

// Local code type
type CodeType = uint32
type CodespaceType = string

const (
	moduleNameDid     = "did"
	moduleNameBonddoc = "bonddoc"
	moduleNameBonds   = "bonds"
	moduleNameIxo     = "dap"
	moduleNamePayment = "payments"
	moduleNameProject = "project"

	CodeInvalidDid         CodeType = 201
	CodeInvalidPubKey      CodeType = 202
	CodeInvalidIssuer      CodeType = 203
	CodeInvalidCredentials CodeType = 204

	// General
	CodeArgumentInvalid                CodeType = 301
	CodeArgumentMissingOrIncorrectType CodeType = 302
	CodeIncorrectNumberOfValues        CodeType = 303

	// Function types and function parameters
	CodeUnrecognizedFunctionType             CodeType = 310
	CodeInvalidFunctionParameter             CodeType = 311
	CodeFunctionNotAvailableForFunctionType  CodeType = 312
	CodeFunctionRequiresNonZeroCurrentSupply CodeType = 313

	// Token/coin names
	CodeReserveTokenInvalid     CodeType = 314
	CodeMaxSupplyDenomInvalid   CodeType = 315
	CodeBondTokenInvalid        CodeType = 316
	CodeReserveDenomsMismatch   CodeType = 317
	CodeInvalidCoinDenomination CodeType = 318
	// Bonds
	CodeBondDoesNotExist        CodeType = 304
	CodeBondAlreadyExists       CodeType = 305
	CodeBondDoesNotAllowSelling CodeType = 306
	CodeDidNotEditAnything      CodeType = 307
	CodeInvalidSwapper          CodeType = 308
	CodeInvalidBond             CodeType = 309
	// Amounts and fees
	CodeInvalidResultantSupply CodeType = 319
	CodeMaxPriceExceeded       CodeType = 320
	CodeSwapAmountInvalid      CodeType = 321
	CodeOrderLimitExceeded     CodeType = 322
	CodeSanityRateViolated     CodeType = 323
	CodeFeeTooLarge            CodeType = 324
	CodeNameDoesNotExist       CodeType = 325
	CodeInternalBondDic        CodeType = 326
	CodeInvalidBasicMsg        CodeType = 150
	CodeBadDataValue           CodeType = 151
	CodeUnauthorizedPermission CodeType = 152
	CodeItemDuplication        CodeType = 153
	CodeItemNotFound           CodeType = 154
	CodeInvalidState           CodeType = 155
	CodeBadWasmExecution       CodeType = 156
	CodeOnlyOneDenomAllowed    CodeType = 157
	CodeInvalidDenom           CodeType = 158
	CodeUnknownClientID        CodeType = 159

	//payments
	CodeInvalidDistribution          CodeType = 101
	CodeInvalidShare                 CodeType = 102
	CodeInvalidPeriod                CodeType = 103
	CodeInvalidPaymentContractAction CodeType = 104
	CodeInvalidDiscount              CodeType = 105
	CodeInvalidDiscountRequest       CodeType = 106
	CodeInvalidPaymentTemplate       CodeType = 107
	CodeInvalidSubscriptionAction    CodeType = 108
	CodeInvalidId                    CodeType = 109
	CodeInvalidArgument              CodeType = 110
	CodeAlreadyExists                CodeType = 111
	CodeInvalidCoin                  CodeType = 112
)

var (
	// ErrUnauthorized is used whenever a request without sufficient
	// authorization is handled.
	ErrorInvalidDidE                         = errors.Register(moduleNameDid, CodeInvalidDid, "invalid did")
	ErrorInvalidPubKey                       = errors.Register(moduleNameDid, CodeInvalidPubKey, "invalid pubkey")
	ErrorInvalidIssuer                       = errors.Register(moduleNameDid, CodeInvalidIssuer, "invalid issuer")
	ErrorInvalidCredentials                  = errors.Register(moduleNameDid, CodeInvalidCredentials, "Data already exist")
	ErrNameDoesNotExist                      = errors.Register(moduleNameBonddoc, CodeNameDoesNotExist, "name does not exist")
	ErrInternalE                             = errors.Register(moduleNameBonddoc, CodeInternalBondDic, "bond did not found")
	ErrGasOverflow                           = errors.Register(moduleNameBonddoc, CodeInvalidDid, "Gas invalid supply")
	ErrArgument                              = errors.Register(moduleNameBonds, CodeArgumentInvalid, "Cannot be empty")
	ErrArgumentMissingOrIncorrectType        = errors.Register(moduleNameBonds, CodeArgumentMissingOrIncorrectType, "Missing or Incorrect Type")
	ErrCodeIncorrectNumberOfValues           = errors.Register(moduleNameBonds, CodeIncorrectNumberOfValues, "Incorrect code number of value")
	ErrUnrecognizedFunctionType              = errors.Register(moduleNameBonds, CodeUnrecognizedFunctionType, "Unrecognized function type")
	ErrCodeInvalidFuncParam                  = errors.Register(moduleNameBonds, CodeInvalidFunctionParameter, "Invalid Function Parameter")
	ErrFunctionNotAvailableForFunctionType   = errors.Register(moduleNameBonds, CodeFunctionNotAvailableForFunctionType, "Function is not available for the function type")
	ErrFunctionRequiresNonZeroCurrentSupply  = errors.Register(moduleNameBonds, CodeFunctionRequiresNonZeroCurrentSupply, "Function requires the current supply to be non zero")
	ErrTokenIsNotAValidReserveTokenCode      = errors.Register(moduleNameBonds, CodeReserveTokenInvalid, "Function requires the current supply to be non zero")
	ErrMaxSupplyDenomDoesNotMatchTokenDenomE = errors.Register(moduleNameBonds, CodeMaxSupplyDenomInvalid, "Max supply denom does not match token denom")
	ErrBondInvalidToken                      = errors.Register(moduleNameBonds, CodeBondTokenInvalid, "bond token is invalid")
	ErrReserveDenomsMismatchE                = errors.Register(moduleNameBonds, CodeReserveDenomsMismatch, "reserve denom mismatch")
	ErroInvalidCoinDenomination              = errors.Register(moduleNameBonds, CodeInvalidCoinDenomination, "wrong coin denomination")
	EInvalidResultantSupply                  = errors.Register(moduleNameBonds, CodeInvalidResultantSupply, "Invalid resultant supply")
	EPriceExceed                             = errors.Register(moduleNameBonds, CodeMaxPriceExceeded, "price exceeded")
	ESwapAmountInvalid                       = errors.Register(moduleNameBonds, CodeSwapAmountInvalid, "invalid amount in swap")
	ErrOrderQuantityLimitExceeded            = errors.Register(moduleNameBonds, CodeOrderLimitExceeded, "Order quantity limits exceeded")
	ErrValuesViolateSanityRate               = errors.Register(moduleNameBonds, CodeSanityRateViolated, "Values violate sanity rate")
	ErrFeesCannotBeOrExceed100Percent        = errors.Register(moduleNameBonds, CodeFeeTooLarge, "Sum of fees is or exceeds 100 percent")
	ErrInvalidBasicMsg                       = errors.Register(moduleNameIxo, CodeInvalidBasicMsg, "Invalid Basic Message")
	ErrBadDataValue                          = errors.Register(moduleNameIxo, CodeBadDataValue, "Bad Data Value")
	ErrUnauthorizedPermission                = errors.Register(moduleNameIxo, CodeUnauthorizedPermission, "Unauthorized Permission")
	ErrItemDuplication                       = errors.Register(moduleNameIxo, CodeItemDuplication, "Item Duplication")
	ErrItemNotFound                          = errors.Register(moduleNameIxo, CodeItemNotFound, "Item Not Found")
	ErrInvalidState                          = errors.Register(moduleNameIxo, CodeInvalidState, "InvalidState")
	ErrBadWasmExecution                      = errors.Register(moduleNameIxo, CodeBadWasmExecution, "Bad Wasm Execution")
	ErrOnlyOneDenomAllowed                   = errors.Register(moduleNameIxo, CodeOnlyOneDenomAllowed, "Only One Denom Allowed")
	ErrInvalidDenom                          = errors.Register(moduleNameIxo, CodeInvalidDenom, "Invalid Denom")
	ErrUnknownClientID                       = errors.Register(moduleNameIxo, CodeUnknownClientID, "Unknown Client ID")
	ErrInvalidDistribution                   = errors.Register(moduleNamePayment, CodeInvalidDistribution, "payment invalid")
	EInvalidShare                            = errors.Register(moduleNamePayment, CodeInvalidShare, "payment invalid")
	EInvalidPeriod                           = errors.Register(moduleNamePayment, CodeInvalidPeriod, "payment invalid")
	EInvalidPaymentCA                        = errors.Register(moduleNamePayment, CodeInvalidPaymentContractAction, "payment invalid")
	EInvalidDiscount                         = errors.Register(moduleNamePayment, CodeInvalidDiscount, "payment invalid")
	EInvalidDiscountReq                      = errors.Register(moduleNamePayment, CodeInvalidDiscountRequest, "payment invalid")
	EInvalidPaymentTemplate                  = errors.Register(moduleNamePayment, CodeInvalidPaymentTemplate, "payment invalid")
	EInvalidSubAction                        = errors.Register(moduleNamePayment, CodeInvalidSubscriptionAction, "payment invalid")
	EInvalidId                               = errors.Register(moduleNamePayment, CodeInvalidId, "payment invalid")
	EInvalidArgs                             = errors.Register(moduleNamePayment, CodeInvalidArgument, "payment invalid")
	EAlreadyExists                           = errors.Register(moduleNamePayment, CodeAlreadyExists, "payment invalid")
	EInvalidCoin                             = errors.Register(moduleNameProject, CodeInvalidCoin, "coin is invalid")
	ErrCodeBondDoesNotExist                  = errors.Register(moduleNameBonds, CodeBondDoesNotExist, "Code bond does not exist")
	ErrCodeBondAlreadyExists                 = errors.Register(moduleNameBonds, CodeBondAlreadyExists, "Code bond already exist")
	ErrCodeBondDoesNotAllowSelling           = errors.Register(moduleNameBonds, CodeBondDoesNotAllowSelling, "Code bond does not allow selling")
	ErrCodeDidNotEditAnything                = errors.Register(moduleNameBonds, CodeDidNotEditAnything, "Did not edit anything from the bond.")
	ErrFromAndToCannotBeTheSameToken_E       = errors.Register(moduleNameBonds, CodeInvalidSwapper, "From and To tokens cannot be the same token.")
	ErrDuplicateReserveToken                 = errors.Register(moduleNameBonds, CodeInvalidBond, "Cannot have duplicate tokens in reserve tokens.")
)

func ErrInvalidAddress(arg string) error {
	return errors.Wrapf(errors.ErrInvalidAddress, arg)
}
func ErrArgumentCannotBeEmpty(argument string) error {
	return errors.Wrapf(ErrArgument, "%s argument cannot be empty", argument)
}
func ErrArgumentCannotBeNegative(arg string) error {
	return errors.Wrapf(ErrArgument, "%s argument cannot be negative", arg)
}
func ErrArgumentMustBePositive(arg string) error {
	return errors.Wrapf(ErrArgument, "%s argument must be a positive value", arg)
}
func ErrMaxSupplyDenomDoesNotMatchTokenDenom() error {
	return errors.Wrap(ErrMaxSupplyDenomDoesNotMatchTokenDenomE, "Max supply denom does not match token denom")
}
func ErrFunctionParameterMissingOrNonInteger(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s parameter is missing or is not an integer", arg)
}
func ErrArgumentMissingOrNonFloat(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s argument is missing or is not a float", arg)
}
func ErrArgumentMissingOrNonInteger(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s argument is missing or is not an integer", arg)
}
func ErrArgumentMissingOrNonUInteger(arg string) error {
	return errors.Wrapf(ErrArgumentMissingOrIncorrectType, "%s argument is missing or is not an unsigned integer", arg)
}
func ErrArgumentMissingOrNonBoolean(arg string) error {
	return errors.Wrapf(ErrCodeIncorrectNumberOfValues, "%s argument is missing or is not true or false", arg)
}
func ErrIncorrectNumberOfReserveTokens(expected int) error {
	return errors.Wrapf(ErrCodeIncorrectNumberOfValues, "Incorrect number of reserve tokens; expected: %d", expected)
}
func ErrorBondDocAlreadyExist() error {
	return errors.Wrap(ErrorInvalidDidE, "Bond doc already exists")
}
func ErrInvalidDid(args string) error {
	return errors.Wrap(ErrorInvalidDidE, args)
}
func ErrInvalidCoins(args string) error {
	return errors.Wrap(EInvalidCoin, args)
}
func ErrIncorrectNumberOfFunctionParameters(expected int) error {
	return errors.Wrapf(ErrCodeIncorrectNumberOfValues, "Incorrect number of function parameters; expected: %d", expected)
}
func ErrInvalidFunctionParameter(parameter string) error {
	return errors.Wrapf(ErrCodeInvalidFuncParam, "Invalid function parameter '%s'", parameter)
}
func ErrTokenIsNotAValidReserveToken(denom string) error {
	return errors.Wrapf(ErrTokenIsNotAValidReserveTokenCode, "Token '%s' is not a valid reserve token", denom)
}
func ErrBondTokenCannotAlsoBeReserveToken() error {
	errMsg := "Token cannot also be a reserve token"
	return errors.Wrap(ErrBondInvalidToken, errMsg)
}
func ErrBondTokenCannotBeStakingToken() error {
	errMsg := "Bond token cannot be staking token"
	return errors.Wrap(ErrBondInvalidToken, errMsg)
}
func ErrBondTokenDoesNotMatchBond() error {
	errMsg := "Bond token does not match bond"
	return errors.Wrap(ErrBondInvalidToken, errMsg)
}
func ErrReserveDenomsMismatch(inputDenoms string, actualDenoms []string) error {
	return errors.Wrapf(ErrReserveDenomsMismatchE, "Denoms in %s do not match reserve denoms; expected: %s", inputDenoms, strings.Join(actualDenoms, ","))
}
func ErrInvalidCoinDenomination(denom string) error {
	return errors.Wrapf(ErroInvalidCoinDenomination, "Invalid coin denomination '%s'", denom)
}
func ErrCannotMintMoreThanMaxSupply() error {
	errMsg := "Cannot mint more tokens than the max supply"
	return errors.Wrap(EInvalidResultantSupply, errMsg)
}
func ErrCannotBurnMoreThanSupply() error {
	errMsg := "Cannot burn more tokens than the current supply"
	return errors.Wrap(EInvalidResultantSupply, errMsg)
}
func ErrMaxPriceExceeded(totalPrice, maxPrice sdk.Coins) error {
	return errors.Wrapf(EPriceExceed, "Actual prices %s exceed max prices %s", totalPrice.String(), maxPrice.String())
}
func ErrSwapAmountTooSmallToGiveAnyReturn(fromToken, toToken string) error {
	return errors.Wrapf(ESwapAmountInvalid, "%s swap amount too small to give any %s return", fromToken, toToken)
}
func ErrSwapAmountCausesReserveDepletion(fromToken, toToken string) error {
	return errors.Wrapf(ESwapAmountInvalid, "%s swap amount too large and causes %s reserve to be depleted", fromToken, toToken)
}
func UnknownRequest(m string) error {
	return errors.Wrap(errors.ErrUnknownRequest, m)
}
func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}
func IntErr(m string) error {
	return errors.Wrap(errors.ErrPanic, m)
}
func ErrJsonMars(m string) error {
	return errors.Wrapf(errors.ErrJSONMarshal, "Json marshall error %s", m)
}
