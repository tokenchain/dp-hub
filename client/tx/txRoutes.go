package tx

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	auRest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
)

func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", auRest.EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/decode", auRest.DecodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/sign_data", SignDataRequest(cliCtx)).Methods("POST")
}
