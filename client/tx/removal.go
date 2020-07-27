package tx

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	genutilrest "github.com/cosmos/cosmos-sdk/x/genutil/client/rest"
	"github.com/gorilla/mux"
	utils2 "github.com/tokenchain/ixo-blockchain/client/utils"
	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/did/ante"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/project"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)
type (
	SignDataReq struct {
		Msg    string `json:"msg" yaml:"msg"`
		PubKey string `json:"pub_key" yaml:"pub_key"`
	}

	SignDataResponse struct {
		SignBytes string      `json:"sign_bytes" yaml:"sign_bytes"`
		Fee       auth.StdFee `json:"fee" yaml:"fee"`
	}
	BroadcastReq struct {
		Tx   string `json:"tx" yaml:"tx"`
		Mode string `json:"mode" yaml:"mode"`
	}
)

func QueryTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		output, err := utils2.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if output.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
		//	rest.PostProcessResponse(w, cliCtx, output)
	}
}

func QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not parse query parameters %s", err.Error()))
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				genutilrest.QueryGenesisTxs(cliCtx, w)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, cliCtx, txs)
			return
		}

		events, page, limit, err = rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

func BroadcastTxRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// The only line in this function different from that in Cosmos SDK
		// is the one below. Instead of codec (JSON) marshalling, hex is used
		// so that the DefaultTxDecoder can successfully recognize the IxoTx
		//
		// txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)

		txBytes, err := hex.DecodeString(strings.TrimPrefix(req.Tx, "0x"))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(req.Mode)

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}

func SignDataRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignDataReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msgBytes, err := hex.DecodeString(strings.TrimPrefix(req.Msg, "0x"))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var msg sdk.Msg
		err = cliCtx.Codec.UnmarshalJSON(msgBytes, &msg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// all messages must be of type ixo.IxoMsg
		ixoMsg, ok := msg.(ante.IxoMsg)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, exported.IntErr("msg must be ixo.IxoMsg").Error())
			return
		}
		msgs := []sdk.Msg{ixoMsg}

		// obtain stdSignMsg (create-project is a special case)
		var stdSignMsg auth.StdSignMsg
		switch ixoMsg.Type() {
		case project.TypeMsgCreateProject:
			stdSignMsg = ixoMsg.(project.MsgCreateProject).ToStdSignMsg(
				project.MsgCreateProjectFee)
		default:
			// Deduce and set signer address
			signerAddress := exported.VerifyKeyToAddr(req.PubKey)
			cliCtx = cliCtx.WithFromAddress(signerAddress)

			txBldr, err := utils.PrepareTxBuilder(auth.NewTxBuilderFromCLI(cliCtx.Input), cliCtx)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			// Build the transaction
			stdSignMsg, err = txBldr.BuildSignMsg(msgs)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			// Create dummy tx with blank signature for fee approximation
			signature := auth.StdSignature{}
			tx := auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee,
				[]auth.StdSignature{signature}, stdSignMsg.Memo)

			// Approximate fee
			fee, err := dap.ApproximateFeeForTx(cliCtx, tx, txBldr.ChainID())
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			stdSignMsg.Fee = fee
		}

		// Produce response from sign bytes and fees
		output := SignDataResponse{
			SignBytes: string(stdSignMsg.Bytes()),
			Fee:       stdSignMsg.Fee,
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}
