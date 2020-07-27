package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/tokenchain/ixo-blockchain/client/utils"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"net/http"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/tokenchain/ixo-blockchain/x/did/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// The .* is necessary so that a slash in the did gets included as part of the did
	r.HandleFunc("/didToAddr/{did:.*}", queryAddressFromDidRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/did/{did}", queryDidDocRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/did", queryAllDidsRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/allDidDocs", queryAllDidDocsRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/checkName/{name}", queryCheckNameSystem(cliCtx)).Methods("GET")
}

type (
	ConfirmBool struct {
		Bool      bool   `json:"bool" yaml:"bool"`
		TimeStamp string `json:"time" yaml:"time"`
	}
)

func queryCheckNameSystem(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		kb, err := utils.GetKeybase(false, cliCtx.Input)
		if err != nil {
			return
		}
		name := vars["name"]
		_, err = kb.Get(strings.ToLower(name))
		result := ConfirmBool{
			Bool:      true,
			TimeStamp: time.Now().String(),
		}

		if err == nil {
			result.Bool = false
		}

		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}
func queryAddressFromDidRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if !exported.IsValidDid(vars["did"]) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("input is not a valid did"))
			return
		}
		//accAddress := ante.DidToAddr(vars["did"])
		key := exported.Did(vars["did"])
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s/%s", types.QuerierRoute, keeper.QueryDidDoc, key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Could't query did. Error: %s", err.Error())))
			return
		}
		var didDoc types.BaseDidDoc
		cliCtx.Codec.MustUnmarshalJSON(res, &didDoc)
		address_dx0 := exported.VerifyKeyToAddrEd25519(didDoc.PubKey)
		rest.PostProcessResponseBare(w, cliCtx, address_dx0)
	}
}
func queryDidDocRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		didAddr := vars["did"]
		key := exported.Did(didAddr)
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s/%s", types.QuerierRoute, keeper.QueryDidDoc, key)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Could't query did. Error: %s", err.Error())))
			return
		}
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			_, _ = w.Write([]byte("No data for respected did address."))
			return
		}

		var didDoc types.BaseDidDoc
		cliCtx.Codec.MustUnmarshalJSON(res, &didDoc)
		rest.PostProcessResponseBare(w, cliCtx, didDoc)
		//rest.PostProcessResponse(w, cliCtx.Codec, didDoc, true)
	}
}

func queryAllDidsRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s", types.QuerierRoute,
			keeper.QueryAllDids)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Could't query did. Error: %s", err.Error())))
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var dids []exported.Did
		cliCtx.Codec.MustUnmarshalJSON(res, &dids)
		rest.PostProcessResponseBare(w, cliCtx, dids)
		//rest.PostProcessResponse(w, cliCtx.Codec, dids, true)
	}
}

func queryAllDidDocsRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s", types.QuerierRoute,
			keeper.QueryAllDidDocs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Could't query did. Error: %s", err.Error())))
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			_, _ = w.Write([]byte("No data present."))
			return
		}

		var didDocs []types.BaseDidDoc
		cliCtx.Codec.MustUnmarshalJSON(res, &didDocs)
		rest.PostProcessResponseBare(w, cliCtx, didDocs)
		//rest.PostProcessResponse(w, cliCtx.Codec, didDocs, true)
	}
}
