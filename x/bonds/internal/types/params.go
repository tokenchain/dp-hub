package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tokenchain/dp-hub/x/did/exported"
)

type (
	Params struct {
		ListingDid exported.Did `json:"listing_did" yaml:"listing_did"`
	}
)

var (
	KeyListingDid                       = []byte("ListingDid")
)

// ParamTable for project module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(ixoDid exported.Did) Params {
	return Params{
		ListingDid: ixoDid,
	}

}

// default project module parameters
func DefaultParams() Params {
	return Params{
		ListingDid: exported.Did("N/A"), // blank
	}
}

// validate params
func ValidateParams(params Params) error {
	if len(params.ListingDid) == 0 {
		return fmt.Errorf("DXP DID cannot be empty ... %s", params.ListingDid)
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Project Params:
  Listing Did: %s
`, p.ListingDid)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyListingDid, &p.ListingDid, listingValidation},
	}
}

func listingValidation(i interface{}) error {
	_, ok := i.(exported.Did)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}