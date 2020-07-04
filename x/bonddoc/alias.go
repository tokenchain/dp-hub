package bonddoc

import (
	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/types"
)

const (
	ModuleName              = types.ModuleName
	QuerierRoute            = types.QuerierRoute
	RouterKey               = types.RouterKey
	StoreKey                = types.StoreKey
	TypeMsgCreateBond       = types.TypeMsgCreateBond
	TypeMsgUpdateBondStatus = types.TypeMsgUpdateBondStatus
	DefaultCodespace        = types.ModuleName
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState

	MsgCreateBond       = types.MsgCreateBond
	MsgUpdateBondStatus = types.MsgUpdateBondStatus

	StoredBondDoc = types.StoredBondDoc
)

var (
	// function aliases
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	// variable aliases
	ModuleCdc = types.ModuleCdc
)
