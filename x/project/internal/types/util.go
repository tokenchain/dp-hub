package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"

	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/ixo/sovrin"
)

func NewMsgCreateProject(senderDid ixo.Did, projectDoc ProjectDoc, projectDid sovrin.SovrinDid) MsgCreateProject {
	return MsgCreateProject{
		SignBytes:  "",
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		PubKey:     projectDid.VerifyKey,
		Data:       projectDoc,
	}
}

func NewMsgUpdateProjectStatus(senderDid ixo.Did, updateProjectStatusDoc UpdateProjectStatusDoc, projectDid sovrin.SovrinDid) MsgUpdateProjectStatus {
	return MsgUpdateProjectStatus{
		SignBytes:  "",
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		Data:       updateProjectStatusDoc,
	}
}

func NewMsgCreateAgent(txHash string, senderDid ixo.Did, createAgentDoc CreateAgentDoc, projectDid sovrin.SovrinDid) MsgCreateAgent {
	return MsgCreateAgent{
		SignBytes:  "",
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createAgentDoc,
	}
}

func NewMsgUpdateAgent(txHash string, senderDid ixo.Did, updateAgentDoc UpdateAgentDoc, projectDid sovrin.SovrinDid) MsgUpdateAgent {
	return MsgUpdateAgent{
		SignBytes:  "",
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       updateAgentDoc,
	}
}

func NewMsgCreateClaim(txHash string, senderDid ixo.Did, createClaimDoc CreateClaimDoc, projectDid sovrin.SovrinDid) MsgCreateClaim {
	return MsgCreateClaim{
		SignBytes:  "",
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createClaimDoc,
	}
}

func NewMsgCreateEvaluation(txHash string, senderDid ixo.Did, createEvaluationDoc CreateEvaluationDoc, projectDid sovrin.SovrinDid) MsgCreateEvaluation {
	return MsgCreateEvaluation{
		SignBytes:  "",
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createEvaluationDoc,
	}
}

func NewMsgWithdrawFunds(senderDid ixo.Did, data WithdrawFundsDoc) MsgWithdrawFunds {
	return MsgWithdrawFunds{
		SignBytes: "",
		SenderDid: senderDid,
		Data:      data,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err sdk.Error) {
	if strings.TrimSpace(value) == "" {
		return false, sdk.ErrUnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
