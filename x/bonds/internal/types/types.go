package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/dp-block/x/bonds/errors"
	"github.com/tokenchain/dp-block/x/did/exported"
)

const (
	TRUE  = "true"
	FALSE = "false"
)

func CheckReserveTokenNames(resTokens []string, token string) error {
	// Check that no token is the same as the main token, no token
	// is duplicate, and that the token is a valid denomination
	uniqueReserveTokens := make(map[string]string)
	for _, r := range resTokens {
		// Check if same as main token
		if r == token {
			return errors.BondTokenCannotAlsoBeReserveToken()
		}

		// Check if duplicate
		if _, ok := uniqueReserveTokens[r]; ok {
			return errors.DuplicateReserveToken()
		} else {
			uniqueReserveTokens[r] = ""
		}

		// Check if can be parsed as coin
		err := CheckCoinDenom(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckNoOfReserveTokens(resTokens []string, fnType string) error {
	// Come up with number of expected reserve tokens
	expectedNoOfTokens, ok := NoOfReserveTokensForFunctionType[fnType]
	if !ok {
		return errors.UnrecognizedFunctionType()
	}

	// Check that number of reserve tokens is correct (if expecting a specific number of tokens)
	if expectedNoOfTokens != AnyNumberOfReserveTokens && len(resTokens) != expectedNoOfTokens {
		return errors.ErrIncorrectNumberOfReserveTokens(expectedNoOfTokens)
	}

	return nil
}

func CheckCoinDenom(denom string) (err error) {
	coin, err2 := sdk.ParseCoin("0" + denom)
	if err2 != nil {
		return exported.IntErr(err2.Error())
	} else if denom != coin.Denom {
		return errors.InvalidCoinDenomination(denom)
	}
	return nil
}

func GetRequiredParamsForFunctionType(fnType string) (fnParams []string, err error) {
	expectedParams, ok := RequiredParamsForFunctionType[fnType]
	if !ok {
		return nil, errors.UnrecognizedFunctionType()
	}
	return expectedParams, nil
}

func GetExceptionsForFunctionType(fnType string) (restrictions FunctionParamRestrictions, err error) {
	restrictions, ok := ExtraParameterRestrictions[fnType]
	if !ok {
		return nil, errors.UnrecognizedFunctionType()
	}
	return restrictions, nil
}
