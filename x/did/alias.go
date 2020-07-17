package did

import (
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
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
	MsgAddDid        = types.MsgAddDid
	MsgAddCredential = types.MsgAddCredential
	IxoMsg           = types.IxoMsg
	IxoTx            = types.IxoTx
	IxoSignature     = types.IxoSignature
	BaseDidDoc       = types.BaseDidDoc
)

var (

	// function aliases
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	RegisterAmino = ed25519.RegisterAmino

	// Tx

	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	UnmarshalIxoDid     = types.UnmarshalIxoDid
	NewDidTxBuild       = types.NewDidTxBuild
	NewIxoTxSingleMsg   = types.NewIxoTxSingleMsg
	NewSignature        = types.NewSignature
	//NewMsgAddCredential = types.NewMsgAddCredential
	//NewMsgAddDid        = types.NewMsgAddDid
	CastTypeSdkTx        = types.CastTypeSdkTx
	DefaultTxDecoder        = types.DefaultTxDecoder

	// variable aliases
	ModuleCdc = types.ModuleCdc
)
