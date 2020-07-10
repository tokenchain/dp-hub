package cli

import (
	"bufio"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"

	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
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

	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass      = "12345678"
	mnemonicEntropySize = 256
)

type CommandDo func(cmd *cobra.Command, args []string) error

func RunMnemonicCmd(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	userEntropy, _ := flags.GetBool(flagUserEntropy)
	var entropySeed []byte
	if userEntropy {
		// prompt the user to enter some entropy
		buf := bufio.NewReader(cmd.InOrStdin())
		inputEntropy, err := input.GetString("> WARNING: Generate at least 256-bits of entropy and enter the results here:", buf)
		if err != nil {
			return err
		}
		if len(inputEntropy) < 43 {
			return fmt.Errorf("256-bits is 43 characters in Base-64, and 100 in Base-6. You entered %v, and probably want more", len(inputEntropy))
		}
		conf, err := input.GetConfirmation(fmt.Sprintf("> Input length: %d", len(inputEntropy)), buf)
		if err != nil {
			return err
		}
		if !conf {
			return nil
		}

		// hash input entropy to get entropy seed
		hashedEntropy := sha256.Sum256([]byte(inputEntropy))
		entropySeed = hashedEntropy[:]
	} else {
		// read entropy seed straight from crypto.Rand
		var err error
		entropySeed, err = bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			return err
		}
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return err
	}
	//	cmd.Println(mnemonic)
	cmd.Println("=======🔑 The passphrase please keep in the secured place. Its an important private key! ")
	cmd.Println(mnemonic)
	cmd.Println("===========================================================================================")
	did_document := exported.MnToDid(mnemonic)

	cmd.Println("=======🔑 generated a new DID document with the above passphrase. Its an important private key! ")
	cmd.Println(did_document.String())

	cmd.Println("=======💳 DID account address: the additional Darkpool identity card")
	cmd.Println(did_document.DidAddress())

	cmd.Println("=======💳 key address for Darkpool wallet address.")
	cmd.Println(did_document.Address().String())
	return nil
}

func getKeybase(transient bool, buf io.Reader) (keys.Keybase, error) {
	if transient {
		return keys.NewInMemory(), nil
	}
	return keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
}

func RunAccMnemonicCmd(cdc *codec.Codec) CommandDo {
	return func(cmd *cobra.Command, args []string) error {
		inBuf := bufio.NewReader(cmd.InOrStdin())

		flags := cmd.Flags()
		//cliCtx := context.NewCLIContext().WithCodec(cdc)

		account := uint32(viper.GetInt(flagAccount))
		index := uint32(viper.GetInt(flagIndex))

		useBIP44 := !viper.IsSet(flagHDPath)
		var hdPath string

		if useBIP44 {
			hdPath = keys.CreateHDPath(account, index).String()
		} else {
			hdPath = viper.GetString(flagHDPath)
		}

		algo := keys.SigningAlgo(viper.GetString(flagKeyAlgo))
		if algo == keys.SigningAlgo("") {
			algo = keys.Secp256k1
		}

		//userEntropy, _ := flags.GetBool(flagUserEntropy)
		isDryRun, _ := flags.GetBool(flagDryRun)
		kb, err := getKeybase(isDryRun, inBuf)

		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		cmd.Println("========Seed====================================================")
		mnemonic, err := bip39.NewMnemonic(entropySeed)
		// create master key and derive first key for keyring
		derivedPriv, err := keys.StdDeriveKey(mnemonic, "", hdPath, algo)
		if err != nil {
			return err
		}
		privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)
		if err != nil {
			return err
		}
		/*
			var info Info
			if encryptPasswd != "" {
					info = keyWriter.writeLocalKey(name, privKey, encryptPasswd, algo)
			} else {
					info = kb.writeOfflineKey(keyWriter, name, privKey.PubKey(), algo)
			}*/
		//signKey := base58.Encode(privateKey[:32])
		//signKey := base58.Encode(entropySeed)
		//cmd.Println("========Sign Key====================================================")
		//cmd.Println(signKey)
		//privA, private_key := ed25519.NewKeyFromSeedToEdPrivateKey(entropySeed)
		//privKey :=NewKeyFromSeed(entropySeed)
		//EDprivKey := PrivKeyToEdPrivateKey(privKey)
		//cmd.Println("========Private Key====================================================")
		//cmd.Println(string(private_key))
		name := args[0]
		info, err := kb.CreateOffline(name, privKey.PubKey(), algo)

		if err != nil {
			return err
		}

		//	cmd.Println(mnemonic)
		did_document := exported.InfoToDidEd25519(info, privKey, derivedPriv)

		cmd.Println("========  Account Info===============================================")
		cmd.Println(info)

		cmd.Println("========User Name====================================================")
		cmd.Println(info.GetName())

		cmd.Println("========Public Key====================================================")
		cmd.Println(info.GetPubKey())

		cmd.Println("========Address===================================================")
		cmd.Println(info.GetAddress())

		cmd.Println("========Algo===================================================")
		cmd.Println(info.GetAlgo())

		cmd.Println("========Path===================================================")
		cmd.Println(info.GetPath())

		cmd.Println("========🔑  Derived Private Key===============================")
		cmd.Println(base58.Encode(derivedPriv))

		cmd.Println("========  Public Key====================================================")
		cmd.Println(base58.Encode(privKey.PubKey().Bytes()))
		cmd.Println(info.GetPubKey().Address().String())

		cmd.Println("========  Address base58====================================================")
		cmd.Println(base58.Encode(info.GetPubKey().Address()))

		cmd.Println("========💳  Address from VerifyKeyToAddr ====================================================")
		cmd.Println(exported.VerifyKeyToAddr(did_document.VerifyKey).String())

		cmd.Println("========💳  Address stored darkpool=====")
		cmd.Println(info.GetAddress().String())

		cmd.Println("=======🔑  Private Key. Its an important private key! ")
		cmd.Println(base58.Encode(privKey.Bytes()))

		cmd.Println("=======🔑  The passphrase please keep in the secured place. Its an important private key! ")
		cmd.Println(mnemonic)

		cmd.Println("=======🔑  generated a new DID document with the above passphrase. Its an important private key! ")
		cmd.Println(did_document.String())

		cmd.Println("=======💳  DID account address: the additional Darkpool identity card")
		cmd.Println(did_document.DidAddress())

		cmd.Println("=======💳  key address for Darkpool wallet address.")
		cmd.Println(did_document.Address().String())

		return nil
	}
}
