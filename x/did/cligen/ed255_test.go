package cligen

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"testing"
)

func Test_veriEd25519(t *testing.T) {
	name := "cosmos"
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	require.Nil(t, err, "memory phrase generated")
	println("========Seed====================================================")
	mnemonic, err := bip39.NewMnemonic(entropySeed)

	var hdPath string
	useBIP44 := !viper.IsSet(flagHDPath)
	if useBIP44 {
		hdPath = keys.CreateHDPath(0, 1).String()
	} else {
		hdPath = viper.GetString(flagHDPath)
	}

	algo := keys.SigningAlgo(viper.GetString(flagKeyAlgo))
	if algo == keys.SigningAlgo("") {
		algo = keys.Secp256k1
	}
	//isDryRun, _ := flags.GetBool(flagDryRun)
	kb, err := getKeybase(true, nil)
	//mnemonic = mnem

	fmt.Println("====== memoric")
	fmt.Println(mnemonic)

	// create master key and derive first key for keyring
	derivedPriv, err := keys.StdDeriveKey(mnemonic, "", hdPath, algo)
	privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)
	info, err := kb.CreateOffline(name, privKey.PubKey(), algo)
	//did_document := exported.InfoToDid(info, privKey, algo)
	did_ed_doc := exported.InfoToDidEd25519(info, privKey, derivedPriv)

	//	privKey_orginal := exported.PrivateKeyToSecp256k1(privKey)
	//	privkey_v3 := exported.SecpPrivKey(derivedPriv)
	privateRecover := exported.RecoverDidEd25519ToPrivateKey(did_ed_doc)

	fmt.Println("========key recover  =========")
	fmt.Println(privateRecover)
	//fmt.Println(len(recover_priv_key_ed))

	fmt.Println("========ENd===================================================")
	//require.Equal(t, cap(privKey.Bytes()), cap(recover_priv_key_ed), "recover key are the same")
	//	require.Equal(t, cap(privKey.Bytes()), cap(recover_priv_key_ed.Bytes()), "recover key success")
	require.Equal(t, 64, len(privateRecover), "recover key success")

}
