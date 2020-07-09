package ed25519

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
)

// ModuleCdc is the codec for the module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterAmino(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterAmino(cdc *amino.Codec) {



}
