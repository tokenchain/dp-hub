package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"strings"
)

func NewMsgSend(toDidOrAddr string, amount sdk.Coins, senderDid did.Did) MsgSend {
	return MsgSend{
		FromDid:     senderDid,
		ToDidOrAddr: toDidOrAddr,
		Amount:      amount,
	}
}

func NewMsgOracleTransfer(fromDid did.Did, toDidOrAddr string, amount sdk.Coins,
	oracleDid did.Did, proof string) MsgOracleTransfer {
	return MsgOracleTransfer{
		OracleDid:   oracleDid,
		FromDid:     fromDid,
		ToDidOrAddr: toDidOrAddr,
		Amount:      amount,
		Proof:       proof,
	}
}

func NewMsgOracleMint(toDidOrAddr string, amount sdk.Coins,
	oracleDid did.Did, proof string) MsgOracleMint {
	return MsgOracleMint{
		OracleDid:   oracleDid,
		ToDidOrAddr: toDidOrAddr,
		Amount:      amount,
		Proof:       proof,
	}
}

func NewMsgOracleBurn(fromDid did.Did, amount sdk.Coins,
	oracleDid did.Did, proof string) MsgOracleBurn {
	return MsgOracleBurn{
		OracleDid: oracleDid,
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
