package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type CodeType = uint32

const (
	codeInvalidProduct          CodeType = 1
	codeTokenPairNotFound       CodeType = 2
	codeDelistOwnerNotMatch     CodeType = 3
	codeInvalidBalanceNotEnough CodeType = 4
	codeInvalidAsset            CodeType = 5
	codeUnknownOperator         CodeType = 6
	codeExistOperator           CodeType = 7
	codeInvalidWebsiteLength    CodeType = 8
	codeInvalidWebsiteURL       CodeType = 9
	codeInvalidProposer         CodeType = 10
)

var (
	Err1  = errors.Register(ModuleName, codeInvalidProduct, "invalid product")
	Err2  = errors.Register(ModuleName, codeTokenPairNotFound, "token pair not found")
	Err3  = errors.Register(ModuleName, codeDelistOwnerNotMatch, "delist owner not match")
	Err4  = errors.Register(ModuleName, codeInvalidBalanceNotEnough, "invalid balance not enough")
	Err5  = errors.Register(ModuleName, codeInvalidAsset, "invalid asset")
	Err6  = errors.Register(ModuleName, codeUnknownOperator, "unknown operator")
	Err7  = errors.Register(ModuleName, codeExistOperator, "exist operator")
	Err8  = errors.Register(ModuleName, codeInvalidWebsiteLength, "invalid website length")
	Err9  = errors.Register(ModuleName, codeInvalidWebsiteURL, "invalid website url")
	Err10 = errors.Register(ModuleName, codeInvalidProposer, "invalid website url")
)

// ErrInvalidProduct returns invalid product error
func ErrInvalidProduct(msg string) error {
	return errors.Wrapf(Err1, "%s :%s", Err1.Error(), msg)
}

// ErrTokenPairNotFound returns token pair not found error
func ErrTokenPairNotFound(msg string) error {
	return errors.Wrapf(Err2, "%s :%s", Err2.Error(), msg)
}

// ErrDelistOwnerNotMatch returns delist owner not match error
func ErrDelistOwnerNotMatch(msg string) error {
	return errors.Wrapf(Err3, "%s :%s", Err3.Error(), msg)
}

// ErrInvalidBalanceNotEnough returns invalid balance not enough error
func ErrInvalidBalanceNotEnough(message string) error {
	return errors.Wrapf(Err4, "%s: %s", Err4.Error(), message)
}

// ErrInvalidAsset returns invalid asset error
func ErrInvalidAsset(message string) error {
	return errors.Wrapf(Err5, "%s: %s", Err5.Error(), message)
}

func ErrUnknownOperator(addr sdk.AccAddress) error {
	return errors.Wrapf(Err6, "%s: %s", Err6.Error(), addr.String())
}

func ErrExistOperator(addr sdk.AccAddress) error {
	return errors.Wrapf(Err7, "%s: %s", Err7.Error(), addr.String())
}

func ErrInvalidWebsiteLength(got, max int) error {
	return errors.Wrapf(Err8, "invalid website length, got length %v, max is %v", got, max)
}

func ErrInvalidWebsiteURL(msg string) error {
	return errors.Wrapf(Err9, "invalid website URL: %s", msg)
}

// ErrTokenPairExisted returns an error when the token pair is existed during the process of listing
// ErrTokenPairExisted returns an error when the token pair is existing during the process of listing
func ErrTokenPairExisted(baseAsset, quoteAsset string) error {
	return errors.Wrapf(Err5, "failed. the token pair exists with %s and %s", baseAsset, quoteAsset)
}
func ErrInvalidProposer(msg string) error {
	return errors.Wrap(Err10, msg)
}

func ErrInvalidProposalContent(msg string) error {
	return errors.Wrapf(govtypes.ErrInvalidProposalContent, "%s. %s", DefaultCodespace, msg)
}

func ErrInvalidProposalType(msg string) error {
	return errors.Wrapf(govtypes.ErrInvalidProposalType, "%s. %s", DefaultCodespace, msg)
}

func ErrInvalidAddress(msg string) error {
	return errors.Wrapf(errors.ErrInvalidAddress, "%s. %s", DefaultCodespace, msg)
}

func ErrInvalidCoins(msg string) error {
	return errors.Wrapf(errors.ErrInvalidCoins, "%s. %s", DefaultCodespace, msg)
}
