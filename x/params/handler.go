package params

import (
	"fmt"
	"github.com/tokenchain/ixo-blockchain/client/utils"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"math"
	"time"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/tokenchain/ixo-blockchain/x/common"
	"github.com/tokenchain/ixo-blockchain/x/params/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParamChangeProposalHandler returns the rollback function of the param proposal handler
func NewParamChangeProposalHandler(k *Keeper) govtypes.Handler {
	return func(ctx sdk.Context, proposalContent govtypes.Content) error {
		switch c := proposalContent.(type) {

		case govtypes.Proposal:
			return handleParameterChangeProposal(ctx, k, &c)

		default:
			errMsg := fmt.Sprintf("unrecognized param proposal content type: %T", c)
			return did.UnknownRequest(errMsg)
		}
	}
}

func handleParameterChangeProposal(ctx sdk.Context, k *Keeper, proposal *govtypes.Proposal) error {
	logger := ctx.Logger().With("module", ModuleName)
	logger.Info("Execute ParameterProposal begin")
	paramProposal := proposal.Content.(types.ParameterChangeProposal)
	curHeight := uint64(ctx.BlockHeight())
	if paramProposal.Height > curHeight {
		k.gk.InsertWaitingProposalQueue(ctx, paramProposal.Height, proposal.ProposalID)
		return nil
	}

	defer k.gk.RemoveFromWaitingProposalQueue(ctx, paramProposal.Height, proposal.ProposalID)
	return changeParams(ctx, k, paramProposal)
}

func changeParams(ctx sdk.Context, k *Keeper, paramProposal types.ParameterChangeProposal) error {
	for _, c := range paramProposal.Changes {
		ss, ok := k.Keeper.GetSubspace(c.Subspace)
		if !ok {
			return types.UnknownSubspace("")
		}
		err := ss.Update(ctx, []byte(c.Key), []byte(c.Value))
		if err != nil {
			return types.ErrSettingParam(c.Key, c.Value, err.Error())
		}
	/*

	ss, ok := k.Keeper.GetSubspace(c.Subspace)
			if !ok {
				return types.UnknownSubspace("")
			}
	var err error
		if len(c.Subkey) == 0 {
			k.Logger(ctx).Info(
				fmt.Sprintf("setting new parameter; key: %s, value: %s", c.Key, c.Value),
			)

			err = ss.Update(ctx, []byte(c.Key), []byte(c.Value))
		} else {
			k.Logger(ctx).Info(
				fmt.Sprintf("setting new parameter; key: %s, subkey: %s, value: %s", c.Key, c.Subspace, c.Value),
			)
			err = ss.UpdateWithSubkey(ctx, []byte(c.Key), []byte(c.Subkey), []byte(c.Value))
		}

		if err != nil {
			return types.ErrSettingParam(c.Key, c.Value, err.Error())
		}*/
	}
	return nil
}

// GetMinDeposit implements ProposalHandler interface
func (keeper Keeper) GetMinDeposit(ctx sdk.Context, content govtypes.Content) (minDeposit sdk.DecCoins) {
	switch content.(type) {
	case types.ParameterChangeProposal:
		minDeposit = keeper.GetParams(ctx).MinDeposit
	}
	return
}

// GetMaxDepositPeriod implements ProposalHandler interface
func (keeper Keeper) GetMaxDepositPeriod(ctx sdk.Context, content govtypes.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.ParameterChangeProposal:
		maxDepositPeriod = keeper.GetParams(ctx).MaxDepositPeriod
	}
	return
}

// GetVotingPeriod implements ProposalHandler interface
func (keeper Keeper) GetVotingPeriod(ctx sdk.Context, content govtypes.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.ParameterChangeProposal:
		votingPeriod = keeper.GetParams(ctx).VotingPeriod
	}
	return
}

// CheckMsgSubmitProposal implements ProposalHandler interface
func (keeper Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govtypes.MsgSubmitProposal) error {
	paramsChangeProposal := msg.Content.(types.ParameterChangeProposal)
	// check message sender is current validator
	if !keeper.sk.IsValidator(ctx, msg.Proposer) {
		return types.InvalidProposer("proposer of ParamChange proposal must be validator")
	}
	// check initial deposit more than or equal to ratio of MinDeposit
	initDeposit := keeper.GetParams(ctx).MinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	init_deposit := utils.ParseDecCoinRounded(initDeposit)
	if err := common.HasSufficientCoins(msg.Proposer, msg.InitialDeposit, init_deposit); err != nil {
		return did.ErrInvalidCoins(fmt.Sprintf("InitialDeposit must not be less than %s", initDeposit.String()))
	}
	// check proposer has sufficient coins
	if err := common.HasSufficientCoins(msg.Proposer, keeper.ck.GetCoins(ctx, msg.Proposer), msg.InitialDeposit); err != nil {
		return did.ErrInvalidCoins(err.Error())
	}

	curHeight := uint64(ctx.BlockHeight())
	maxHeight := keeper.GetParams(ctx).MaxBlockHeight
	if maxHeight == 0 {
		maxHeight = math.MaxInt64 - paramsChangeProposal.Height
	}
	if paramsChangeProposal.Height < curHeight || paramsChangeProposal.Height > curHeight+maxHeight {
		return types.InvalidHeight(paramsChangeProposal.Height, curHeight, maxHeight)
	}
	// run simulation with cache context
	cacheCtx, _ := ctx.CacheContext()
	return changeParams(cacheCtx, &keeper, paramsChangeProposal)
}

// nolint
func (keeper Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govtypes.Proposal) {
}
func (keeper Keeper) VoteHandler(ctx sdk.Context, proposal govtypes.Proposal, vote govtypes.Vote) (string, error) {
	return "", nil
}
func (keeper Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govtypes.Proposal) {
}
func (keeper Keeper) RejectedHandler(ctx sdk.Context, content govtypes.Content) {
}
