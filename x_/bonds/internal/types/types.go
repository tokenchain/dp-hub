package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ro "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x"
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
			return x.ErrBondTokenCannotAlsoBeReserveToken()
		}

		// Check if duplicate
		if _, ok := uniqueReserveTokens[r]; ok {
			return ro.Wrap(x.ErrDuplicateReserveToken, "")
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
		return ro.Wrap(x.ErrUnrecognizedFunctionType, "")
	}

	// Check that number of reserve tokens is correct (if expecting a specific number of tokens)
	if expectedNoOfTokens != AnyNumberOfReserveTokens && len(resTokens) != expectedNoOfTokens {
		return x.ErrIncorrectNumberOfReserveTokens(expectedNoOfTokens)
	}

	return nil
}

func CheckCoinDenom(denom string) (err error) {
	coin, err2 := sdk.ParseCoin("0" + denom)
	if err2 != nil {
		return x.IntErr(err2.Error())
	} else if denom != coin.Denom {
		return x.ErrInvalidCoinDenomination(denom)
	}
	return nil
}

func GetRequiredParamsForFunctionType(fnType string) (fnParams []string, err error) {
	expectedParams, ok := RequiredParamsForFunctionType[fnType]
	if !ok {
		return nil, ro.Wrap(x.ErrUnrecognizedFunctionType, "")
	}
	return expectedParams, nil
}
