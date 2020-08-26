package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgList{}, "darkpool/dex/MsgList", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "darkpool/dex/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "darkpool/dex/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "darkpool/dex/MsgTransferTradingPairOwnership", nil)
	cdc.RegisterConcrete(DelistProposal{}, "darkpool/dex/DelistProposal", nil)
	cdc.RegisterConcrete(MsgCreateOperator{}, "darkpool/dex/CreateOperator", nil)
	cdc.RegisterConcrete(MsgUpdateOperator{}, "darkpool/dex/UpdateOperator", nil)
}

// ModuleCdc represents generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
