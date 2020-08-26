package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	proposalTypeDelist = "Delist"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeDelist)
	govtypes.RegisterProposalTypeCodec(DelistProposal{}, "okchain/dex/DelistProposal")

}

// Assert DelistProposal implements govtypes.Content at compile-time
var _ govtypes.Content = (*DelistProposal)(nil)

// DelistProposal represents delist proposal object
type DelistProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	BaseAsset   string         `json:"base_asset" yaml:"base_asset"`
	QuoteAsset  string         `json:"quote_asset" yaml:"quote_asset"`
}

// NewDelistProposal create a new delist proposal object
func NewDelistProposal(title, description string, proposer sdk.AccAddress, baseAsset, quoteAsset string) DelistProposal {
	return DelistProposal{
		Title:       title,
		Description: description,
		Proposer:    proposer,
		BaseAsset:   baseAsset,
		QuoteAsset:  quoteAsset,
	}
}

// GetTitle returns title of delist proposal object
func (drp DelistProposal) GetTitle() string {
	return drp.Title
}

// GetDescription returns description of delist proposal object
func (drp DelistProposal) GetDescription() string {
	return drp.Description
}

// ProposalRoute returns route key of delist proposal object
func (DelistProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of delist proposal object
func (DelistProposal) ProposalType() string {
	return proposalTypeDelist
}

// ValidateBasic validates delist proposal
func (drp DelistProposal) ValidateBasic() error {
	if len(strings.TrimSpace(drp.Title)) == 0 {
		return ErrInvalidProposalContent("failed to submit delist proposal because title is blank")
	}
	if len(drp.Title) > govtypes.MaxTitleLength {
		return ErrInvalidProposalContent(fmt.Sprintf("failed to submit delist proposal because title is longer than max length of %d", govtypes.MaxTitleLength))
	}

	if len(drp.Description) == 0 {
		return ErrInvalidProposalContent("failed to submit delist proposal because  description is blank")
	}

	if len(drp.Description) > govtypes.MaxDescriptionLength {
		return ErrInvalidProposalContent(fmt.Sprintf("failed to submit delist proposal because  description is longer than max length of %d", govtypes.MaxDescriptionLength))
	}

	if drp.ProposalType() != proposalTypeDelist {
		return ErrInvalidProposalType(drp.ProposalType())
	}

	if drp.Proposer.Empty() {
		return ErrInvalidAddress(drp.Proposer.String())
	}

	if drp.BaseAsset == drp.QuoteAsset {
		return ErrInvalidCoins(fmt.Sprintf("failed to submit delist proposal because baseasset is same as quoteasset"))
	}

	return nil
}

// String converts delist proposal object to string
func (drp DelistProposal) String() string {
	return fmt.Sprintf(`DelistProposal:
 Title:               %s
 Description:         %s
 Type:                %s
 Proposer:            %s
 ListAsset            %s
 QuoteAsset           %s
`, drp.Title, drp.Description,
		drp.ProposalType(), drp.Proposer,
		drp.BaseAsset, drp.QuoteAsset,
	)
}
