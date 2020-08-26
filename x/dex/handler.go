package dex

import (
	"fmt"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"strconv"

	"github.com/tokenchain/ixo-blockchain/x/common/perf"
	"github.com/tokenchain/ixo-blockchain/x/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
)

// NewHandler handles all "dex" type messages.
func NewHandler(k IKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", ModuleName)

		var handlerFun func() (*sdk.Result, error)
		var name string
		switch msg := msg.(type) {
		case MsgList:
			name = "handleMsgList"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgList(ctx, k, msg, logger)
			}
		case MsgDeposit:
			name = "handleMsgDeposit"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgDeposit(ctx, k, msg, logger)
			}
		case MsgWithdraw:
			name = "handleMsgWithDraw"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgWithDraw(ctx, k, msg, logger)
			}
		case MsgTransferOwnership:
			name = "handleMsgTransferOwnership"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgTransferOwnership(ctx, k, msg, logger)
			}
		case MsgCreateOperator:
			name = "handleMsgCreateOperator"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgCreateOperator(ctx, k, msg, logger)
			}
		case MsgUpdateOperator:
			name = "handleMsgUpdateOperator"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgUpdateOperator(ctx, k, msg, logger)
			}
		default:
			errMsg := fmt.Sprintf("unrecognized dex message type: %T", msg)
			return &sdk.Result{}, did.UnknownRequest(errMsg)
		}

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, ModuleName, name, seq)
		return handlerFun()
	}
}

func handleMsgList(ctx sdk.Context, keeper IKeeper, msg MsgList, logger log.Logger) (*sdk.Result, error) {

	if !keeper.GetTokenKeeper().TokenExist(ctx, msg.ListAsset) ||
		!keeper.GetTokenKeeper().TokenExist(ctx, msg.QuoteAsset) {
		return &sdk.Result{}, did.ErrInvalidCoins(fmt.Sprintf("%s or %s is not valid", msg.ListAsset, msg.QuoteAsset))
	}

	if _, exists := keeper.GetOperator(ctx, msg.Owner); !exists {
		return &sdk.Result{}, types.ErrUnknownOperator(msg.Owner)
	}

	tokenPair := &TokenPair{
		BaseAssetSymbol:  msg.ListAsset,
		QuoteAssetSymbol: msg.QuoteAsset,
		InitPrice:        msg.InitPrice,
		MaxPriceDigit:    int64(DefaultMaxPriceDigitSize),
		MaxQuantityDigit: int64(DefaultMaxQuantityDigitSize),
		MinQuantity:      sdk.MustNewDecFromStr("0.00000001"),
		Owner:            msg.Owner,
		Delisting:        false,
		Deposits:         DefaultTokenPairDeposit,
		BlockHeight:      ctx.BlockHeight(),
	}

	// check whether a specific token pair exists with the symbols of base asset and quote asset
	// Note: aaa_bbb and bbb_aaa are actually one token pair
	if keeper.GetTokenPair(ctx, fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)) != nil ||
		keeper.GetTokenPair(ctx, fmt.Sprintf("%s_%s", tokenPair.QuoteAssetSymbol, tokenPair.BaseAssetSymbol)) != nil {
		return &sdk.Result{}, types.ErrTokenPairExisted(tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	}

	// deduction fee
	feeCoins, _ := keeper.GetParams(ctx).ListFee.TruncateDecimal()
	err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.GetFeeCollector(), sdk.Coins{feeCoins})
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeCoins))
	}

	err2 := keeper.SaveTokenPair(ctx, tokenPair)
	if err2 != nil {
		return &sdk.Result{}, did.IntErr(fmt.Sprintf("failed to SaveTokenPair: %s", err2.Error()))
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgList: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(AttributeKeyListAsset, tokenPair.BaseAssetSymbol),
			sdk.NewAttribute(AttributeKeyQuoteAsset, tokenPair.QuoteAssetSymbol),
			sdk.NewAttribute(AttributeKeyInitPrice, tokenPair.InitPrice.String()),
			sdk.NewAttribute(AttributeKeyMaxPriceDigit, strconv.FormatInt(tokenPair.MaxPriceDigit, 10)),
			sdk.NewAttribute(AttributeKeyMaxSizeDigit, strconv.FormatInt(tokenPair.MaxQuantityDigit, 10)),
			sdk.NewAttribute(AttributeKeyMinTradeFee, tokenPair.MinQuantity.String()),
			sdk.NewAttribute(AttributeKeyDelisting, fmt.Sprintf("%t", tokenPair.Delisting)),
			sdk.NewAttribute(AttributeKeyChargedFees, feeCoins.String()),
		),
	)

	return &sdk.Result{}, nil
}

