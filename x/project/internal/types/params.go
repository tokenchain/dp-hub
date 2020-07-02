package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

var (
	KeyIxoDid                       = []byte("IxoDid")
	KeyProjectMinimumInitialFunding = []byte("ProjectMinimumInitialFunding")
)

type Params struct {
	IxoDid                       types.Did `json:"dp_did" yaml:"dp_did"`
	ProjectMinimumInitialFunding sdk.Dec   `json:"project_minimum_initial_funding" yaml:"project_minimum_initial_funding"`
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyIxoDid, &p.IxoDid, ixoValidation},
		{KeyProjectMinimumInitialFunding, &p.ProjectMinimumInitialFunding, projectminiValidation},
	}
}
func (p Params) String() string {
	return fmt.Sprintf(`Project Params:
  Ixo Did: %s
  Project Minimum Initial Funding: %s`, p.ProjectMinimumInitialFunding, p.IxoDid)
}

// ParamTable for project module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(pif sdk.Dec, pixDid types.Did) Params {
	return Params{
		IxoDid:                       pixDid,
		ProjectMinimumInitialFunding: pif,
	}
}

// default project module parameters
func DefaultParams() Params {
	return Params{
		IxoDid:                       types.Did(""),                          // Blank
		ProjectMinimumInitialFunding: sdk.NewDec(500).Mul(types.IxoDecimals), // 500.000
	}
}

// validate params
func ValidateParams(params Params) error {
	if len(params.IxoDid) == 0 {
		return fmt.Errorf("DAP did cannot be empty...")
	}
	if params.ProjectMinimumInitialFunding.LT(sdk.ZeroDec()) {
		return fmt.Errorf("Project parameter ProjectMinimumInitialFunding should be positive, is %s ", params.ProjectMinimumInitialFunding.String())
	}
	return nil
}
func ixoValidation(i interface{}) error {
	_, ok := i.(types.Did)
	if !ok {
		return fmt.Errorf("ixoValidation Invalid parameter type: %T.", i)
	}
	return nil
}
func projectminiValidation(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("pv invalid params type: %T .", i)
	}
/*	if params.LT(sdk.ZeroDec()) {
		return fmt.Errorf("Project params project minimum initial funding should be positive, is %s .", params.String())
	}*/
	return nil
}
