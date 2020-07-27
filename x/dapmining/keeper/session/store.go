package session

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gorilla/sessions"
	"github.com/tendermint/dex-demo/embedded/session"

	"github.com/tendermint/tendermint/crypto"
	"net/http"
	"sync"
)

const sessionName = "uex_session"

var SessionStore = sessions.NewCookieStore(generateSessionKey())

// TODO: pull from config
func generateSessionKey() []byte {
	var out [32]byte
	return out[:]
}

var kb *keys.Info
var currID string
var mtx sync.RWMutex

func GetKBFromSession(r *http.Request) (*keys.Info, error) {
	id, err := GetStr(r, KeybaseIDKey)
	if err != nil {
		return nil, err
	}
	kb := GetKB(id)
	if kb == nil {
		return nil, errors.Wrap(errors.New("accountmine", 1, "no keybase found"), "no keybase found")
	}
	return kb, nil
}

func MustGetKBFromSession(r *http.Request) *keys.Info {
	kb, err := GetKBFromSession(r)
	if err != nil {
		panic(err)
	}
	return kb
}

func MustGetKBPassphraseFromSession(r *http.Request) string {
	return session.MustGetStr(r, KeybasePassphraseKey)
}

func GetKB(id string) *keys.Info {
	mtx.RLock()
	defer mtx.RUnlock()
	if currID != id {
		return nil
	}

	return kb
}
/*
func NewHotKeybase(name string, passphrase string, pk crypto.PrivKey) *Keybase {
	armor := mintkey.EncryptArmorPrivKey(pk, passphrase, name)
	addr := sdk.AccAddress(pk.PubKey().Address())

	return accounts.Account{}
}*/

func ReplaceKB(name string, passphrase string, pk crypto.PrivKey) string {
	mtx.Lock()
	defer mtx.Unlock()
	currID = ReadStr32()
	//kb = NewHotKeybase(name, passphrase, pk)

	return currID
}
