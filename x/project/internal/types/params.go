package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

// Parameter store keys
var (
	KeyIxoDid                       = []byte("IxoDid")
	KeyProjectMinimumInitialFunding = []byte("ProjectMinimumInitialFunding")
)

// project parameters
type Params struct {
	IxoDid                       types.Did `json:"dp_did" yaml:"dp_did"`
	ProjectMinimumInitialFunding sdk.Int   `json:"project_minimum_initial_funding" yaml:"project_minimum_initial_funding"`
}

// ParamTable for project module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(projectMinimumInitialFunding sdk.Int, ixoDid types.Did) Params {
	return Params{
		IxoDid:                       ixoDid,
		ProjectMinimumInitialFunding: projectMinimumInitialFunding,
	}

}

// default project module parameters
func DefaultParams() Params {
	return Params{
		IxoDid:                       types.Did(""), // blank
		ProjectMinimumInitialFunding: sdk.OneInt(),  // 1
	}
}

// validate params
func ValidateParams(params Params) error {
	if len(params.IxoDid) == 0 {
		return fmt.Errorf("ixo did cannot be empty...")
	}
	if params.ProjectMinimumInitialFunding.LT(sdk.ZeroInt()) {
		return fmt.Errorf("project parameter ProjectMinimumInitialFunding should be positive, is %s ", params.ProjectMinimumInitialFunding.String())
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Project Params:
  Ixo Did: %s
  Project Minimum Initial Funding: %s

`, p.ProjectMinimumInitialFunding, p.IxoDid)
}
func ixoValidation(i interface{}) error {
	_, ok := i.(types.Did)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func projectminiValidation(i interface{}) error {
	_, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyIxoDid, &p.IxoDid, ixoValidation},
		{KeyProjectMinimumInitialFunding, &p.ProjectMinimumInitialFunding, projectminiValidation},
	}
}
