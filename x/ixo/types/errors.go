package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
var (
	ErrInvalidBasicMsg        = sdkerrors.Register(ModuleName, 1, "InvalidBasicMsg")
	ErrBadDataValue           = sdkerrors.Register(ModuleName, 2, "BadDataValue")
	ErrUnauthorizedPermission = sdkerrors.Register(ModuleName, 3, "UnauthorizedPermission")
	ErrItemDuplication        = sdkerrors.Register(ModuleName, 4, "ItemDuplication")
	ErrItemNotFound           = sdkerrors.Register(ModuleName, 5, "ItemNotFound")
	ErrInvalidState           = sdkerrors.Register(ModuleName, 6, "InvalidState")
	ErrBadWasmExecution       = sdkerrors.Register(ModuleName, 7, "BadWasmExecution")
	ErrOnlyOneDenomAllowed    = sdkerrors.Register(ModuleName, 8, "OnlyOneDenomAllowed")
	ErrInvalidDenom           = sdkerrors.Register(ModuleName, 9, "InvalidDenom")
	ErrUnknownClientID        = sdkerrors.Register(ModuleName, 10, "UnknownClientID")
)*/

const (
	DefaultCodespace       sdk.CodespaceType = ModuleName
	InvalidBasicMsg        sdk.CodeType      = 150
	BadDataValue           sdk.CodeType      = 151
	UnauthorizedPermission sdk.CodeType      = 152
	ItemDuplication        sdk.CodeType      = 153
	ItemNotFound           sdk.CodeType      = 154
	InvalidState           sdk.CodeType      = 155
	BadWasmExecution       sdk.CodeType      = 156
	OnlyOneDenomAllowed    sdk.CodeType      = 157
	InvalidDenom           sdk.CodeType      = 158
	UnknownClientID        sdk.CodeType      = 159
)

func ErrInvalidBasicMsg(msg string) sdk.Error {
	errMsg := fmt.Sprintf("invalid basic message %s", msg)
	return sdk.NewError(DefaultCodespace, InvalidBasicMsg, errMsg)
}

func ErrBadDataValue(msg string) sdk.Error {
	errMsg := fmt.Sprintf("bad data value %s", msg)
	return sdk.NewError(DefaultCodespace, BadDataValue, errMsg)
}

func ErrUnauthorizedPermission(msg string) sdk.Error {
	errMsg := fmt.Sprintf("permission deny %s", msg)
	return sdk.NewError(DefaultCodespace, UnauthorizedPermission, errMsg)
}

func ErrItemDuplication(msg string) sdk.Error {
	errMsg := fmt.Sprintf("item duplicated %s", msg)
	return sdk.NewError(DefaultCodespace, ItemDuplication, errMsg)
}

func ErrItemNotFound(msg string) sdk.Error {
	errMsg := fmt.Sprintf("item is not found %s", msg)
	return sdk.NewError(DefaultCodespace, ItemNotFound, errMsg)
}

func ErrInvalidState(msg string) sdk.Error {
	errMsg := fmt.Sprintf("invalid state %s", msg)
	return sdk.NewError(DefaultCodespace, InvalidState, errMsg)
}

func ErrBadWasmExecution(msg string) sdk.Error {
	errMsg := fmt.Sprintf("bad wasm execution %s", msg)
	return sdk.NewError(DefaultCodespace, BadWasmExecution, errMsg)
}

func ErrOnlyOneDenomAllowed(msg string) sdk.Error {
	errMsg := fmt.Sprintf("only one denom allowed %s", msg)
	return sdk.NewError(DefaultCodespace, OnlyOneDenomAllowed, errMsg)
}

func ErrInvalidDenom(msg string) sdk.Error {
	errMsg := fmt.Sprintf("invalid denom %s", msg)
	return sdk.NewError(DefaultCodespace, InvalidDenom, errMsg)
}

func ErrUnknownClientID() sdk.Error {
	errMsg := fmt.Sprintf("unknown client ID")
	return sdk.NewError(DefaultCodespace, UnknownClientID, errMsg)
}
