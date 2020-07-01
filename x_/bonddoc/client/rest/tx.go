package rest

import (
	"encoding/json"
	"fmt"
	types2 "github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bond", createBondRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/updateBondStatus", updateBondStatusRequestHandler(cliCtx)).Methods("PUT")
}

func createBondRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		senderDid := r.URL.Query().Get("senderDid")
		bondDocParam := r.URL.Query().Get("bondDoc")
		didDocParam := r.URL.Query().Get("didDoc")
		mode := r.URL.Query().Get("mode")

		var bondDoc types.BondDoc
		err := json.Unmarshal([]byte(bondDocParam), &bondDoc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Could not unmarshall bondDoc into struct. Error: %s", err.Error())))
			return
		}

		didDoc, err := types2.UnmarshalSovrinDid(didDocParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(mode)
		msg := types.NewMsgCreateBond(senderDid, bondDoc, didDoc)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, didDoc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func updateBondStatusRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

		bondStatus := types.BondStatus(status)
		if bondStatus != types.PreIssuanceStatus &&
			bondStatus != types.OpenStatus &&
			bondStatus != types.SuspendedStatus &&
			bondStatus != types.ClosedStatus &&
			bondStatus != types.SettlementStatus &&
			bondStatus != types.EndedStatus {
			_, _ = w.Write([]byte("The status must be one of 'PREISSUANCE', " +
				"'OPEN', 'SUSPENDED', 'CLOSED', 'SETTLEMENT' or 'ENDED'"))
			return
		}

		updateBondStatusDoc := types.UpdateBondStatusDoc{
			Status: bondStatus,
		}

		msg := types.NewMsgUpdateBondStatus(senderDid, updateBondStatusDoc, sovrinDid)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
