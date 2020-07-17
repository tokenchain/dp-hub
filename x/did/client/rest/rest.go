package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/did", createDidRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/credential", addCredentialRequestHandler(cliCtx)).Methods("POST")
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// The .* is necessary so that a slash in the did gets included as part of the did
	r.HandleFunc("/didToAddr/{did:.*}", queryAddressFromDidRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/did/{did}", queryDidDocRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/did", queryAllDidsRequestHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/allDidDocs", queryAllDidDocsRequestHandler(cliCtx)).Methods("GET")
}
