package cligen

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"io"
)

const (
	flagDryRun      = "dry-run"
	flagUserEntropy = "unsafe-entropy"
	flagInteractive = "interactive"
	flagRecover     = "recover"
	flagNoBackup    = "no-backup"
	flagAccount     = "account"
	flagIndex       = "index"
	flagMultisig    = "multisig"
	flagNoSort      = "nosort"
	flagHDPath      = "hd-path"
	flagKeyAlgo     = "algo"

	mnemonicEntropySize = 256
	mnem                = "poem recycle raw armor pill emotion clutch test glad process phrase know butter clever idea allow castle exchange turn bundle intact clown cricket crane"
)

func getKeybase(transient bool, buf io.Reader) (keys.Keybase, error) {
	if transient {
		return keys.NewInMemory(), nil
	}
	return keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
}
