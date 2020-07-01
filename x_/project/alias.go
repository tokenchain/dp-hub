package project

import (
	"github.com/tokenchain/ixo-blockchain/x/project/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/project/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey

	DefaultCodespace = types.DefaultCodespace
	PaidoutStatus    = types.PaidoutStatus
	FundedStatus     = types.FundedStatus

	TypeMsgCreateProject = types.TypeMsgCreateProject

	MsgCreateProjectFee            = types.MsgCreateProjectFee
	MsgCreateProjectTransactionFee = types.MsgCreateProjectTransactionFee
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState

	MsgCreateProject       = types.MsgCreateProject
	MsgUpdateProjectStatus = types.MsgUpdateProjectStatus
	MsgCreateAgent         = types.MsgCreateAgent
	MsgUpdateAgent         = types.MsgUpdateAgent
	MsgCreateClaim         = types.MsgCreateClaim
	MsgCreateEvaluation    = types.MsgCreateEvaluation
	MsgWithdrawFunds       = types.MsgWithdrawFunds

	StoredProjectDoc  = types.StoredProjectDoc
	WithdrawalInfo    = types.WithdrawalInfo
	AccountMap        = types.AccountMap
	InternalAccountID = types.InternalAccountID
)

var (
	// function aliases
	NewKeeper      = keeper.NewKeeper
	NewQuerier     = keeper.NewQuerier
	ParamKeyTable  = types.ParamKeyTable
	NewParams      = types.NewParams
	DefaultParams  = types.DefaultParams
	ValidateParams = types.ValidateParams
	RegisterCodec  = types.RegisterCodec

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	// variable aliases
	ModuleCdc = types.ModuleCdc
)
