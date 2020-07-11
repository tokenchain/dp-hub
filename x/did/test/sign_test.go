package test

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"strings"
	"testing"
)


type addrData struct {
	Mnemonic string
	Master   string
	Seed     string
	Priv     string
	Pub      string
	Addr     string
}

func Test_dk(t *testing.T) {

	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	println("========Seed====================================================")
	mnemonic, err := bip39.NewMnemonic(entropySeed)

	//account := uint32(viper.GetInt(flagAccount))
	//index := uint32(viper.GetInt(flagIndex))

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
	require.Nil(t, err, "KB pass ")
	mnemonic = sample_mnom
	// create master key and derive first key for keyring
	derivedPriv, err := keys.StdDeriveKey(mnemonic, "", hdPath, algo)
	require.Nil(t, err, "pass 1 ")
	privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)
	require.Nil(t, err, "pass 2 ")
	/*
		var info Info
		if encryptPasswd != "" {
				info = keyWriter.writeLocalKey(name, privKey, encryptPasswd, algo)
		} else {
				info = kb.writeOfflineKey(keyWriter, name, privKey.PubKey(), algo)
		}*/
	//signKey := base58.Encode(privateKey[:32])
	//signKey := base58.Encode(entropySeed)
	//println("========Sign Key====================================================")
	//println(signKey)
	//privA, private_key := ed25519.NewKeyFromSeedToEdPrivateKey(entropySeed)
	//privKey :=NewKeyFromSeed(entropySeed)
	//EDprivKey := PrivKeyToEdPrivateKey(privKey)
	//println("========Private Key====================================================")
	//println(string(private_key))

	name := "cosmos"
	info, err := kb.CreateOffline(name, privKey.PubKey(), keys.Secp256k1)

	require.Nil(t, err, "offline creation pass 3")
	//	println(mnemonic)
	did_document := exported.InfoToDid(info, privKey, algo)
	println("========💳  Address from VerifyKeyToAddr ====================================================")
	println(exported.VerifyKeyToAddr(did_document.VerifyKey).String())
	println("========💳  Address stored darkpool=====")
	println(info.GetAddress().String())

	println("========  Account Info===============================================")
	println(info)

	println("========User Name====================================================")
	println(info.GetName())

	println("========Public Key====================================================")
	println(info.GetPubKey())

	println("========Address===================================================")
	println(info.GetAddress().String())

	println("========Algo===================================================")
	println(info.GetAlgo())

	println("========Path===================================================")
	println(info)
	println("========🔑  Derived Private Key====================================================")
	println(base58.Encode(derivedPriv))
	println("========  Public Key====================================================")
	println(base58.Encode(privKey.PubKey().Bytes()))
	println(base58.Encode(info.GetPubKey().Bytes()[:]))

	println(info.GetPubKey().Address().String())

	println("========  Address base58====================================================")
	println(base58.Encode(info.GetPubKey().Address()))

	println("=======🔑  Private Key. Its an important private key! ")
	println(base58.Encode(privKey.Bytes()))
	println("===========================================================================================")
	println("=======🔑  The passphrase please keep in the secured place. Its an important private key! ")
	println(mnemonic)
	println("===========================================================================================")

	println("=======🔑  generated a new DID document with the above passphrase. Its an important private key! ")
	println(did_document.String())

	println("=======💳  DID account address: the additional Darkpool identity card")
	println(did_document.DidAddress())

	println("=======💳  key address for Darkpool wallet address.")
	println(did_document.Address().String())

}



func Test_recover(t *testing.T) {
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
	mnemonic = sample_mnom
	// create master key and derive first key for keyring
	derivedPriv, err := keys.StdDeriveKey(mnemonic, "", hdPath, algo)
	privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)

	info, err := kb.CreateOffline(name, privKey.PubKey(), algo)

	did_document := exported.InfoToDid(info, privKey, algo)
	recover_privKey:=exported.RecoverDidSecpK1ToPrivateKey(did_document)

	//var recover_privKey ed25519tm.PrivKeyEd25519
	fmt.Println(recover_privKey)
	require.Equal(t, 32, len(recover_privKey), "is now the same")
}

func Test_development(t *testing.T) {
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
	mnemonic = sample_mnom
	// create master key and derive first key for keyring
	derivedPriv, err := keys.StdDeriveKey(mnemonic, "", hdPath, algo)
	privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)

	info, err := kb.CreateOffline(name, privKey.PubKey(), algo)

	did_document := exported.InfoToDid(info, privKey, algo)

	//var recover_privKey ed25519tm.PrivKeyEd25519
	var recover_privKey secp256k1.PrivKeySecp256k1
	p1, _ := hex.DecodeString(strings.ToLower(did_document.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(did_document.Secret.SignKey))
	privKey_orginal := exported.PrivateKeyToSecp256k1(privKey)
	privkey_v3 := exported.SecpPrivKey(derivedPriv)

	fmt.Println("========start result  ===================================================")
	fmt.Println(did_document)
	fmt.Println(did_document.Secret.EncryptionPrivateKey)
	fmt.Println(p1, len(p1))
	fmt.Println(did_document.Secret.SignKey)
	fmt.Println(p2, len(p2))

	copy(recover_privKey[:], p1)
	copy(recover_privKey[24:], p2)


	fmt.Println("========secp original  =========")
	fmt.Println(privkey_v3)
	fmt.Println(len(privkey_v3))

	fmt.Println("========key original  =========")
	fmt.Println(privKey_orginal)
	fmt.Println(len(privKey_orginal))

	fmt.Println("========key recover  =========")
	fmt.Println(recover_privKey)
	fmt.Println(len(recover_privKey))

	fmt.Println("========direct private key =========")
	fmt.Println(privKey.Bytes())
	fmt.Println(len(privKey.Bytes()))

	fmt.Println(len(base58.Decode(did_document.Secret.EncryptionPrivateKey)))
	fmt.Println(len(base58.Decode(did_document.Secret.SignKey)))
	//fmt.Println(privKey.Equals(recover_privKey))
	//fmt.Println(privKey.PubKey().Address().String())

	fmt.Println("========ENd===================================================")
	//require.Equal(t, len(privKey.Bytes()), len(recover_privKey), "recover key are the same")
	//	require.Equal(t, len(privKey.Bytes()), len(recover_privKey.Bytes()), "recover key success")
	require.Equal(t, privKey, recover_privKey, "recover key success")
}

