package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Assert ParameterChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = ParameterChangeProposal{}

func init() {
	govtypes.RegisterProposalType(sdkparams.ProposalTypeChange)
	govtypes.RegisterProposalTypeCodec(ParameterChangeProposal{}, "darkpool/params/ParameterChangeProposal")
}

// ParameterChangeProposal is the struct of param change proposal
type ParameterChangeProposal struct {
	sdkparams.ParameterChangeProposal
	Height uint64 `json:"height" yaml:"height"`
}

// NewParameterChangeProposal creates a new instance of ParameterChangeProposal
func NewParameterChangeProposal(title, description string, changes []types.ParamChange, height uint64,
) ParameterChangeProposal {
	return ParameterChangeProposal{
		ParameterChangeProposal: sdkparams.NewParameterChangeProposal(title, description, changes),
		Height:                  height,
	}
}

// ValidateBasic validates the parameter change proposal
func (pcp ParameterChangeProposal) ValidateBasic() error {
	if len(strings.TrimSpace(pcp.Title)) == 0 {
		return ErrInvalidProposalContent( "proposal title cannot be blank")
	}
	if len(pcp.Title) > govtypes.MaxTitleLength {
		return ErrInvalidProposalContent(fmt.Sprintf("proposal title is longer than max length of %d", govtypes.MaxTitleLength))
	}

	if len(pcp.Description) == 0 {
		return ErrInvalidProposalContent( "proposal description cannot be blank")
	}

	if len(pcp.Description) > govtypes.MaxDescriptionLength {
		return ErrInvalidProposalContent(
			fmt.Sprintf("proposal description is longer than max length of %d", govtypes.MaxDescriptionLength))
	}

	if pcp.ProposalType() != sdkparams.ProposalTypeChange {
		return ErrInvalidProposalType( pcp.ProposalType())
	}

	if len(pcp.Changes) != 1 {
		return ErrInvalidMaxProposalNum(fmt.Sprintf("one proposal can only change one pair of parameter"))
	}

	return sdkparams.ValidateChanges(pcp.Changes)
}
