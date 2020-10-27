package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/tokenchain/dp-hub/client/utils"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/bonds", queryBondsHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}", RestBondDid),
		queryBondHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/batch", RestBondDid),
		queryBatchHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/last_batch", RestBondDid),
		queryLastBatchHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/current_price", RestBondDid),
		queryCurrentPriceHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/current_reserve", RestBondDid),
		queryCurrentReserveHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/price/{%s}", RestBondDid, RestBondAmount),
		queryCustomPriceHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/buy_price/{%s}", RestBondDid, RestBondAmount),
		queryBuyPriceHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/sell_return/{%s}", RestBondDid, RestBondAmount),
		querySellReturnHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/bonds/{%s}/swap_return/{%s}/{%s}", RestBondDid, RestFromTokenWithAmount, RestToToken),
		querySwapReturnHandler(cliCtx, queryRoute),
	).Methods("GET")
}

func queryBondsHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/bonds", queryRoute)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBondHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/bond/%s", queryRoute, bondDid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBatchHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/batch/%s", queryRoute, bondDid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryLastBatchHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/last_batch/%s", queryRoute, bondDid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCurrentPriceHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]

		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/current_price/%s", queryRoute, bondDid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCurrentReserveHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]

		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/current_reserve/%s", queryRoute, bondDid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCustomPriceHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		bondAmount := vars[RestBondAmount]

		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/custom_price/%s/%s", queryRoute, bondDid, bondAmount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBuyPriceHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		bondAmount := vars[RestBondAmount]

		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/buy_price/%s/%s", queryRoute, bondDid, bondAmount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySellReturnHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		bondAmount := vars[RestBondAmount]

		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/sell_return/%s/%s", queryRoute, bondDid, bondAmount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySwapReturnHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bondDid := vars[RestBondDid]
		fromTokenWithAmount := vars[RestFromTokenWithAmount]
		toToken := vars[RestToToken]

		reserveCoinWithAmount, err := sdk.ParseCoin(fromTokenWithAmount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		res, _, err := utils.QueryWithData(cliCtx, "custom/%s/swap_return/%s/%s/%s/%s", queryRoute, bondDid, reserveCoinWithAmount.Denom, reserveCoinWithAmount.Amount.String(), toToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
