package bonds

import (
	"github.com/tokenchain/dp-hub/x/bonds/internal/keeper"
	"github.com/tokenchain/dp-hub/x/bonds/internal/types"
	"github.com/tokenchain/dp-hub/x/did/exported"
)

//noinspection GoNameStartsWithPackageName
const (
	DefaultCodespace = types.ModuleName

	BondsMintBurnAccount       = types.BondsMintBurnAccount
	BatchesIntermediaryAccount = types.BatchesIntermediaryAccount

	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

//noinspection GoNameStartsWithPackageName
var (
	// function aliases
	RegisterInvariants = keeper.RegisterInvariants
	NewKeeper          = keeper.NewKeeper
	NewQuerier         = keeper.NewQuerier
	RegisterCodec      = types.RegisterCodec

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	// variable aliases
	ModuleCdc            = types.ModuleCdc
	BondsKeyPrefix       = types.BondsKeyPrefix
	BatchesKeyPrefix     = types.BatchesKeyPrefix
	LastBatchesKeyPrefix = types.LastBatchesKeyPrefix
)

type (
	Keeper        = keeper.Keeper
	Bond          = types.Bond
	CodeType      = exported.CodeType
	MsgCreateBond = types.MsgCreateBond
	MsgEditBond   = types.MsgEditBond
	MsgBuy        = types.MsgBuy
	MsgSell       = types.MsgSell
	MsgSwap       = types.MsgSwap
	GenesisState  = types.GenesisState
)
