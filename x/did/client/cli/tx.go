package cli

import (
	"fmt"
	"github.com/tokenchain/dp-hub/x/did/ante"
	"github.com/tokenchain/dp-hub/x/did/exported"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tokenchain/dp-hub/x/did/internal/types"
)

func GetCmdAddDidDoc(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-did-doc [sovrin-did]",
		Short: "Add a new SovrinDid from the full json document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sovrinDid, err := exported.UnmarshalDxpDid(args[0])
			if err != nil {
				return err
			}
			fmt.Println(sovrinDid)
			msg := types.NewMsgAddDid(sovrinDid.Did, sovrinDid.GetPubKey())
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithFromAddress(sovrinDid.Address())
			return ante.NewDidTxBuild(cliCtx, msg, sovrinDid).CompleteAndBroadcastTxCLI()
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
			return ante.NewDidTxBuild(cliCtx, msg, sovrinDid).CompleteAndBroadcastTxCLI()
		},
	}
}

/*
func GetCmdDidGenerate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "generate-doc [username]",
		Short: "Generate new did document for the given name",
		RunE:  RunGenerationNewDoc(cdc),
	}
}
*/
func GetCmdAccDidGenerate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "generate-offline [name]",
		Short: "Generate did document offline",
		RunE:  runGenerationOffline(cdc),
	}
}
