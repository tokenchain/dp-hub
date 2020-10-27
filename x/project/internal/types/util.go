package types

import (
	"github.com/tokenchain/dp-hub/x/did/exported"

	"strings"
)

func NewMsgCreateProject(senderDid exported.Did, projectDoc ProjectDoc, projectDid exported.IxoDid) MsgCreateProject {
	return MsgCreateProject{
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		PubKey:     projectDid.GetPubKey(),
		Data:       projectDoc,
	}
}

func NewMsgUpdateProjectStatus(senderDid exported.Did, updateProjectStatusDoc UpdateProjectStatusDoc, projectDid exported.IxoDid) MsgUpdateProjectStatus {
	return MsgUpdateProjectStatus{
		TxHash:     "",
		SenderDid:  senderDid,
		ProjectDid: projectDid.Did,
		Data:       updateProjectStatusDoc,
	}
}

func NewMsgCreateAgent(txHash string, senderDid exported.Did, createAgentDoc CreateAgentDoc, projectDid exported.IxoDid) MsgCreateAgent {
	return MsgCreateAgent{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createAgentDoc,
	}
}

func NewMsgUpdateAgent(txHash string, senderDid exported.Did, updateAgentDoc UpdateAgentDoc, projectDid exported.IxoDid) MsgUpdateAgent {
	return MsgUpdateAgent{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       updateAgentDoc,
	}
}

func NewMsgCreateClaim(txHash string, senderDid exported.Did, createClaimDoc CreateClaimDoc, projectDid exported.IxoDid) MsgCreateClaim {
	return MsgCreateClaim{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createClaimDoc,
	}
}

func NewMsgCreateEvaluation(txHash string, senderDid exported.Did, createEvaluationDoc CreateEvaluationDoc, projectDid exported.IxoDid) MsgCreateEvaluation {
	return MsgCreateEvaluation{
		ProjectDid: projectDid.Did,
		TxHash:     txHash,
		SenderDid:  senderDid,
		Data:       createEvaluationDoc,
	}
}

func NewMsgWithdrawFunds(senderDid exported.Did, data WithdrawFundsDoc) MsgWithdrawFunds {
	return MsgWithdrawFunds{
		SenderDid: senderDid,
		Data:      data,
	}
}

func CheckNotEmpty(value string, name string) (valid bool, err error) {
	if strings.TrimSpace(value) == "" {
		return false, exported.UnknownRequest(name + " is empty.")
	} else {
		return true, nil
	}
}
