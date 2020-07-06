package types

import (
	er "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
	"strings"
)

func NewMsgCreateBond(senderDid types.Did, bondDoc BondDoc, bondDid types.SovrinDid) MsgCreateBond {
	return MsgCreateBond{
		TxHash:    "",
		SenderDid: senderDid,
		BondDid:   bondDid.Did,
		PubKey:    bondDid.VerifyKey,
		Data:      bondDoc,
	}
}

func NewMsgUpdateBondStatus(senderDid types.Did, updateBondStatusDoc UpdateBondStatusDoc, bondDid types.SovrinDid) MsgUpdateBondStatus {
	return MsgUpdateBondStatus{
		SenderDid: senderDid,
		BondDid:   bondDid.Did,
		Data:      updateBondStatusDoc,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, er.Wrapf(er.ErrUnknownRequest, "%s is empty", name)
	} else {
		return true, nil
	}
}
