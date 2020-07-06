package dapmining

import (
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/dapmining/keeper"
	"github.com/tokenchain/ixo-blockchain/x/dapmining/types"
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
	CodeType     = x.CodeType
	GenesisState = types.GenesisState
)
