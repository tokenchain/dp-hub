package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"strings"
)

func NewMsgSend(toDidOrAddr string, amount sdk.Coins, senderDid exported.Did) MsgSend {
	return MsgSend{
		FromDid:     senderDid,
		ToDidOrAddr: toDidOrAddr,
		Amount:      amount,
	}
}

func NewMsgOracleTransfer(fromDid exported.Did, toDidOrAddr string, amount sdk.Coins,
	oracleDid exported.Did, proof string) MsgOracleTransfer {
	return MsgOracleTransfer{
		OracleDid:   oracleDid,
		FromDid:     fromDid,
		ToDidOrAddr: toDidOrAddr,
		Amount:      amount,
		Proof:       proof,
	}
}

func NewMsgOracleMint(toDidOrAddr string, amount sdk.Coins,
	oracleDid exported.Did, proof string) MsgOracleMint {
	return MsgOracleMint{
		OracleDid:   oracleDid,
		ToDidOrAddr: toDidOrAddr,
		Amount:      amount,
		Proof:       proof,
	}
}

func NewMsgOracleBurn(fromDid exported.Did, amount sdk.Coins,
	oracleDid exported.Did, proof string) MsgOracleBurn {
	return MsgOracleBurn{
		OracleDid: oracleDid,
		FromDid:   fromDid,
		Amount:    amount,
		Proof:     proof,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, exported.UnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
