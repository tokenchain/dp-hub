package fees

import (
	"github.com/tokenchain/ixo-blockchain/x/fees/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/fees/internal/types"
)

const (
	ModuleName = types.ModuleName

	FeeRemainderPool = types.FeeRemainderPool

	FeePrefix          = types.FeeIdPrefix
	FeeContractPrefix  = types.FeeContractIdPrefix
	SubscriptionPrefix = types.SubscriptionIdPrefix

	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey

	FeeClaimTransaction      = types.FeeClaimTransaction
	FeeEvaluationTransaction = types.FeeEvaluationTransaction
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState

	FeeType = types.FeeType

	Fee               = types.Fee
	FeeContract       = types.FeeContract
	Distribution      = types.Distribution
	DistributionShare = types.DistributionShare

	Discount  = types.Discount
	Discounts = types.Discounts

	Subscription = types.Subscription
	Period       = types.Period
	BlockPeriod  = types.BlockPeriod
	TimePeriod   = types.TimePeriod

	MsgSetFeeContractAuthorisation = types.MsgSetFeeContractAuthorisation
	MsgCreateFee                   = types.MsgCreateFee
	MsgCreateFeeContract           = types.MsgCreateFeeContract
	MsgCreateSubscription          = types.MsgCreateSubscription
	MsgGrantFeeDiscount            = types.MsgGrantFeeDiscount
	MsgRevokeFeeDiscount           = types.MsgRevokeFeeDiscount
	MsgChargeFee                   = types.MsgChargeFee
)

var (
	// function aliases
	NewKeeper      = keeper.NewKeeper
	NewQuerier     = keeper.NewQuerier
	ParamKeyTable  = types.ParamKeyTable
	NewParams      = types.NewParams
	DefaultParams  = types.DefaultParams
	ValidateParams = types.ValidateParams

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	NewFee                   = types.NewFee
	NewFeeContract           = types.NewFeeContract
	NewFeeContractNoDiscount = types.NewFeeContractNoDiscount
	NewDistribution          = types.NewDistribution
	NewDistributionShare     = types.NewDistributionShare

	NewDiscount  = types.NewDiscount
	NewDiscounts = types.NewDiscounts

	NewSubscription = types.NewSubscription
	NewBlockPeriod  = types.NewBlockPeriod
	NewTimePeriod   = types.NewTimePeriod

	// variable aliases
	ModuleCdc                               = types.ModuleCdc
	KeyIxoFactor                            = types.KeyIxoFactor
	KeyNodeFeePercentage                    = types.KeyNodeFeePercentage
	KeyClaimFeeAmount                       = types.KeyClaimFeeAmount
	KeyEvaluationFeeAmount                  = types.KeyEvaluationFeeAmount
	KeyInitiationFeeAmount                  = types.KeyInitiationFeeAmount
	KeyInitiationNodeFeePercentage          = types.KeyInitiationNodeFeePercentage
	KeyServiceAgentRegistrationFeeAmount    = types.KeyServiceAgentRegistrationFeeAmount
	KeyEvaluationAgentRegistrationFeeAmount = types.KeyEvaluationAgentRegistrationFeeAmount
	KeyEvaluationPayFeePercentage           = types.KeyEvaluationPayFeePercentage
	KeyEvaluationPayNodeFeePercentage       = types.KeyEvaluationPayNodeFeePercentage

)
