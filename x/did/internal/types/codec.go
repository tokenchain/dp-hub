package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddDid{}, "did/MsgAddDid", nil)
	cdc.RegisterConcrete(MsgAddCredential{}, "did/MsgAddCredential", nil)
	// TODO: https://github.com/tokenchain/ixo-blockchain/issues/76
	cdc.RegisterConcrete(BaseDidDoc{}, "darkpool/BaseDidDoc", nil)
	//cdc.RegisterConcrete(ante.IxoTx{}, "darkpool/IxoTx", nil)
	//cdc.RegisterConcrete(DidCredential{}, "did/DidCredential", nil)

}

// ModuleCdc is the codec for the module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
