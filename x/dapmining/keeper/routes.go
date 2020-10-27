package keeper

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/tokenchain/dp-hub/x/dapmining/keeper/session"
	"io/ioutil"
	"net/http"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	sub := r.PathPrefix("/auth").Subrouter()
	sub.HandleFunc("/login", loginHandler()).Methods("POST")
	sub.Handle("/logout", DefaultAuthMW(logoutHandler())).Methods("POST")
	sub.HandleFunc("/csrf_token", csrfTokenHandler()).Methods("GET")
	sub.Handle("/me", DefaultAuthMW(meHandler(ctx, ctx.Codec))).Methods("GET")
}

func loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var req LoginRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accountName := viper.GetString(AccountNameFlag)
		if req.Username != accountName {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		kbID, hotPW, err := authorize(req.Password)
		if err != nil {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		err = session.SetStrings(w, r, session.KeybaseIDKey, kbID, session.KeybasePassphraseKey, hotPW)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store, err := session.SessionStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delete(store.Values, session.KeybaseIDKey)
		err = store.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func csrfTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := GetCSRFToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(tok))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


func meHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		/*	addr := session.MustGetKBFromSession(r).
			res := &MeResponse{Address: addr}
			resB := cdc.MustMarshalJSON(res)
			rest.PostProcessResponseBare(w, ctx, resB)*/
	}
}

func authorize(passphrase string) (string, string, error) {
	kb := keys.NewInMemoryKeyBase()
	accountName := viper.GetString(AccountNameFlag)
	pk, err := kb.ExportPrivateKeyObject(accountName, passphrase)
	if err != nil {
		return "", "", err
	}
	hotPassphrase := session.ReadStr32()
	return session.ReplaceKB(accountName, hotPassphrase, pk), hotPassphrase, nil
}
