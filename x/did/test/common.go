package test

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/app"
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

	mnemonicEntropySize         = 256
	sample_mnom                 = "ignore sing before romance shiver hidden away despair soda gas moon merit borrow ten orbit sibling blame again pair estate siege dose horror rough"
	sample_did                  = "did:dxp:VrsU9cUAcYgF7f397xtjsX"
	sample_verifyKey            = "GjKLRmDSCLALj28519q8XwKTmJTfFpobEsWCCKWHhzut"
	sample_encryptionPublicKey  = "2Pb4bkbk1oXTpypzuZFABUtDgDP8VCZEiBVuXJsVbgYb"
	sample_seed                 = "74fd93fdd7508e6b2fc9f4e1ac8cef727003f2c36e3f1acf1fcb104658da8f42"
	sample_signKey              = "8sgZQSCiu8GHTveWb1mfusT1KbaFCksriHoXhyDzwahF"
	sample_encryptionPrivateKey = "8sgZQSCiu8GHTveWb1mfusT1KbaFCksriHoXhyDzwahF"
	sample_address              = "dx01nyx8wn3qelmdpykjcjqnw22zdmu9pjt9us0y73"
	sample_pub                  = "dx0pub1zcjduepqaxmxmerk2lw76qxlsf8cc0rzst5hfgy4a3xmvnaxgmkczc9hv30sh4yj2y"
)

/**
{
  "did": "did:dxp:2ou16SbYWkAKDKwUfpQbZX",
  "verifyKey": "2gFhgUd59Ki3aP9dQhMCjR6sphHSuUxR65U4xxbxUT6227FQy78teaDmDt",
  "encryptionPublicKey": "SpeEfV7pvZsW9rs2FRZTGm",
  "secret": {
    "seed": "bddd61295f0a62a1bdfef1537e37c6061beeadbec6e833e6b7115da85711a258",
    "signKey": "0EA56632FB7323E6C7B36D2506AD852ED11B5EF545DC4B5AA2F3B9C3E7421EBA",
    "encryptionPrivateKey": "E1B0F79B20B38E6EAD1D4441BEF7FA3CC4EEF9A1B3B7CCD200A5344F9D2BD3000EBC71A252B38E6EAD1D4441BEF7FA3CC4EEF9A1B3B7CCD200A5344F9D2BD3000EBC71A252"
  }
}
*/
func getKeybase(transient bool, buf io.Reader) (keys.Keybase, error) {
	if transient {
		return keys.NewInMemory(), nil
	}
	return keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)
	if start < 0 || end > length || start > end {
		return ""
	}
	if start == 0 && end == length {
		return source
	}
	return string(r[start:end])
}

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, true)

	return app, ctx
}
func setPrefix() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()
}
