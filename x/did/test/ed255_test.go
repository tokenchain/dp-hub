package test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"testing"
)

func Test_recover_ke(t *testing.T) {
	setPrefix()
	name := "cosmos"
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	require.Nil(t, err, "memory phrase generated")
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	did_document := exported.NewDidGeneratorBuilder().WithMem(mnemonic).WithName(name).Build()
	fmt.Println("========🔑 Account Info Save This To a secured place===============================================")
	fmt.Println(did_document)

	fmt.Println("=======🔑  The passphrase please keep in the secured place. Its an important private key! ")
	fmt.Println(mnemonic)
}
func Test_recover_public_key(t *testing.T) {
	setPrefix()
	doc := exported.NewDapDid(
		sample_did,
		sample_verifyKey,
		sample_encryptionPublicKey,
		sample_seed,
		sample_signKey,
		sample_encryptionPrivateKey,
		sample_address,
		sample_pub,
		"cosmos",
	)

	//err := dap.SignAndBroadcastTxCli()

	/*	var recover_pub [32]byte
		name := substring(doc.Did, 8, len(doc.Did))
		p1 := base58.Decode(name)
		p2 := base58.Decode(doc.EncryptionPublicKey)
		copy(recover_pub[:], p1)
		copy(recover_pub[16:], p2)*/
	/*
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
		config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
		config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
		config.Seal()*/
	//both parts are filled
	r := exported.RecoverDidToEd25519PubKey(doc)
	//only the first part filled
	//f := exported.RecoverDidEd25519ToPrivateKey(doc)
	//get the private key
	fmt.Println(doc.FromAddressDx0().String())
	fmt.Println(doc.Address().String())

	fmt.Println("========key recover  pubkey from pub key=========")
	fmt.Println(doc.FromPubKeyDx0())
	fmt.Println(doc.GetPubKey())

	fmt.Println("========key recover  pubkey=========")
	fmt.Println(doc.GetPubKeyByte())

	fmt.Println("========key recover  private key=========")
	fmt.Println(doc.GetPriKeyByte())
	pri := doc.GetPriKeyByte()
	fmt.Println(len(pri), pri)
	require.Equal(t, r, ed25519.PrivKeyEd25519(doc.GetPriKeyByte()))
	fmt.Println("========it is now name  =========")
	//	fmt.Println(name, doc.Address().String())

	fmt.Println("========key recover  =========")
	fmt.Println(len(r), r)

	fmt.Println("========key recover  =========")

	fmt.Println("========key recover  =========")
	require.Nil(t, nil, "successfull signed")
}
func Test_ed25519_development(t *testing.T) {
	setPrefix()

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

	fmt.Println("====== memoric  ============")
	fmt.Println(mnemonic)

	// create master key and derive first key for keyring
	derivedPriv, err := keys.StdDeriveKey(mnemonic, "", hdPath, algo)
	privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)
	info, err := kb.CreateOffline(name, privKey.PubKey(), algo)

	docCombine := exported.InfoToDidEd25519(info, derivedPriv)

	//	privKey_orginal := exported.PrivateKeyToSecp256k1(privKey)
	//	privkey_v3 := exported.SecpPrivKey(derivedPriv)
	//privateRecover := exported.RecoverDidToPrivateKeyClassic(docCombine)
	//cosmosPrivateKey := exported.RecoverDidToCosmosPrivateKey(docCombine)

	fmt.Println("========key recover  =========")
	fmt.Println(docCombine)
	fmt.Println(docCombine.Dpinfo.DpAddress)
	//fmt.Println(privateRecover)

	fmt.Println(info.GetAddress().String())
	fmt.Println(docCombine.Dpinfo.PubKey)
	//fmt.Println(len(recover_priv_key_ed))

	fmt.Println("========cosmos private key check  =========")

	fmt.Println(len(derivedPriv), derivedPriv)
	fmt.Println(len(privKey.Bytes()), privKey)
	//fmt.Println(len(cosmosPrivateKey), cosmosPrivateKey.PubKey().Address().Bytes())
	//fmt.Println(privKey.Equals(secp256k1.PrivKeySecp256k1(cosmosPrivateKey)))

	fmt.Println("========end game now  ===========================")
	//require.Equal(t, cap(privKey.Bytes()), cap(recover_priv_key_ed), "recover key are the same")
	//	require.Equal(t, cap(privKey.Bytes()), cap(recover_priv_key_ed.Bytes()), "recover key success")
	//require.Equal(t, 64, len(privateRecover), "recover key success")
	//require.Equal(t, privKey, cosmosPrivateKey, "recover private cosmos key success")
	//require.Equal(t, info.GetAddress().String(), docCombine.DpInfo.GetAddress().String(), "keybase account established")

}
