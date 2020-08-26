// nolint
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/tokenchain/ixo-blockchain/x/dex/keeper
// ALIASGEN: github.com/tokenchain/ixo-blockchain/x/dex/types
package dex

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tokenchain/ixo-blockchain/x/common/version"
	"github.com/tokenchain/ixo-blockchain/x/dex/keeper"
	"github.com/tokenchain/ixo-blockchain/x/dex/types"
)

const (
	ModuleName                  = types.ModuleName
	DefaultCodespace            = types.DefaultCodespace
	DefaultParamspace           = types.DefaultParamspace
	TokenPairStoreKey           = types.TokenPairStoreKey
	QuerierRoute                = types.QuerierRoute
	RouterKey                   = types.RouterKey
	StoreKey                    = types.StoreKey
	DefaultMaxPriceDigitSize    = types.DefaultMaxPriceDigitSize
	DefaultMaxQuantityDigitSize = types.DefaultMaxQuantityDigitSize
	AuthFeeCollector            = auth.FeeCollectorName

	EventTypeCreateBond                = "create_bond"
	EventTypeEditBond                  = "edit_bond"
	EventTypeInitSwapper               = "init_swapper"
	EventTypeBuy                       = "buy"
	EventTypeSell                      = "sell"
	EventTypeSwap                      = "swap"
	EventTypeOrderCancel               = "order_cancel"
	EventTypeOrderFulfill              = "order_fulfill"
	AttributeKeyBondDid                = "bond_did"
	AttributeKeyToken                  = "token"
	AttributeKeyName                   = "name"
	AttributeKeyDescription            = "description"
	AttributeKeyFunctionType           = "function_type"
	AttributeKeyFunctionParameters     = "function_parameters"
	AttributeKeyReserveTokens          = "reserve_tokens"
	AttributeKeyReserveAddress         = "reserve_address"
	AttributeKeyTxFeePercentage        = "tx_fee_percentage"
	AttributeKeyExitFeePercentage      = "exit_fee_percentage"
	AttributeKeyFeeAddress             = "fee_address"
	AttributeKeyMaxSupply              = "max_supply"
	AttributeKeyOrderQuantityLimits    = "order_quantity_limits"
	AttributeKeySanityRate             = "sanity_rate"
	AttributeKeySanityMarginPercentage = "sanity_margin_percentage"
	AttributeKeyAllowSells             = "allow_sells"
	AttributeKeyBatchBlocks            = "batch_blocks"
	AttributeKeyMaxPrices              = "max_prices"
	AttributeKeySwapFromToken          = "from_token"
	AttributeKeySwapToToken            = "to_token"
	AttributeKeyOrderType              = "order_type"
	AttributeKeyAddress                = "address"
	AttributeKeyCancelReason           = "cancel_reason"
	AttributeKeyTokensMinted           = "tokens_minted"
	AttributeKeyTokensBurned           = "tokens_burned"
	AttributeKeyTokensSwapped          = "tokens_swapped"
	AttributeKeyChargedPrices          = "charged_prices"
	AttributeKeyChargedFees            = "tx-fee"
	AttributeKeyDelisting              = "delisting"
	AttributeKeyMinTradeFee            = "min-trade-size"
	AttributeKeyMaxSizeDigit           = "max-size-digit"
	AttributeKeyMaxPriceDigit          = "max-price-digit"
	AttributeKeyInitPrice              = "init-price"
	AttributeKeyQuoteAsset             = "quote-asset"
	AttributeKeyListAsset              = "list-asset"
	AttributeKeyReturnedToAddress      = "returned_to_address"
	AttributeKeyNewBondTokenBalance    = "new_bond_token_balance"
	AttributeValueBuyOrder             = "buy"
	AttributeValueSellOrder            = "sell"
	AttributeValueSwapOrder            = "swap"
	AttributeValueCategory             = ModuleName
)

type (
	// Keepers
	Keeper              = keeper.Keeper
	IKeeper             = keeper.IKeeper
	SupplyKeeper        = keeper.SupplyKeeper
	TokenKeeper         = keeper.TokenKeeper
	StakingKeeper       = keeper.StakingKeeper
	BankKeeper          = keeper.BankKeeper
	ProtocolVersionType = version.ProtocolVersionType
	// Messages
	MsgList              = types.MsgList
	MsgDeposit           = types.MsgDeposit
	MsgWithdraw          = types.MsgWithdraw
	MsgTransferOwnership = types.MsgTransferOwnership
	MsgUpdateOperator    = types.MsgUpdateOperator
	MsgCreateOperator    = types.MsgCreateOperator
	TokenPair            = types.TokenPair
	Params               = types.Params
	WithdrawInfo         = types.WithdrawInfo
	WithdrawInfos        = types.WithdrawInfos
	DEXOperator          = types.DEXOperator
	DEXOperators         = types.DEXOperators
)

var (
	ModuleCdc               = types.ModuleCdc
	DefaultTokenPairDeposit = types.DefaultTokenPairDeposit

	RegisterCodec       = types.RegisterCodec
	NewQuerier          = keeper.NewQuerier
	NewKeeper           = keeper.NewKeeper
	GetBuiltInTokenPair = keeper.GetBuiltInTokenPair
	DefaultParams       = types.DefaultParams

	NewMsgList     = types.NewMsgList
	NewMsgDeposit  = types.NewMsgDeposit
	NewMsgWithdraw = types.NewMsgWithdraw

	ErrInvalidProduct      = types.ErrInvalidProduct
	ErrTokenPairNotFound   = types.ErrTokenPairNotFound
	ErrDelistOwnerNotMatch = types.ErrDelistOwnerNotMatch
)
