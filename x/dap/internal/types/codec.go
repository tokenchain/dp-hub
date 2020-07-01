package types


import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	//RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}


// DidKeeper defines the did contract that must be fulfilled throughout the ixo module
type DidKeeper interface {
	GetDidDoc(ctx sdk.Context, did exported.Did) (exported.DidDoc, sdk.Error)
	SetDidDoc(ctx sdk.Context, did exported.DidDoc) (err sdk.Error)
	AddDidDoc(ctx sdk.Context, did exported.DidDoc)
	AddCredentials(ctx sdk.Context, did exported.Did, credential exported.DidCredential) (err sdk.Error)
	GetAllDidDocs(ctx sdk.Context) (didDocs []exported.DidDoc)
	GetAddDids(ctx sdk.Context) (dids []exported.Did)
}
