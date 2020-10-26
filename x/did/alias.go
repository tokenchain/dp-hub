package did

import (
	"github.com/tokenchain/dp-block/x/did/ed25519"
	"github.com/tokenchain/dp-block/x/did/exported"
	"github.com/tokenchain/dp-block/x/did/internal/keeper"
	"github.com/tokenchain/dp-block/x/did/internal/types"
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
/*	Did          = exported.Did
	DidDoc       = exported.DidDoc
	IxoDid       = exported.IxoDid*/

	MsgAddDid        = types.MsgAddDid
	MsgAddCredential = types.MsgAddCredential

/*	IxoTx        = ante.IxoTx
	IxoSignature = ante.IxoSignature
	IxoMsg       = ante.IxoMsg
	PubKeyGetter = ante.PubKeyGetter*/
)

var (
/*	NewDefaultPubKeyGetter = ante.NewDefaultPubKeyGetter
	DefaultAnteHandler     = ante.DefaultAnteHandler
	DidAnteHandler         = ante.DidAnteHandler
	NewDidTxBuild          = ante.NewDidTxBuild
	NewSignature           = ante.NewSignature
	DefaultTxDecoder       = ante.DefaultTxDecoder
	NewIxoTxSingleMsg      = ante.NewIxoTxSingleMsg
	DidToAddr              = ante.DidToAddr*/

	// function aliases
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	RegisterAmino = ed25519.RegisterAmino

	// Tx

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	UnmarshalIxoDid     = types.UnmarshalIxoDid

	// variable aliases
	ModuleCdc = types.ModuleCdc

	ErrorInvalidDid = exported.ErrorInvalidDidE
)
