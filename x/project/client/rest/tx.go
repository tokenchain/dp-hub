package rest

import (
	"encoding/json"
	"fmt"
	"github.com/tokenchain/dp-hub/x/did/exported"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/tokenchain/dp-hub/x/dap"
	"github.com/tokenchain/dp-hub/x/project/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/project", createProjectRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/updateProjectStatus", updateProjectStatusRequestHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/createAgent", createAgentRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/createClaim", createClaimRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/createEvaluation", createEvaluationRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/withdrawFunds", withdrawFundsRequestHandler(cliCtx)).Methods("POST")
}
func writeHeadf(w http.ResponseWriter, code int, format string, i ...interface{}) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(fmt.Sprintf(format, i...)))
}
func writeHead(w http.ResponseWriter, code int, txt string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(txt))
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
			writeHeadf(w, http.StatusBadRequest, "Could not unmarshall projectDoc into struct. Error: %s", err.Error())
			return
		}

		didDoc, err := exported.UnmarshalDxpDid(didDocParam)
		if err != nil {
			writeHeadf(w, http.StatusBadRequest, "Bad Request. Error: %s", err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)
		msg := types.NewMsgCreateProject(senderDid, projectDoc, didDoc)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, didDoc)
		if err != nil {
			writeHeadf(w, http.StatusInternalServerError, "Internal Server Error: %s", err.Error())
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

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
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

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
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

		projectDid, err := exported.UnmarshalDxpDid(projectDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
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

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, projectDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
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

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		createClaimDoc := types.CreateClaimDoc{
			ClaimID: claimId,
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		msg := types.NewMsgCreateClaim(txHash, senderDid, createClaimDoc, sovrinDid)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
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

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
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

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
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

		senderDid, err := exported.UnmarshalDxpDid(senderDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		var data types.WithdrawFundsDoc
		err = json.Unmarshal([]byte(dataParam), &data)
		if err != nil {
			writeHeadf(w, http.StatusBadRequest, "Could not unmarshall data into struct. Error: %s", err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)

		msg := types.NewMsgWithdrawFunds(senderDid.Did, data)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, senderDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
