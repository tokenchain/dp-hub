package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/tokenchain/ixo-blockchain/x/did/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/ixo/sovrin"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/did", createDidRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/credential", addCredentialRequestHandler(cliCtx)).Methods("POST")
}

func createDidRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		didDocParam := r.URL.Query().Get("didDoc")
		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		sovrinDid, err := sovrin.UnmarshalSovrinDid(didDocParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		msg := types.NewMsgAddDid(sovrinDid.Did, sovrinDid.VerifyKey)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func addCredentialRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		did := r.URL.Query().Get("did")
		didDocParam := r.URL.Query().Get("signerDidDoc")
		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		_, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute,
			keeper.QueryDidDoc, did), nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("The did is not found"))
			return
		}

		sovrinDid, err := sovrin.UnmarshalSovrinDid(didDocParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		t := time.Now()
		issued := t.Format(time.RFC3339)

		credTypes := []string{"Credential", "ProofOfKYC"}

		msg := types.NewMsgAddCredential(did, credTypes, sovrinDid.Did, issued)

		output, err := ixo.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
