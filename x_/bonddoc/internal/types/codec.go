package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateBond{}, "bonddoc/CreateBond", nil)
	cdc.RegisterConcrete(MsgUpdateBondStatus{}, "bonddoc/UpdateBondStatus", nil)

	// TODO: https://github.com/tokenchain/ixo-blockchain/issues/76
	// cdc.RegisterInterface((*StoredBondDoc)(nil), nil)
}

// ModuleCdc is the codec for the module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
