package types


import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types"
)

var (
	ErrNameDoesNotExist = sdkerrors.NewError(ModuleName, 1, "name does not exist")
)