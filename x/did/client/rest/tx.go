package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/tokenchain/ixo-blockchain/client/utils"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/did", createDidRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/credential", addCredentialRequestHandler(cliCtx)).Methods("POST")
}

func writeHeadf(w http.ResponseWriter, code int, format string, i ...interface{}) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(fmt.Sprintf(format, i...)))
}
func writeHead(w http.ResponseWriter, code int, txt string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(txt))
}
func createDidRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		didDocParam := r.URL.Query().Get("didDoc")
		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		sovrinDid, err := exported.UnmarshalDxpDid(didDocParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgAddDid(sovrinDid.Did, sovrinDid.GetPubKey())

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
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

		_, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s/%s", types.QuerierRoute, keeper.QueryDidDoc, did)
		if err != nil {
			writeHead(w, http.StatusBadRequest, "The did is not found")
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(didDocParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		t := time.Now()
		issued := t.Format(time.RFC3339)

		credTypes := []string{"Credential", "ProofOfKYC"}

		msg := types.NewMsgAddCredential(did, credTypes, sovrinDid.Did, issued)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
