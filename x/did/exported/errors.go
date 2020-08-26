package exported

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// Local code type
type CodeType = uint32
type CodespaceType = string

const (
	moduleNameDid     = "did"
	moduleNameBonddoc = "bonddoc"
	moduleNameIxo     = "dap"
	moduleNamePayment = "payments"
	moduleNameProject = "project"

	CodeOK                           CodeType = 0
	CodeInvalidDid                   CodeType = 201
	CodeInvalidPubKey                CodeType = 202
	CodeInvalidIssuer                CodeType = 203
	CodeInvalidCredentials           CodeType = 204
	CodeNameDoesNotExist             CodeType = 325
	CodeInternalBondDic              CodeType = 326
	CodeInvalidBasicMsg              CodeType = 150
	CodeBadDataValue                 CodeType = 151
	CodeUnauthorizedPermission       CodeType = 152
	CodeItemDuplication              CodeType = 153
	CodeItemNotFound                 CodeType = 154
	CodeInvalidState                 CodeType = 155
	CodeBadWasmExecution             CodeType = 156
	CodeOnlyOneDenomAllowed          CodeType = 157
	CodeInvalidDenom                 CodeType = 158
	CodeUnknownClientID              CodeType = 159
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
	CodeInsufficientCoins            CodeType = 113
)

var (
	// ErrUnauthorized is used whenever a request without sufficient
	// authorization is handled.
	ErrorInvalidDidE          = errors.Register(moduleNameDid, CodeInvalidDid, "invalid did")
	ErrorInvalidPubKey        = errors.Register(moduleNameDid, CodeInvalidPubKey, "invalid pubkey")
	ErrorInvalidIssuer        = errors.Register(moduleNameDid, CodeInvalidIssuer, "invalid issuer")
	ErrorInvalidCredentials   = errors.Register(moduleNameDid, CodeInvalidCredentials, "Data already exist")
	ErrNameDoesNotExist       = errors.Register(moduleNameBonddoc, CodeNameDoesNotExist, "name does not exist")
	ErrInternalE              = errors.Register(moduleNameBonddoc, CodeInternalBondDic, "bond did not found")
	ErrGasOverflow            = errors.Register(moduleNameBonddoc, CodeInvalidDid, "Gas invalid supply")
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
	ErrInsufficientCoins      = errors.Register(moduleNameIxo, CodeInsufficientCoins, "Unknown Client ID")
	ErrInvalidDistribution    = errors.Register(moduleNamePayment, CodeInvalidDistribution, "payment invalid")
	EInvalidShare             = errors.Register(moduleNamePayment, CodeInvalidShare, "payment invalid")
	EInvalidPeriod            = errors.Register(moduleNamePayment, CodeInvalidPeriod, "payment invalid")
	EInvalidPaymentCA         = errors.Register(moduleNamePayment, CodeInvalidPaymentContractAction, "payment invalid")
	EInvalidDiscount          = errors.Register(moduleNamePayment, CodeInvalidDiscount, "payment invalid")
	EInvalidDiscountReq       = errors.Register(moduleNamePayment, CodeInvalidDiscountRequest, "payment invalid")
	EInvalidPaymentTemplate   = errors.Register(moduleNamePayment, CodeInvalidPaymentTemplate, "payment invalid")
	EInvalidSubAction         = errors.Register(moduleNamePayment, CodeInvalidSubscriptionAction, "payment invalid")
	EInvalidId                = errors.Register(moduleNamePayment, CodeInvalidId, "payment invalid")
	EInvalidArgs              = errors.Register(moduleNamePayment, CodeInvalidArgument, "payment invalid")
	EAlreadyExists            = errors.Register(moduleNamePayment, CodeAlreadyExists, "payment invalid")
	EInvalidCoin              = errors.Register(moduleNameProject, CodeInvalidCoin, "coin is invalid")
)

func ErrInvalidDid(args string) error {
	return errors.Wrap(ErrorInvalidDidE, args)
}
func ErrInvalidCoins(args string) error {
	return errors.Wrap(EInvalidCoin, args)
}
func ErrInvalidAddress(arg string) error {
	return errors.Wrapf(errors.ErrInvalidAddress, arg)
}
func InvalidAddress(m string) error {
	return errors.Wrap(errors.ErrInvalidAddress, m)
}
func UnknownRequest(m string) error {
	return errors.Wrap(errors.ErrUnknownRequest, m)
}
func UnknownRequestErr(m string, err error) error {
	return errors.Wrap(errors.ErrUnknownRequest, fmt.Sprintf("%s :%s", m, err.Error()))
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
func ErrUnmarshalJson(m string) error {
	return errors.Wrapf(errors.ErrJSONUnmarshal, "technical error in %s", m)
}
func InsufficientCoins(m string) error {
	return errors.Wrapf(ErrInsufficientCoins, "%s: %s", moduleNameIxo, m)
}
