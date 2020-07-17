package exported

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "did"
)

type CodeType = uint32
type CodespaceType = string

const (
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

	CodeInvalidDid         CodeType = 201
	CodeInvalidPubKey      CodeType = 202
	CodeInvalidIssuer      CodeType = 203
	CodeInvalidCredentials CodeType = 204
	CodeGasOverflow        CodeType = 205

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

	moduleNameBonddoc = "bonddoc"
	moduleNameIxo     = "dap"
	moduleNamePayment = "payments"
	moduleNameProject = "project"
)

var (
	// ErrUnauthorized is used whenever a request without sufficient
	// authorization is handled.
	ErrInvalidBasicMsg        = errors.Register(moduleNameIxo, CodeInvalidBasicMsg, "Invalid Basic Message")
	ErrBadDataValue           = errors.Register(moduleNameIxo, CodeBadDataValue, "Bad Data Value")
	ErrUnauthorizedPermission = errors.Register(moduleNameIxo, CodeUnauthorizedPermission, "Unauthorized Permission")
	ErrItemDuplication        = errors.Register(moduleNameIxo, CodeItemDuplication, "Item Duplication")
	ErrItemNotFound           = errors.Register(moduleNameIxo, CodeItemNotFound, "Item Not Found")
	ErrInvalidState           = errors.Register(moduleNameIxo, CodeInvalidState, "InvalidState")
	ErrBadWasmExecution       = errors.Register(moduleNameIxo, CodeBadWasmExecution, "Bad Wasm Execution")
	ErrOnlyOneDenomAllowed    = errors.Register(moduleNameIxo, CodeOnlyOneDenomAllowed, "Only One Denom Allowed")
	ErrInvalidDenom           = errors.Register(moduleNameIxo, CodeInvalidDenom, "Invalid Denom")
	ErrUnknownClientID        = errors.Register(moduleNameIxo, CodeUnknownClientID, "Unknown Client ID")

	ErrorInvalid            = errors.Register(ModuleName, CodeInvalidDid, "invalid did")
	ErrorInvalidCredentials = errors.Register(ModuleName, CodeInvalidCredentials, "Data already exist")
	ErrGasOverflow          = errors.Register(ModuleName, CodeGasOverflow, "Gas invalid supply")
	ErrInvalidIssuer        = errors.Register(ModuleName, CodeInvalidIssuer, "Invalid did issuer")
	ErrInvalidPubKey        = errors.Register(ModuleName, CodeInvalidPubKey, "Invalid public key")

	EInvalidCoin = errors.Register(moduleNameProject, CodeInvalidCoin, "coin is invalid")

	EInvalidShare           = errors.Register(moduleNamePayment, CodeInvalidShare, "payment invalid")
	EInvalidPeriod          = errors.Register(moduleNamePayment, CodeInvalidPeriod, "payment invalid")
	EInvalidPaymentCA       = errors.Register(moduleNamePayment, CodeInvalidPaymentContractAction, "payment invalid")
	EInvalidDiscount        = errors.Register(moduleNamePayment, CodeInvalidDiscount, "payment invalid")
	EInvalidDiscountReq     = errors.Register(moduleNamePayment, CodeInvalidDiscountRequest, "payment invalid")
	EInvalidPaymentTemplate = errors.Register(moduleNamePayment, CodeInvalidPaymentTemplate, "payment invalid")
	EInvalidSubAction       = errors.Register(moduleNamePayment, CodeInvalidSubscriptionAction, "payment invalid")
	EInvalidId              = errors.Register(moduleNamePayment, CodeInvalidId, "payment invalid")
	EInvalidArgs            = errors.Register(moduleNamePayment, CodeInvalidArgument, "payment invalid")
	EAlreadyExists          = errors.Register(moduleNamePayment, CodeAlreadyExists, "payment invalid")
	ErrInvalidDistribution  = errors.Register(moduleNamePayment, CodeInvalidDistribution, "payment invalid")
)