func handleMsgDeposit(ctx sdk.Context, keeper IKeeper, msg MsgDeposit, logger log.Logger) (*sdk.Result, error) {
	if sdkErr := keeper.Deposit(ctx, msg.Product, msg.Depositor, msg.Amount); sdkErr != nil {
		return &sdk.Result{}, sdkErr
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgDeposit: BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return &sdk.Result{}, nil
}

func handleMsgWithDraw(ctx sdk.Context, keeper IKeeper, msg MsgWithdraw, logger log.Logger) (*sdk.Result, error) {
	if sdkErr := keeper.Withdraw(ctx, msg.Product, msg.Depositor, msg.Amount); sdkErr != nil {
		return &sdk.Result{}, sdkErr
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgWithDraw: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTransferOwnership(ctx sdk.Context, keeper IKeeper, msg MsgTransferOwnership,
	logger log.Logger) (*sdk.Result, error) {

	if _, exist := keeper.GetOperator(ctx, msg.FromAddress); !exist {
		return &sdk.Result{}, types.ErrUnknownOperator(msg.FromAddress)
	}

	if _, exist := keeper.GetOperator(ctx, msg.ToAddress); !exist {
		return &sdk.Result{}, types.ErrUnknownOperator(msg.ToAddress)
	}

	if sdkErr := keeper.TransferOwnership(ctx, msg.Product, msg.FromAddress, msg.ToAddress); sdkErr != nil {
		return &sdk.Result{}, sdkErr
	}

	// deduction fee
	feeCoins, _ := keeper.GetParams(ctx).TransferOwnershipFee.TruncateDecimal()
	err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.FromAddress, keeper.GetFeeCollector(), sdk.Coins{feeCoins})
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeCoins))
	}

	logger.Debug(fmt.Sprintf("successfully handleMsgTransferOwnership: "+
		"BlockHeight: %d, Msg: %+v", ctx.BlockHeight(), msg))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(AttributeKeyChargedFees, feeCoins.String()),
		),
	)
	return &sdk.Result{}, nil
}

func handleMsgCreateOperator(ctx sdk.Context, keeper IKeeper, msg MsgCreateOperator, logger log.Logger) (*sdk.Result, error) {

	logger.Debug(fmt.Sprintf("handleMsgCreateOperator msg: %+v", msg))

	if _, isExist := keeper.GetOperator(ctx, msg.Owner); isExist {
		return &sdk.Result{}, types.ErrExistOperator(msg.Owner)
	}
	operator := types.DEXOperator{
		Address:            msg.Owner,
		HandlingFeeAddress: msg.HandlingFeeAddress,
		Website:            msg.Website,
		InitHeight:         ctx.BlockHeight(),
		TxHash:             fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes())),
	}
	keeper.SetOperator(ctx, operator)

	// deduction fee
	feeCoins, _ := sdk.DecCoins{keeper.GetParams(ctx).RegisterOperatorFee}.TruncateDecimal()
	err := keeper.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Owner, keeper.GetFeeCollector(), feeCoins)
	if err != nil {
		return &sdk.Result{}, did.InsufficientCoins(fmt.Sprintf("insufficient fee coins(need %s)",
			feeCoins.String()))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(AttributeKeyChargedFees, feeCoins.String()),
		),
	)

	return &sdk.Result{}, nil
}

func handleMsgUpdateOperator(ctx sdk.Context, keeper IKeeper, msg MsgUpdateOperator, logger log.Logger) (*sdk.Result, error) {

	logger.Debug(fmt.Sprintf("handleMsgUpdateOperator msg: %+v", msg))

	operator, isExist := keeper.GetOperator(ctx, msg.Owner)
	if !isExist {
		return &sdk.Result{}, types.ErrUnknownOperator(msg.Owner)
	}
	if !operator.Address.Equals(msg.Owner) {
		return &sdk.Result{}, did.Unauthorized("Not the operator's owner")
	}

	operator.HandlingFeeAddress = msg.HandlingFeeAddress
	operator.Website = msg.Website

	keeper.SetOperator(ctx, operator)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
		),
	)

	return &sdk.Result{}, nil
}
