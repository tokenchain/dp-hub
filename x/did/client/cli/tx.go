package cli

import (
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdAddDidDoc(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-did-doc [sovrin-did]",
		Short: "Add a new DxpDid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sovrinDid, err := did.UnmarshalDxpDid(args[0])
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(did.DidToAddr(sovrinDid.Did))

			msg := types.NewMsgAddDid(sovrinDid.Did, sovrinDid.VerifyKey)
			return ixo.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
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

			sovrinDid, err := did.UnmarshalDxpDid(args[1])
			if err != nil {
				return err
			}

			t := time.Now()
			issued := t.Format(time.RFC3339)

			credTypes := []string{"Credential", "ProofOfKYC"}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(did.DidToAddr(sovrinDid.Did))

			msg := types.NewMsgAddCredential(didAddr, credTypes, sovrinDid.Did, issued)
			return ixo.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}
