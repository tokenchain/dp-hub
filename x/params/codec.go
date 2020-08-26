package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tokenchain/ixo-blockchain/x/params/types"
)

// ModuleCdc is the codec of module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(types.ParameterChangeProposal{}, "darkpool/params/ParameterChangeProposal", nil)
}
