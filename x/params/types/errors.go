package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
)

// const
const (
	CodeInvalidMaxProposalNum did.CodeType = 4
	CodeInvalidHeight         did.CodeType = 5
)

var (
	Err10 = errors.Register(types.ModuleName, CodeInvalidMaxProposalNum, "invalid product")
	Err11 = errors.Register(types.ModuleName, CodeInvalidHeight, "invalid height")
)

// ErrInvalidMaxProposalNum returns error when the number of params to change are out of limit
func ErrInvalidMaxProposalNum(msg string) error {
	return errors.Wrap(Err10, msg)
}

//return sdk.NewError(codespace, CodeInvalidMaxProposalNum, msg)
func ErrInvalidProposalContent(msg string) error {
	return errors.Wrapf(govtypes.ErrInvalidProposalContent, "%s: %s", types.ModuleName, msg)
}
func ErrInvalidProposalType(msg string) error {
	return errors.Wrapf(govtypes.ErrInvalidProposalType, "%s: %s", types.ModuleName, msg)
}
func UnknownSubspace(msg string) error {
	return errors.Wrapf(sdkparams.ErrUnknownSubspace, "%s: %s", types.ModuleName, msg)
}
func ErrSettingParam(key, value, msg string) error {
	return errors.Wrapf(sdkparams.ErrSettingParameter, "%s: %s / %s / %s", types.ModuleName, key, value, msg)
}
func InvalidProposer(msg string) error {
	return errors.Wrapf(govtypes.ErrInvalidProposalType, "%s: %s", types.ModuleName, msg)
}
func InvalidHeight(height, c, d uint64) error {
	return errors.Wrapf(Err11, "%s: %d,%d,%d", types.ModuleName, height, c, d)
}
