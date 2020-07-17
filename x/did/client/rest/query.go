package rest

import (
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	didtypes "github.com/tokenchain/ixo-blockchain/x/did/internal/types"

	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	rest "github.com/tokenchain/ixo-blockchain/client"
)

func queryAddressFromDidRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		if !exported.IsValidDid(vars["did"]) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("input is not a valid did"))
			return
		}

		accAddress := exported.DidToAddr(vars["did"])

		rest.PostProcessResponse(w, cliCtx.Codec, accAddress, true)
	}
}

func queryDidDocRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		didAddr := vars["did"]
		key := exported.Did(didAddr)
		res, _, err := cliCtx.QueryWithData(getQuery(key), nil)
		if err != nil {
			writeHeadf(w, http.StatusInternalServerError, "Could't query did. Error: %s", err.Error())
			return
		}
		if len(res) == 0 {
			writeHead(w, http.StatusNoContent, "No data for respected did address.")
			return
		}

		var didDoc didtypes.BaseDidDoc
		cliCtx.Codec.MustUnmarshalJSON(res, &didDoc)

		rest.PostProcessResponse(w, cliCtx.Codec, didDoc, true)
	}
}

func queryAllDidsRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		res, _, err := cliCtx.QueryWithData(getQuerySr(), nil)
		if err != nil {
			writeHeadf(w, http.StatusInternalServerError, "Could't query did. Error: %s", err.Error())
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var dids []exported.Did
		cliCtx.Codec.MustUnmarshalJSON(res, &dids)

		rest.PostProcessResponse(w, cliCtx.Codec, dids, true)
	}
}

func queryAllDidDocsRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		res, _, err := cliCtx.QueryWithData(getQueryAll(), nil)
		if err != nil {
			writeHeadf(w, http.StatusInternalServerError, "Could't query did. Error: %s", err.Error())
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			_, _ = w.Write([]byte("No data present."))
			return
		}

		var didDocs []didtypes.BaseDidDoc
		cliCtx.Codec.MustUnmarshalJSON(res, &didDocs)

		rest.PostProcessResponse(w, cliCtx.Codec, didDocs, true)
	}
}
