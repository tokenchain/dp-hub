package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"strings"
)

func NewMsgSend(toDid types.Did, amount sdk.Coins, senderDid types.SovrinDid) MsgSend {
	return MsgSend{
		PubKey:  senderDid.VerifyKey,
		FromDid: senderDid.Did,
		ToDid:   toDid,
		Amount:  amount,
	}
}

func NewMsgOracleTransfer(fromDid, toDid types.Did, amount sdk.Coins,
	oracleDid types.SovrinDid, proof string) MsgOracleTransfer {
	return MsgOracleTransfer{
		PubKey:    oracleDid.VerifyKey,
		OracleDid: oracleDid.Did,
		FromDid:   fromDid,
		ToDid:     toDid,
		Amount:    amount,
		Proof:     proof,
	}
}

func NewMsgOracleMint(toDid types.Did, amount sdk.Coins,
	oracleDid types.SovrinDid, proof string) MsgOracleMint {
	return MsgOracleMint{
		PubKey:    oracleDid.VerifyKey,
		OracleDid: oracleDid.Did,
		ToDid:     toDid,
		Amount:    amount,
		Proof:     proof,
	}
}

func NewMsgOracleBurn(fromDid types.Did, amount sdk.Coins,
	oracleDid types.SovrinDid, proof string) MsgOracleBurn {
	return MsgOracleBurn{
		PubKey:    oracleDid.VerifyKey,
		OracleDid: oracleDid.Did,
		FromDid:   fromDid,
		Amount:    amount,
		Proof:     proof,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, x.UnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
