package dapmining

import (
	"github.com/tokenchain/dp-hub/x/dapmining/keeper"
	"github.com/tokenchain/dp-hub/x/dapmining/types"
	"github.com/tokenchain/dp-hub/x/did/exported"
)

//noinspection GoNameStartsWithPackageName
const (
	DefaultCodespace = types.ModuleName

	BondsMintBurnAccount = types.DapDistributionsReward

	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

//noinspection GoNameStartsWithPackageName
var (
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	RegisterCodec       = types.RegisterCodec
	ModuleCdc           = types.ModuleCdc
	RegisterInvariants  = keeper.RegisterInvariants
	NewQuerier          = keeper.NewQuerier
)

type (
	Keeper       = keeper.Keeper
	CodeType     = exported.CodeType
	GenesisState = types.GenesisState
)
