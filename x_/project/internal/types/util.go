package types

import (
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"strings"
)

func NewMsgCreateProject(senderDid types.Did, projectDoc ProjectDoc, projectDid types.SovrinDid) MsgCreateProject {
	return MsgCreateProject{
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		PubKey:     projectDid.VerifyKey,
		Data:       projectDoc,
	}
}

func NewMsgUpdateProjectStatus(senderDid types.Did, updateProjectStatusDoc UpdateProjectStatusDoc, projectDid types.SovrinDid) MsgUpdateProjectStatus {
	return MsgUpdateProjectStatus{
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		Data:       updateProjectStatusDoc,
	}
}

func NewMsgCreateAgent(txHash string, senderDid types.Did, createAgentDoc CreateAgentDoc, projectDid types.SovrinDid) MsgCreateAgent {
	return MsgCreateAgent{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createAgentDoc,
	}
}

func NewMsgUpdateAgent(txHash string, senderDid types.Did, updateAgentDoc UpdateAgentDoc, projectDid types.SovrinDid) MsgUpdateAgent {
	return MsgUpdateAgent{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       updateAgentDoc,
	}
}

func NewMsgCreateClaim(txHash string, senderDid types.Did, createClaimDoc CreateClaimDoc, projectDid types.SovrinDid) MsgCreateClaim {
	return MsgCreateClaim{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createClaimDoc,
	}
}

func NewMsgCreateEvaluation(txHash string, senderDid types.Did, createEvaluationDoc CreateEvaluationDoc, projectDid types.SovrinDid) MsgCreateEvaluation {
	return MsgCreateEvaluation{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createEvaluationDoc,
	}
}

func NewMsgWithdrawFunds(senderDid types.Did, data WithdrawFundsDoc) MsgWithdrawFunds {
	return MsgWithdrawFunds{
		SenderDid: senderDid,
		Data:      data,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, x.UnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
