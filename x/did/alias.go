package did

import (
	"github.com/tokenchain/dp-hub/x"
	"github.com/tokenchain/dp-hub/x/did/internal/keeper"
	"github.com/tokenchain/dp-hub/x/did/internal/types"
)

const (
	ModuleName       = types.ModuleName
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	StoreKey         = types.StoreKey
	DefaultCodespace = types.ModuleName
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
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

	ErrorInvalidDid = x.ErrorInvalidDidE
)
