package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"strings"

	"github.com/tokenchain/ixo-blockchain/x/ixo/sovrin"
)

func NewMsgCreateBond(senderDid ixo.Did, bondDoc BondDoc, bondDid sovrin.SovrinDid) MsgCreateBond {
	return MsgCreateBond{
		SignBytes: "",
		TxHash:    "",
		SenderDid: senderDid,
		BondDid:   bondDid.Did,
		PubKey:    bondDid.VerifyKey,
		Data:      bondDoc,
	}
}

func NewMsgUpdateBondStatus(senderDid ixo.Did, updateBondStatusDoc UpdateBondStatusDoc, bondDid sovrin.SovrinDid) MsgUpdateBondStatus {
	return MsgUpdateBondStatus{
		SignBytes: "",
		SenderDid: senderDid,
		BondDid:   bondDid.Did,
		Data:      updateBondStatusDoc,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err sdk.Error) {
	if strings.TrimSpace(value) == "" {
		return false, sdk.ErrUnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
