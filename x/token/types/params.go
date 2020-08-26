package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/common"
	"github.com/tokenchain/ixo-blockchain/x/params"
)

const (
	DefaultFeeIssue  = "2500"
	DefaultFeeMint   = "10"
	DefaultFeeBurn   = "10"
	DefaultFeeModify = "0"
	DefaultFeeChown  = "10"
)

var (
	KeyFeeIssue  = []byte("FeeIssue")
	KeyFeeMint   = []byte("FeeMint")
	KeyFeeBurn   = []byte("FeeBurn")
	KeyFeeModify = []byte("FeeModify")
	KeyFeeChown  = []byte("FeeChown")
)

var _ params.ParamSet = &Params{}

// mint parameters
type Params struct {
	FeeIssue  sdk.DecCoin `json:"issue_fee"`
	FeeMint   sdk.DecCoin `json:"mint_fee"`
	FeeBurn   sdk.DecCoin `json:"burn_fee"`
	FeeModify sdk.DecCoin `json:"modify_fee"`
	FeeChown  sdk.DecCoin `json:"transfer_ownership_fee"`
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyFeeIssue, &p.FeeIssue, decCoinValid},
		{KeyFeeMint, &p.FeeMint, decCoinValid},
		{KeyFeeBurn, &p.FeeBurn, decCoinValid},
		{KeyFeeModify, &p.FeeModify, decCoinValid},
		{KeyFeeChown, &p.FeeChown, decCoinValid},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		FeeIssue:  sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeIssue)),
		FeeMint:   sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeMint)),
		FeeBurn:   sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeBurn)),
		FeeModify: sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeModify)),
		FeeChown:  sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeChown)),
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("FeeIssue: %s\n", p.FeeIssue))
	sb.WriteString(fmt.Sprintf("FeeMint: %s\n", p.FeeMint))
	sb.WriteString(fmt.Sprintf("FeeBurn: %s\n", p.FeeBurn))
	sb.WriteString(fmt.Sprintf("FeeModify: %s\n", p.FeeModify))
	sb.WriteString(fmt.Sprintf("FeeChown: %s\n", p.FeeChown))

	return sb.String()
}

func decCoinValid(i interface{}) error {
	_, ok := i.(sdk.DecCoins)
	if !ok {
		return fmt.Errorf("dec coins: %T", i)
	}
	return nil
}
