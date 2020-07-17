package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"net/http"
	"strings"

	"github.com/tokenchain/ixo-blockchain/x/payments/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/payments/createPaymentTemplate", createPaymentTemplateHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/payments/createPaymentContract", createPaymentContractHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/payments/createSubscription", createSubscriptionHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/payments/setPaymentContractAuthorisation", setPaymentContractAuthorisationHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/payments/grantDiscount", grantDiscountHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/payments/revokeDiscount", revokeDiscountHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/payments/effectPayment", effectPaymentHandler(cliCtx)).Methods("POST")
}

const (
	TRUE  = "true"
	FALSE = "false"
)

func parseBool(boolStr, boolName string) (bool, error) {
	boolStr = strings.ToLower(strings.TrimSpace(boolStr))
	if boolStr == TRUE {
		return true, nil
	} else if boolStr == FALSE {
		return false, nil
	} else {
		return false, exported.ErrInvalidArgument(fmt.Sprintf("%s is not a valid bool (true/false)", boolName))
	}
}
func writeHeadf(w http.ResponseWriter, code int, format string, i ...interface{}) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(fmt.Sprintf(format, i...)))
}
func writeHead(w http.ResponseWriter, code int, txt string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(txt))
}
func createPaymentTemplateHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		templateJsonParam := r.URL.Query().Get("paymentTemplateJson")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		var template types.PaymentTemplate
		err := ctx.Codec.UnmarshalJSON([]byte(templateJsonParam), &template)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreatePaymentTemplate(template, sovrinDid.Did)

		//output, err := auth.SignAndBroadcastTxRest(ctx, msg, sovrinDid)
		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}

func createPaymentContractHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		contractIdParam := r.URL.Query().Get("paymentContractId")
		templateIdParam := r.URL.Query().Get("paymentTemplateId")
		payerAddrParam := r.URL.Query().Get("payerAddr")
		canDeauthoriseParam := r.URL.Query().Get("canDeauthorise")
		discountIdParam := r.URL.Query().Get("discountId")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		payerAddr, err := sdk.AccAddressFromBech32(payerAddrParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		canDeauthorise, err := parseBool(canDeauthoriseParam, "canDeauthorise")
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		discountId, err := sdk.ParseUint(discountIdParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreatePaymentContract(templateIdParam,
			contractIdParam, payerAddr, canDeauthorise, discountId, sovrinDid.Did)

		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}

func createSubscriptionHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		subIdParam := r.URL.Query().Get("subId")
		contractIdParam := r.URL.Query().Get("paymentContractId")
		maxPeriodsParam := r.URL.Query().Get("maxPeriods")
		periodParam := r.URL.Query().Get("period")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		maxPeriods, err := sdk.ParseUint(maxPeriodsParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		var period types.Period
		err = ctx.Codec.UnmarshalJSON([]byte(periodParam), &period)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateSubscription(subIdParam, contractIdParam,
			maxPeriods, period, sovrinDid.Did)

		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}

func setPaymentContractAuthorisationHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		contractIdParam := r.URL.Query().Get("paymentContractId")
		authorisedParam := r.URL.Query().Get("authorised")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		authorised, err := parseBool(authorisedParam, "authorised")
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSetPaymentContractAuthorisation(contractIdParam,
			authorised, sovrinDid.Did)

		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}

func grantDiscountHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		contractIdParam := r.URL.Query().Get("paymentContractId")
		discountIdParam := r.URL.Query().Get("discountId")
		recipientAddrParam := r.URL.Query().Get("recipientAddr")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		discountId, err := sdk.ParseUint(discountIdParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		recipientAddr, err := sdk.AccAddressFromBech32(recipientAddrParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgGrantDiscount(contractIdParam, discountId,
			recipientAddr, sovrinDid.Did)

		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}

func revokeDiscountHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		contractIdParam := r.URL.Query().Get("paymentContractId")
		holderAddrParam := r.URL.Query().Get("holderAddr")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		holderAddr, err := sdk.AccAddressFromBech32(holderAddrParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRevokeDiscount(contractIdParam, holderAddr, sovrinDid.Did)

		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}

func effectPaymentHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		contractIdParam := r.URL.Query().Get("paymentContractId")
		sovrinDidParam := r.URL.Query().Get("sovrinDid")

		mode := r.URL.Query().Get("mode")
		ctx = ctx.WithBroadcastMode(mode)

		sovrinDid, err := exported.UnmarshalDxpDid(sovrinDidParam)
		if err != nil {
			writeHead(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgEffectPayment(contractIdParam, sovrinDid.Did)

		output, err := did.NewDidTxBuild(ctx, msg, sovrinDid).SignAndBroadcastTxRest()
		if err != nil {
			writeHead(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, output)
	}
}
