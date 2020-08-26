package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
)

const (
	CodeInvalidPriceDigit       did.CodeType = 1
	CodeInvalidMinTradeSize     did.CodeType = 2
	CodeInvalidDexList          did.CodeType = 3
	CodeInvalidBalanceNotEnough did.CodeType = 4
	CodeInvalidHeight           did.CodeType = 5
	CodeInvalidAsset            did.CodeType = 6
	CodeInvalidCommon           did.CodeType = 7
	CodeBlockedRecipient        did.CodeType = 8
	CodeSendDisabled            did.CodeType = 9
)

var (
	ErrInvalidPriceDigit       = errors.Register(ModuleName, CodeInvalidPriceDigit, "invalid price")
	ErrSendDisabled            = errors.Register(ModuleName, CodeSendDisabled, "invalid price")
	ErrBlockedRecipient        = errors.Register(ModuleName, CodeBlockedRecipient, "invalid price")
	ErrInvalidCommon           = errors.Register(ModuleName, CodeInvalidCommon, "invalid price")
	ErrInvalidAsset            = errors.Register(ModuleName, CodeInvalidAsset, "invalid price")
	ErrInvalidHeight           = errors.Register(ModuleName, CodeInvalidHeight, "invalid price")
	ErrInvalidBalanceNotEnough = errors.Register(ModuleName, CodeInvalidBalanceNotEnough, "invalid price")
	ErrInvalidDexList          = errors.Register(ModuleName, CodeInvalidDexList, "invalid price")
	ErrInvalidMinTradeSize     = errors.Register(ModuleName, CodeInvalidMinTradeSize, "invalid price")
)

// ErrBlockedRecipient returns an error when a transfer is tried on a blocked recipient
func BlockedRecipient(blockedAddr string) error {
	return errors.Wrapf(ErrInvalidPriceDigit, "failed. %s is not allowed to receive transactions", blockedAddr)
}

// ErrSendDisabled returns an error when the transaction sending is disabled in bank module
func SendDisabled() error {
	return errors.Wrap(ErrSendDisabled, "failed. send transactions are currently disabled")
}

func InvalidDexList(message string) error {
	return errors.Wrap(ErrInvalidDexList, message)
}

func InvalidBalanceNotEnough(message string) error {
	return errors.Wrap(ErrInvalidBalanceNotEnough, message)
}

func InvalidHeight(h, ch, max int64) error {
	return errors.Wrapf(ErrInvalidHeight, "Height %d must be greater than current block height %d and less than %d + %d.", h, ch, ch, max)
}

func InvalidCommon(message string) error {
	return errors.Wrap(ErrInvalidCommon, message)
}

func InvalidAsset(message string) error {
	return errors.Wrap(ErrInvalidAsset, message)
}
func InvalidMinTradeSize(message string) error {
	return errors.Wrap(ErrInvalidMinTradeSize, message)
}
