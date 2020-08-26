package token

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tokenchain/ixo-blockchain/client/utils"
	"github.com/tokenchain/ixo-blockchain/x/common/perf"
	"github.com/tokenchain/ixo-blockchain/x/common/version"
	"github.com/tokenchain/ixo-blockchain/x/dex"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/token/keeper"
	"github.com/tokenchain/ixo-blockchain/x/token/types"
)

// NewTokenHandler returns a handler for "token" type messages.
func NewTokenHandler(keeper keeper.Keeper, protocolVersion version.ProtocolVersionType) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		//logger := ctx.Logger().With("module", "token")
		// NOTE msg already has validate basic run
		var name string
		var handlerFun func() (*sdk.Result, error)
		logger := ctx.Logger().With("module", "token")
		switch msg := msg.(type) {
		case types.MsgTokenIssue:
			name = "handleMsgTokenIssue"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenIssue(ctx, keeper, msg, logger)
			}

		case types.MsgTokenBurn:
			name = "handleMsgTokenBurn"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenBurn(ctx, keeper, msg, logger)
			}

		case types.MsgTokenMint:
			name = "handleMsgTokenMint"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenMint(ctx, keeper, msg, logger)
			}

		case types.MsgMultiSend:
			name = "handleMsgMultiSend"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgMultiSend(ctx, keeper, msg, logger)
			}

		case types.MsgSend:
			name = "handleMsgSend"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgSend(ctx, keeper, msg, logger)
			}

		case types.MsgTransferOwnership:
			name = "handleMsgTokenChown"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenChown(ctx, keeper, msg, logger)
			}

		case types.MsgTokenModify:
			name = "handleMsgTokenModify"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTokenModify(ctx, keeper, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("Unrecognized token Msg type: %v", msg.Type())
			return &sdk.Result{}, did.UnknownRequest(errMsg)
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgTokenIssue(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTokenIssue, logger log.Logger) (*sdk.Result, error) {
	// check upper bound
	totalSupply, err := sdk.NewDecFromStr(msg.TotalSupply)
	if err != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("invalid total supply(%s)", msg.TotalSupply))
	}
	if totalSupply.GT(sdk.NewDec(types.TotalSupplyUpperbound)) {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("total-supply(%s) exceeds the upper limit(%d)",
			msg.TotalSupply, types.TotalSupplyUpperbound))
	}

	token := types.Token{
		Description:         msg.Description,
		OriginalSymbol:      msg.OriginalSymbol,
		WholeName:           msg.WholeName,
		OriginalTotalSupply: totalSupply,
		Owner:               msg.Owner,
		Mintable:            msg.Mintable,
	}

	// generate a random symbol
	newName, valid := keeper.AddTokenSuffix(ctx, keeper, msg.OriginalSymbol)
	if !valid {
		return &sdk.Result{}, did.ErrInvalidCoins(fmt.Sprintf(
			"temporarily failed to generate a unique symbol for %s. Try again.",
			msg.OriginalSymbol))
	}

	token.Symbol = newName
	totalsupply, ok := sdk.NewIntFromString(msg.TotalSupply)
	coins := sdk.NewCoin(token.Symbol, totalsupply)
	// set supply
	err = keeper.SupplyKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{coins})
	if err != nil || !ok {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("supply mint coins error:%s", err.Error()))
	}

	// send coins to owner
	err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, token.Owner, sdk.Coins{coins})
	if err != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("supply send coins error:%s", err.Error()))
	}

	// set token info
	keeper.NewToken(ctx, token)

	// deduction fee
	t := keeper.GetParams(ctx).FeeIssue
	feeDecCoins := utils.ParseDecCoinRounded(sdk.DecCoins{t})
	err = keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, token.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String()))
	}

	var name = "handleMsgTokenIssue"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Description:%s,Symbol:%s,OriginalSymbol:%s,TotalSupply:%s,Owner:%v,Mintable:%v>\n"+
			"                           result<Owner have enough okts to issue %s>\n",
			ctx.BlockHeight(), name,
			msg.Description, msg.Symbol, msg.OriginalSymbol, msg.TotalSupply, msg.Owner, msg.Mintable,
			token.Symbol))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(dex.AttributeKeyChargedFees, keeper.GetParams(ctx).FeeIssue.String()),
			sdk.NewAttribute("symbol", token.Symbol),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenBurn(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTokenBurn, logger log.Logger) (*sdk.Result, error) {

	token := keeper.GetTokenInfo(ctx, msg.Amount.Denom)

	// check owner
	if !token.Owner.Equals(msg.Owner) {
		return &sdk.Result{}, did.Unauthorized("Not the token's owner")
	}



	t :=  msg.Amount
	subCoins := utils.ParseDecCoinRounded(sdk.DecCoins{t})



	// send coins to moduleAcc
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, types.ModuleName, subCoins)
	if err != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("supply send coins error:%s", err.Error()))
	}

	// set supply
	err = keeper.SupplyKeeper.BurnCoins(ctx, types.ModuleName, subCoins)
	if err != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("supply burn coins error:%s", err.Error()))
	}

	// deduction fee

	f := keeper.GetParams(ctx).FeeBurn
	feeBurnCoins := utils.ParseDecCoinRounded(sdk.DecCoins{f})


	err = keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.feeCollectorName, feeBurnCoins)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeBurnCoins.String()))
	}

	var name = "handleMsgTokenBurn"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Owner:%s,Symbol:%s,Amount:%s>\n"+
			"                           result<Owner have enough okts to burn %s>\n",
			ctx.BlockHeight(), name,
			msg.Owner, msg.Amount.Denom, msg.Amount,
			msg.Amount.Denom))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(dex.AttributeKeyChargedFees, feeBurnCoins.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenMint(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTokenMint, logger log.Logger) (*sdk.Result, error) {
	token := keeper.GetTokenInfo(ctx, msg.Amount.Denom)
	// check owner
	if !token.Owner.Equals(msg.Owner) {
		return &sdk.Result{}, did.Unauthorized(fmt.Sprintf("%s is not the owner of token(%s)",
			msg.Owner.String(), msg.Amount.Denom))
	}

	// check whether token is mintable
	if !token.Mintable {
		return &sdk.Result{}, did.Unauthorized(fmt.Sprintf("token(%s) is not mintable", token.Symbol))
	}

	mintCoins := utils.ParseDecCoinRounded(sdk.DecCoins{msg.Amount})
	// set supply
	err := keeper.SupplyKeeper.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("supply mint coins error:%s", err.Error()))
	}

	// send coins to acc
	err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Owner, mintCoins)
	if err != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("supply send coins error:%s", err.Error()))
	}

	// deduction fee
	feeDecCoins := utils.ParseDecCoinRounded(sdk.DecCoins{keeper.GetParams(ctx).FeeMint})
	err = keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String()))
	}

	name := "handleMsgTokenMint"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Owner:%s,Symbol:%s,Amount:%s>\n"+
			"                           result<Owner have enough okts to Mint %s>\n",
			ctx.BlockHeight(), name,
			msg.Owner, msg.Amount.Denom, msg.Amount,
			msg.Amount.Denom))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(dex.AttributeKeyChargedFees, feeDecCoins.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMultiSend(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgMultiSend, logger log.Logger) (*sdk.Result, error) {
	if !keeper.bankKeeper.GetSendEnabled(ctx) {
		return &sdk.Result{}, types.SendDisabled()
	}

	var transfers string
	var coinNum int
	for _, transferUnit := range msg.Transfers {
		coinNum += len(transferUnit.Coins)
		err := keeper.SendCoinsFromAccountToAccount(ctx, msg.From, transferUnit.To, transferUnit.Coins)
		if err != nil {
			return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient coins(need %s)",
				transferUnit.Coins.String()))
		}
		transfers += fmt.Sprintf("                          msg<To:%s,Coin:%s>\n", transferUnit.To, transferUnit.Coins)
	}

	name := "handleMsgMultiSend"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<From:%s>\n"+
			transfers+
			"                           result<Owner have enough okts to send multi txs>\n",
			ctx.BlockHeight(), name,
			msg.From))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName)),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSend(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSend, logger log.Logger) (*sdk.Result, error) {
	if !keeper.bankKeeper.GetSendEnabled(ctx) {
		return &sdk.Result{}, types.SendDisabled()
	}

	err := keeper.SendCoinsFromAccountToAccount(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient coins(need %s)",
			msg.Amount.String()))
	}

	var name = "handleMsgSend"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<From:%s,To:%s,Amount:%s>\n"+
			"                           result<Owner have enough okts to send a tx>\n",
			ctx.BlockHeight(), name,
			msg.FromAddress, msg.ToAddress, msg.Amount))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName)),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenChown(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTransferOwnership, logger log.Logger) (*sdk.Result, error) {
	tokenInfo := keeper.GetTokenInfo(ctx, msg.Symbol)

	if !tokenInfo.Owner.Equals(msg.FromAddress) {
		return &sdk.Result{}, did.Unauthorized(fmt.Sprintf("%s is not the owner of token(%s)",
			msg.FromAddress.String(), msg.Symbol))
	}

	// first remove it from the raw owner
	keeper.DeleteUserToken(ctx, tokenInfo.Owner, tokenInfo.Symbol)

	tokenInfo.Owner = msg.ToAddress
	keeper.NewToken(ctx, tokenInfo)

	// deduction fee
	t :=  keeper.GetParams(ctx).FeeChown
	feeDecCoins := utils.ParseDecCoinRounded(sdk.DecCoins{t})


	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FromAddress, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String()))
	}

	var name = "handleMsgTokenChown"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<From:%s,To:%s,Symbol:%s,ToSign:%s>\n"+
			"                           result<Owner have enough okts to transfer the %s>\n",
			ctx.BlockHeight(), name,
			msg.FromAddress, msg.ToAddress, msg.Symbol, msg.ToSignature,
			msg.Symbol))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(dex.AttributeKeyChargedFees, keeper.GetParams(ctx).FeeChown.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTokenModify(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgTokenModify, logger log.Logger) (*sdk.Result, error) {
	token := keeper.GetTokenInfo(ctx, msg.Symbol)
	// check owner
	if !token.Owner.Equals(msg.Owner) {
		return &sdk.Result{}, did.Unauthorized(fmt.Sprintf("%s is not the owner of token(%s)",
			msg.Owner.String(), msg.Symbol))
	}
	if !msg.IsWholeNameModified && !msg.IsDescriptionModified {
		return &sdk.Result{}, did.IntErr("nothing modified")
	}
	// modify
	if msg.IsWholeNameModified {
		token.WholeName = msg.WholeName
	}
	if msg.IsDescriptionModified {
		token.Description = msg.Description
	}

	keeper.UpdateToken(ctx, token)

	// deduction fee
	t :=  keeper.GetParams(ctx).FeeModify
	feeDecCoins := utils.ParseDecCoinRounded(sdk.DecCoins{t})

	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.feeCollectorName, feeDecCoins)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeDecCoins.String()))
	}

	name := "handleMsgTokenModify"
	if logger != nil {
		logger.Debug(fmt.Sprintf("BlockHeight<%d>, handler<%s>\n"+
			"                           msg<Owner:%s,Symbol:%s,WholeName:%s,Description:%s>\n"+
			"                           result<Owner have enough okts to edit %s>\n",
			ctx.BlockHeight(), name,
			msg.Owner, msg.Symbol, msg.WholeName, msg.Description,
			msg.Symbol))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(dex.AttributeKeyChargedFees, keeper.GetParams(ctx).FeeModify.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
