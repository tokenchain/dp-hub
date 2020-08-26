package keeper

import (
	"fmt"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/tokenchain/ixo-blockchain/x/common"
	"github.com/tokenchain/ixo-blockchain/x/dex/types"
)

// GetMinDeposit returns min deposit
func (k Keeper) GetMinDeposit(ctx sdk.Context, content gov.Content) (minDeposit sdk.DecCoins) {
	if _, ok := content.(types.DelistProposal); ok {
		minDeposit = k.GetParams(ctx).DelistMinDeposit
	}
	return
}

// GetMaxDepositPeriod returns max deposit period
func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content gov.Content) (maxDepositPeriod time.Duration) {
	if _, ok := content.(types.DelistProposal); ok {
		maxDepositPeriod = k.GetParams(ctx).DelistMaxDepositPeriod
	}
	return
}

// GetVotingPeriod returns voting period
func (k Keeper) GetVotingPeriod(ctx sdk.Context, content gov.Content) (votingPeriod time.Duration) {
	if _, ok := content.(types.DelistProposal); ok {
		votingPeriod = k.GetParams(ctx).DelistVotingPeriod
	}
	return
}

// check msg Delist proposal
func (k Keeper) checkMsgDelistProposal(ctx sdk.Context, delistProposal types.DelistProposal, proposer sdk.AccAddress, initialDeposit sdk.DecCoins) error {
	// check the proposer of the msg is a validator
	if !k.stakingKeeper.IsValidator(ctx, proposer) {
		return types.ErrInvalidProposer("failed to submit proposal because the proposer of delist proposal should be a validator")
	}

	// check the propose of the msg is equal the proposer in proposal content
	if !proposer.Equals(delistProposal.Proposer) {
		return types.ErrInvalidProposer("failed to submit proposal because the proposer of proposal msg should be equal the proposer in proposal content")
	}

	// check whether the baseAsset is in the Dex list
	queryTokenPair := k.GetTokenPair(ctx, fmt.Sprintf("%s_%s", delistProposal.BaseAsset, delistProposal.QuoteAsset))
	if queryTokenPair == nil {
		return types.ErrTokenPairNotFound(fmt.Sprintf("failed to submit proposal because the asset with base asset '%s' and quote asset '%s' didn't exist on the Dex", delistProposal.BaseAsset, delistProposal.QuoteAsset))
	}

	// check the initial deposit
	localMinDeposit := k.GetParams(ctx).DelistMinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	err := common.HasSufficientDecCoins(proposer, initialDeposit, localMinDeposit)

	if err != nil {
		return types.ErrInvalidAsset(fmt.Sprintf("failed to submit proposal because initial deposit should be more than %s", localMinDeposit.String()))
	}

	// check whether the proposer can afford the initial deposit
	coinProposer := sdk.NewDecCoinsFromCoins(k.bankKeeper.GetCoins(ctx, proposer)...)
	err = common.HasSufficientDecCoins(proposer, coinProposer, initialDeposit)
	if err != nil {
		return types.ErrInvalidBalanceNotEnough(fmt.Sprintf("failed to submit proposal because proposer %s didn't have enough coins to pay for the initial deposit %s", proposer, initialDeposit))
	}
	return nil
}

// CheckMsgSubmitProposal validates MsgSubmitProposal
func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) (sdkErr error) {
	switch content := msg.Content.(type) {
	case types.DelistProposal:
		decCoins := sdk.NewDecCoinsFromCoins(msg.InitialDeposit...)
		sdkErr = k.checkMsgDelistProposal(ctx, content, msg.Proposer, decCoins)
	default:
		errContent := fmt.Sprintf("unrecognized dex proposal content type: %T", content)
		sdkErr = did.UnknownRequest(errContent)
	}
	return
}

// nolint
func (k Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govTypes.Proposal) {}

// VoteHandler handles  delist proposal when voted
func (k Keeper) VoteHandler(ctx sdk.Context, proposal govTypes.Proposal, vote govTypes.Vote) (string, error) {
	if _, ok := proposal.Content.(types.DelistProposal); ok {
		delistProposal := proposal.Content.(types.DelistProposal)
		tokenPairName := delistProposal.BaseAsset + "_" + delistProposal.QuoteAsset
		if k.IsTokenPairLocked(ctx, tokenPairName) {
			errContent := fmt.Sprintf("the trading pair (%s) is locked, please retry later", tokenPairName)
			return "", did.IntErr(errContent)
		}
	}
	return "", nil
}

// RejectedHandler handles delist proposal when rejected
func (k Keeper) RejectedHandler(ctx sdk.Context, content govTypes.Content) {
	if content, ok := content.(types.DelistProposal); ok {
		tokenPairName := fmt.Sprintf("%s_%s", content.BaseAsset, content.QuoteAsset)
		//update the token info from the store
		tokenPair := k.GetTokenPair(ctx, tokenPairName)
		if tokenPair == nil {
			ctx.Logger().Error(fmt.Sprintf("token pair %s does not exist", tokenPairName))
			return
		}
		tokenPair.Delisting = false
		k.UpdateTokenPair(ctx, tokenPairName, tokenPair)
	}
}

// AfterDepositPeriodPassed handles delist proposal when passed
func (k Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govTypes.Proposal) {
	if content, ok := proposal.Content.(types.DelistProposal); ok {
		tokenPairName := fmt.Sprintf("%s_%s", content.BaseAsset, content.QuoteAsset)
		// change the status of the token pair in the store
		tokenPair := k.GetTokenPair(ctx, tokenPairName)
		if tokenPair == nil {
			ctx.Logger().Error(fmt.Sprintf("token pair %s does not exist", tokenPairName))
			return
		}
		tokenPair.Delisting = true
		k.UpdateTokenPair(ctx, tokenPairName, tokenPair)
	}
}

// RemoveFromActiveProposalQueue removes active proposal in queue
func (k Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	k.govKeeper.RemoveFromActiveProposalQueue(ctx, proposalID, endTime)
}
