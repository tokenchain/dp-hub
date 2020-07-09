package cli

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

const (
	flagUserEntropy     = "unsafe-entropy"
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

func RunAccMnemonicCmd(cdc *codec.Codec) CommandDo {
	return func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		cliCtx := context.NewCLIContext().WithCodec(cdc)

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

		privKey := ed25519.NewKeyFromSeed(entropySeed)
		EDprivKey := ed25519.PrivKeyToEdPrivateKey(privKey)
		name := args[0]
		info, err := cliCtx.Keybase.CreateOffline(name, EDprivKey.PubKey(), keys.Ed25519)
		if err != nil {
			return err
		}

		mnemonic, err := bip39.NewMnemonic(entropySeed)
		if err != nil {
			return err
		}
		//	cmd.Println(mnemonic)
		did_document := exported.MnToDid(mnemonic)

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

		cmd.Println("=======🔑 Private key. Its an important private key! ")
		cmd.Println(EDprivKey.String())
		cmd.Println("===========================================================================================")

		cmd.Println("=======🔑 The passphrase please keep in the secured place. Its an important private key! ")
		cmd.Println(mnemonic)
		cmd.Println("===========================================================================================")

		cmd.Println("=======🔑 generated a new DID document with the above passphrase. Its an important private key! ")
		cmd.Println(did_document.String())

		cmd.Println("=======💳 DID account address: the additional Darkpool identity card")
		cmd.Println(did_document.DidAddress())

		cmd.Println("=======💳 key address for Darkpool wallet address.")
		cmd.Println(did_document.Address().String())

		return nil
	}
}
