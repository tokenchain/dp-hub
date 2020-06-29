package rest

import (
	"encoding/json"
	"fmt"
	types2 "github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/project/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/project", createProjectRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/updateProjectStatus", updateProjectStatusRequestHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/createAgent", createAgentRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/createClaim", createClaimRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/createEvaluation", createEvaluationRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/withdrawFunds", withdrawFundsRequestHandler(cliCtx)).Methods("POST")
}

func createProjectRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		senderDid := r.URL.Query().Get("senderDid")
		projectDocParam := r.URL.Query().Get("projectDoc")
		didDocParam := r.URL.Query().Get("didDoc")
		mode := r.URL.Query().Get("mode")

		var projectDoc types.ProjectDoc
		err := json.Unmarshal([]byte(projectDocParam), &projectDoc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Could not unmarshall projectDoc into struct. Error: %s", err.Error())))
			return
		}

		didDoc, err := types2.UnmarshalSovrinDid(didDocParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)
		msg := types.NewMsgCreateProject(senderDid, projectDoc, didDoc)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, didDoc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func updateProjectStatusRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		senderDid := r.URL.Query().Get("senderDid")
		status := r.URL.Query().Get("status")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")
		mode := r.URL.Query().Get("mode")

		sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		projectStatus := types.ProjectStatus(status)
		if projectStatus != types.CreatedProject &&
			projectStatus != types.PendingStatus &&
			projectStatus != types.FundedStatus &&
			projectStatus != types.StartedStatus &&
			projectStatus != types.StoppedStatus &&
			projectStatus != types.PaidoutStatus {
			_, _ = w.Write([]byte("The status must be one of 'CREATED', " +
				"'PENDING', 'FUNDED', 'STARTED', 'STOPPED' or 'PAIDOUT'"))
			return
		}

		updateProjectStatusDoc := types.UpdateProjectStatusDoc{
			Status: projectStatus,
		}

		msg := types.NewMsgUpdateProjectStatus(senderDid, updateProjectStatusDoc, sovrinDid)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func createAgentRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		txHash := r.URL.Query().Get("txHash")
		senderDid := r.URL.Query().Get("senderDid")
		agentDid := r.URL.Query().Get("agentDid")
		role := r.URL.Query().Get("role")
		projectDidParam := r.URL.Query().Get("projectDid")
		mode := r.URL.Query().Get("mode")

		projectDid, err := types2.UnmarshalSovrinDid(projectDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		if role != "SA" && role != "EA" && role != "IA" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "The role must be one of 'SA', 'EA' or 'IA'")
			return
		}

		createAgentDoc := types.CreateAgentDoc{
			AgentDid: agentDid,
			Role:     role,
		}

		msg := types.NewMsgCreateAgent(txHash, senderDid, createAgentDoc, projectDid)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, projectDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func createClaimRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		txHash := r.URL.Query().Get("txHash")
		senderDid := r.URL.Query().Get("senderDid")
		claimId := r.URL.Query().Get("claimId")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")
		mode := r.URL.Query().Get("mode")

		sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		createClaimDoc := types.CreateClaimDoc{
			ClaimID: claimId,
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		msg := types.NewMsgCreateClaim(txHash, senderDid, createClaimDoc, sovrinDid)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func createEvaluationRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		txHash := r.URL.Query().Get("txHash")
		senderDid := r.URL.Query().Get("senderDid")
		claimDid := r.URL.Query().Get("claimDid")
		status := r.URL.Query().Get("status")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")
		mode := r.URL.Query().Get("mode")

		sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		claimStatus := types.ClaimStatus(status)
		if claimStatus != types.PendingClaim && claimStatus != types.ApprovedClaim && claimStatus != types.RejectedClaim {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "The status must be one of '0' (Pending), '1' (Approved) or '2' (Rejected)")
			return
		}

		createEvaluationDoc := types.CreateEvaluationDoc{
			ClaimID: claimDid,
			Status:  claimStatus,
		}

		msg := types.NewMsgCreateEvaluation(txHash, senderDid, createEvaluationDoc, sovrinDid)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func withdrawFundsRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		senderDidParam := r.URL.Query().Get("senderDid")
		dataParam := r.URL.Query().Get("data")
		mode := r.URL.Query().Get("mode")

		senderDid, err := types2.UnmarshalSovrinDid(senderDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var data types.WithdrawFundsDoc
		err = json.Unmarshal([]byte(dataParam), &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Could not unmarshall data into struct. Error: %s", err.Error())))
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		msg := types.NewMsgWithdrawFunds(senderDid.Did, data)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, senderDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
