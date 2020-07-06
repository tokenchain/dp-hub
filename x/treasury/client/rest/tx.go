package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	types2 "github.com/tokenchain/ixo-blockchain/x/dap/types"
	"net/http"

	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/treasury/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/treasury/send", sendRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/treasury/oracleTransfer", oracleTransferRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/treasury/oracleMint", oracleMintRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/treasury/oracleBurn", oracleBurnRequestHandler(cliCtx)).Methods("POST")
}

func sendRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		toDidParam := r.URL.Query().Get("toDid")
		amountParam := r.URL.Query().Get("amount")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		coins, err := sdk.ParseCoins(amountParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		msg := types.NewMsgSend(toDidParam, coins, sovrinDid)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, sovrinDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func oracleTransferRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		fromDidParam := r.URL.Query().Get("fromDid")
		toDidParam := r.URL.Query().Get("toDid")
		amountParam := r.URL.Query().Get("amount")
		oracleDidParam := r.URL.Query().Get("oracleDid")
		proofParam := r.URL.Query().Get("proof")

		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		coins, err := sdk.ParseCoins(amountParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		oracleDid, err := types2.UnmarshalSovrinDid(oracleDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		msg := types.NewMsgOracleTransfer(fromDidParam, toDidParam, coins, oracleDid, proofParam)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, oracleDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func oracleMintRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		toDidParam := r.URL.Query().Get("toDid")
		amountParam := r.URL.Query().Get("amount")
		oracleDidParam := r.URL.Query().Get("oracleDid")
		proofParam := r.URL.Query().Get("proof")

		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		coins, err := sdk.ParseCoins(amountParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		oracleDid, err := types2.UnmarshalSovrinDid(oracleDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		msg := types.NewMsgOracleMint(toDidParam, coins, oracleDid, proofParam)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, oracleDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

func oracleBurnRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		fromDidParam := r.URL.Query().Get("fromDid")
		amountParam := r.URL.Query().Get("amount")
		oracleDidParam := r.URL.Query().Get("oracleDid")
		proofParam := r.URL.Query().Get("proof")

		mode := r.URL.Query().Get("mode")
		cliCtx = cliCtx.WithBroadcastMode(mode)

		coins, err := sdk.ParseCoins(amountParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		oracleDid, err := types2.UnmarshalSovrinDid(oracleDidParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		msg := types.NewMsgOracleBurn(fromDidParam, coins, oracleDid, proofParam)

		output, err := dap.SignAndBroadcastTxRest(cliCtx, msg, oracleDid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}
