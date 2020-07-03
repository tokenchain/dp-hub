package types

import (
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"strings"
)

func NewMsgCreateProject(senderDid did.Did, projectDoc ProjectDoc, projectDid did.DxpDid) MsgCreateProject {
	return MsgCreateProject{
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		PubKey:     projectDid.VerifyKey,
		Data:       projectDoc,
	}
}

func NewMsgUpdateProjectStatus(senderDid did.Did, updateProjectStatusDoc UpdateProjectStatusDoc, projectDid did.DxpDid) MsgUpdateProjectStatus {
	return MsgUpdateProjectStatus{
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		Data:       updateProjectStatusDoc,
	}
}

func NewMsgCreateAgent(txHash string, senderDid did.Did, createAgentDoc CreateAgentDoc, projectDid did.DxpDid) MsgCreateAgent {
	return MsgCreateAgent{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createAgentDoc,
	}
}

func NewMsgUpdateAgent(txHash string, senderDid did.Did, updateAgentDoc UpdateAgentDoc, projectDid did.DxpDid) MsgUpdateAgent {
	return MsgUpdateAgent{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       updateAgentDoc,
	}
}

func NewMsgCreateClaim(txHash string, senderDid did.Did, createClaimDoc CreateClaimDoc, projectDid did.DxpDid) MsgCreateClaim {
	return MsgCreateClaim{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createClaimDoc,
	}
}

func NewMsgCreateEvaluation(txHash string, senderDid did.Did, createEvaluationDoc CreateEvaluationDoc, projectDid did.DxpDid) MsgCreateEvaluation {
	return MsgCreateEvaluation{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createEvaluationDoc,
	}
}

func NewMsgWithdrawFunds(senderDid did.Did, data WithdrawFundsDoc) MsgWithdrawFunds {
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
