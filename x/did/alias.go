package did

import (
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
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
	Did          = exported.Did
	DidDoc       = exported.DidDoc
	IxoDid       = exported.IxoDid

	MsgAddDid        = types.MsgAddDid
	MsgAddCredential = types.MsgAddCredential
)

var (

	// function aliases
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	RegisterAmino = ed25519.RegisterAmino

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	UnmarshalIxoDid     = types.UnmarshalIxoDid

	// variable aliases
	ModuleCdc = types.ModuleCdc

	ErrorInvalidDid = x.ErrorInvalidDidE
)
