package errors

type CodeType = uint32

const (
	ModuleName = "bonds"
	// Bonds
	CodeBondDoesNotExist        CodeType = 304
	CodeBondAlreadyExists       CodeType = 305
	CodeBondDoesNotAllowSelling CodeType = 306
	CodeDidNotEditAnything      CodeType = 307
	CodeInvalidSwapper          CodeType = 308
	CodeInvalidBond             CodeType = 309
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
	CodeParseCodeInvalid        CodeType = 325

	// Amounts and fees
	CodeInvalidResultantSupply CodeType = 319
	CodeMaxPriceExceeded       CodeType = 320
	CodeSwapAmountInvalid      CodeType = 321
	CodeOrderLimitExceeded     CodeType = 322
	CodeSanityRateViolated     CodeType = 323
	CodeFeeTooLarge            CodeType = 324
)
