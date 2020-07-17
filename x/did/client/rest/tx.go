package rest

import (
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	didtypes "github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

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

		msg := didtypes.NewMsgAddDid(sovrinDid.Did, sovrinDid.GetPubKey())

		output, err := didtypes.NewDidTxBuild(cliCtx, msg, sovrinDid).SignAndBroadcastTxRest()
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
		did_q := r.URL.Query().Get("did")
		didDocParam := r.URL.Query().Get("signerDidDoc")
		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		_, _, err := cliCtx.QueryWithData(getQuery(did_q), nil)
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

		msg := didtypes.NewMsgAddCredential(did_q, credTypes, sovrinDid.Did, issued)

		output, err := didtypes.NewDidTxBuild(cliCtx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
