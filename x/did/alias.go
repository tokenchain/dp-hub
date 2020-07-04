package did

import (
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

const (
	ModuleName       = types.ModuleName
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	StoreKey         = types.StoreKey
	DefaultCodespace = types.ModuleName
	QueryDidDoc      = keeper.QueryDidDoc
	QueryAllDids     = keeper.QueryAllDids
	QueryAllDidDocs  = keeper.QueryAllDidDocs
)

type (
	Keeper           = keeper.Keeper
	GenesisState     = types.GenesisState
	Did              = types.Did
	Claim            = types.Claim
	DidCredential    = types.DidCredential
	DidDoc           = types.DidDoc
	BaseDidDoc       = types.BaseDidDoc
	DxpDid           = types.DxpDid
	SovrinSecret     = types.SovrinSecret
	MsgAddDid        = types.MsgAddDid
	MsgAddCredential = types.MsgAddCredential
)

var (
	// function aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	// variable aliases
	ErrorInvalidDid = x.ErrorInvalidDidE
	ModuleCdc       = types.ModuleCdc
	DidKey          = types.DidKey
	IsValidDid      = types.IsValidDid
	ValidDid        = types.ValidDid
)
