package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/tokenchain/dp-block/client/utils"
	aute2 "github.com/tokenchain/dp-block/x/did/ante"
	"github.com/tokenchain/dp-block/x/did/internal/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tokenchain/dp-block/x/did/exported"
)

const (
	flagDryRun       = "dry-run"
	flagUserEntropy  = "unsafe-entropy"
	flagInteractive  = "interactive"
	flagRecover      = "recover"
	flagNoBackup     = "no-backup"
	flagAccount      = "account"
	flagIndex        = "index"
	flagMultisig     = "multisig"
	flagNoSort       = "nosort"
	flagHDPath       = "hd-path"
	flagGenerateOnly = "generate-only"
	flagKeyAlgo      = "algo"

	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass      = "12345678"
	mnemonicEntropySize = 256
)

type CommandDo func(cmd *cobra.Command, args []string) error

func runGenerationOffline(cdc *codec.Codec) CommandDo {
	return func(cmd *cobra.Command, args []string) error {
		inBuf := bufio.NewReader(cmd.InOrStdin())

		flags := cmd.Flags()
		//cliCtx := context.NewCLIContext().WithCodec(cdc)

		account := uint32(viper.GetInt(flagAccount))
		index := uint32(viper.GetInt(flagIndex))
		generateOnly, _ := flags.GetBool(flagGenerateOnly)

		useBIP44 := !viper.IsSet(flagHDPath)
		var mnemonic string
		var docCombine exported.IxoDid

		algo := keys.SigningAlgo(viper.GetString(flagKeyAlgo))
		if algo == keys.SigningAlgo("") {
			algo = keys.Secp256k1
		}

		//userEntropy, _ := flags.GetBool(flagUserEntropy)
		isDryRun, _ := flags.GetBool(flagDryRun)
		kb, err := utils.GetKeybase(isDryRun, inBuf)

		//privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)
		didBuilder := exported.NewDidGeneratorBuilder()

		name := args[0]
		_, err = kb.Get(name)
		if err == nil {
			// account exists, ask for user confirmation
			response, err2 := input.GetConfirmation(fmt.Sprintf("Override the existing name %s", name), inBuf)
			if err2 != nil {
				return err2
			}
			if !response {
				return errors.New("aborted, not going to override this name")
			}
		}
		didBuilder = didBuilder.WithName(name).Debug()
		if useBIP44 {
			docCombine = didBuilder.BuildDocBIP44(account, index, "")
		} else {
			docCombine = didBuilder.BuildDocHD(viper.GetString(flagHDPath), "")
		}

		mnemonic = didBuilder.GetMnemonicString()

		docInfo, err := kb.CreateOffline(name, docCombine.FromPubKeyDx0(), keys.Ed25519)
		if err != nil {
			fmt.Println("failed to register key with name: ", name)
		}
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("========🔑 Account Info Save This To a secured place===============================================")
		cmd.Println(docCombine)

		cmd.Println("=======🔑  The passphrase please keep in the secured place. Its an important private key! ")
		cmd.Println(mnemonic)

		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")
		cmd.Println("")

		msg2i := fmt.Sprintf("🔐  Do you want to go ahead and make this on the block? %s please go ahead and make the first transaction and send DAPs to the below account\n%s", docInfo.GetName(), docInfo.GetAddress().String())

		response, err2 := input.GetConfirmation(msg2i, inBuf)
		if err2 != nil {
			return err2
		}
		if !response {
			return errors.New("aborted.")
		}
		msg := types.NewMsgAddDid(docCombine.Did, docCombine.GetPubKey())
		cliCtx := context.NewCLIContext().WithCodec(cdc).WithFromAddress(docCombine.Address())
		preheat := aute2.NewDidTxBuild(cliCtx, msg, docCombine)

		if generateOnly {
			return preheat.DebugTxDecode()
		}

		return preheat.CompleteAndBroadcastTxCLI()
	}
}

/*
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
	cmd.Println(did_document.Address().String())*/
