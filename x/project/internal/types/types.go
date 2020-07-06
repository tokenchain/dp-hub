package types

import (
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	InternalAccountID          string
	AccountMap                 map[InternalAccountID]sdk.AccAddress
	ProjectStatus              string
	ProjectStatusTransitionMap map[ProjectStatus][]ProjectStatus
)

func (id InternalAccountID) ToAddressKey(projectDid types.Did) string {
	return projectDid + "/" + string(id)
}

type StoredProjectDoc interface {
	GetEvaluatorPay() int64
	GetProjectDid() types.Did
	GetSenderDid() types.Did
	GetPubKey() string
	GetStatus() ProjectStatus
	SetStatus(status ProjectStatus)
}

const (
	NullStatus     ProjectStatus = ""
	CreatedProject ProjectStatus = "CREATED"
	PendingStatus  ProjectStatus = "PENDING"
	FundedStatus   ProjectStatus = "FUNDED"
	StartedStatus  ProjectStatus = "STARTED"
	StoppedStatus  ProjectStatus = "STOPPED"
	PaidoutStatus  ProjectStatus = "PAIDOUT"
)

var StateTransitions = initStateTransitions()

func initStateTransitions() ProjectStatusTransitionMap {
	return ProjectStatusTransitionMap{
		NullStatus:     {CreatedProject},
		CreatedProject: {PendingStatus},
		PendingStatus:  {CreatedProject, FundedStatus},
		FundedStatus:   {StartedStatus},
		StartedStatus:  {StoppedStatus},
		StoppedStatus:  {PaidoutStatus},
	}
}

func (next ProjectStatus) IsValidProgressionFrom(prev ProjectStatus) bool {
	validStatuses := StateTransitions[prev]
	for _, v := range validStatuses {
		if v == next {
			return true
		}
	}
	return false
}

type WithdrawalInfo struct {
	ActionID     string    `json:"actionID" yaml:"actionID"`
	ProjectDid   types.Did `json:"projectDid" yaml:"projectDid"`
	RecipientDid types.Did `json:"recipientDid" yaml:"recipientDid"`
	Amount       sdk.Coin  `json:"amount" yaml:"amount"`
}

type UpdateProjectStatusDoc struct {
	Status          ProjectStatus `json:"status" yaml:"status"`
	EthFundingTxnID string        `json:"ethFundingTxnID" yaml:"ethFundingTxnID"`
}

type ProjectDoc struct {
	NodeDid              string        `json:"nodeDid" yaml:"nodeDid"`
	RequiredClaims       string        `json:"requiredClaims" yaml:"requiredClaims"`
	EvaluatorPayPerClaim string        `json:"evaluatorPayPerClaim" yaml:"evaluatorPayPerClaim"`
	ServiceEndpoint      string        `json:"serviceEndpoint" yaml:"serviceEndpoint"`
	CreatedOn            string        `json:"createdOn" yaml:"createdOn"`
	CreatedBy            string        `json:"createdBy" yaml:"createdBy"`
	Status               ProjectStatus `json:"status" yaml:"status"`
}

func (pd ProjectDoc) GetEvaluatorPay() int64 {
	if pd.EvaluatorPayPerClaim == "" {
		return 0
	} else {
		i, err := strconv.ParseInt(pd.EvaluatorPayPerClaim, 10, 64)
		if err != nil {
			panic(err)
		}

		return i
	}
}

type CreateAgentDoc struct {
	AgentDid types.Did `json:"did" yaml:"did"`
	Role     string    `json:"role" yaml:"role"`
}

type AgentStatus = string

const (
	PendingAgent  AgentStatus = "0"
	ApprovedAgent AgentStatus = "1"
	RevokedAgent  AgentStatus = "2"
)

type UpdateAgentDoc struct {
	Did    types.Did   `json:"did" yaml:"did"`
	Status AgentStatus `json:"status" yaml:"status"`
	Role   string      `json:"role" yaml:"role"`
}

type CreateClaimDoc struct {
	ClaimID string `json:"claimID" yaml:"claimID"`
}

type ClaimStatus string

const (
	PendingClaim  ClaimStatus = "0"
	ApprovedClaim ClaimStatus = "1"
	RejectedClaim ClaimStatus = "2"
)

type CreateEvaluationDoc struct {
	ClaimID string      `json:"claimID" yaml:"claimID"`
	Status  ClaimStatus `json:"status" yaml:"status"`
}

type WithdrawFundsDoc struct {
	ProjectDid   types.Did `json:"projectDid" yaml:"projectDid"`
	RecipientDid types.Did `json:"recipientDid" yaml:"recipientDid"`
	Amount       sdk.Int   `json:"amount" yaml:"amount"`
	IsRefund     bool      `json:"isRefund" yaml:"isRefund"`
}
