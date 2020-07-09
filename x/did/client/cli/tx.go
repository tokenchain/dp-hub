package cli

import (
	"fmt"
	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

func GetCmdAddDidDoc(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-did-doc [sovrin-did]",
		Short: "Add a new SovrinDid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sovrinDid, err := exported.UnmarshalDxpDid(args[0])
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithFromAddress(sovrinDid.Address())
			msg := types.NewMsgAddDid(sovrinDid.Did, sovrinDid.VerifyKey)
			return dap.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}

func GetCmdAddCredential(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-kyc-credential [did] [signer-did-doc]",
		Short: "Add a new KYC Credential for a Did by the signer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			didAddr := args[0]

			sovrinDid, err := exported.UnmarshalDxpDid(args[1])
			if err != nil {
				return err
			}
			fmt.Println("Confirmed its a valid did document... ")
			t := time.Now()
			issued := t.Format(time.RFC3339)
			credTypes := []string{"Credential", "ProofOfKYC"}
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithFromAddress(sovrinDid.Address())

			msg := types.NewMsgAddCredential(didAddr, credTypes, sovrinDid.Did, issued)
			return dap.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}

func GetCmdDidGenerate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "generate-did-doc",
		Short: "Query all DID documents",
		RunE:  RunMnemonicCmd,
	}
}

func GetCmdAccDidGenerate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "generate-did-acc [name]",
		Short: "Query all DID documents",
		RunE:  RunAccMnemonicCmd(cdc),
	}
}
