package bonds

import (
	"github.com/tokenchain/ixo-blockchain/x/bonds/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/bonds/internal/types"
)

//noinspection GoNameStartsWithPackageName
const (
	DefaultCodespace = types.DefaultCodespace

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
	Keeper       = keeper.Keeper
	CodeType     = types.CodeType
	GenesisState = types.GenesisState
)
