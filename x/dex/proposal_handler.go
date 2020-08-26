package dex

import (
	"fmt"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/tokenchain/ixo-blockchain/x/dex/types"
)

// NewProposalHandler handles "gov" type message in "dex"
func NewProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, content govTypes.Content) error {
		switch c := content.(type) {
		case types.DelistProposal:
			return handleDelistProposal(ctx, k, c)
		default:
			errMsg := fmt.Sprintf("unrecognized param proposal content type: %s", c)
			return did.UnknownRequest(errMsg)
		}
	}
}

func handleDelistProposal(ctx sdk.Context, keeper *Keeper, p types.DelistProposal) error {
	logger := ctx.Logger().With("module", types.ModuleName)
	logger.Debug("execute DelistProposal begin")

	tokenPairName := fmt.Sprintf("%s_%s", p.BaseAsset, p.QuoteAsset)
	tokenPair := keeper.GetTokenPair(ctx, tokenPairName)
	if tokenPair == nil {
		return ErrTokenPairNotFound(fmt.Sprintf("%+v", p))
	}
	if keeper.IsTokenPairLocked(ctx, tokenPairName) {
		errContent := fmt.Sprintf("unexpected state, the trading pair (%s) is locked", tokenPairName)
		return did.IntErr(errContent)
	}

	// withdraw
	if tokenPair.Deposits.IsPositive() {
		err := keeper.Withdraw(ctx, tokenPair.Name(), tokenPair.Owner, tokenPair.Deposits)
		if err != nil {
			return did.IntErr(fmt.Sprintf("failed to withdraw deposits:%s error: %V",
				tokenPair.Deposits.String(), err))
		}
	}

	// delete the token pair by its name from store and cache
	keeper.DeleteTokenPairByName(ctx, tokenPair.Owner, tokenPairName)
	// remove the delistProposal from the active proposal queue
	keeper.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
	//= ====== ====== ======== ======
	ctx.EventManager().EmitEvent(sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute("token-pair-deleted", tokenPairName),
			))

	return nil
}
