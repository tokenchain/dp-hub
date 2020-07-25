package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/did/ante"
	"github.com/tokenchain/ixo-blockchain/x/treasury/internal/types"
)

func GetCmdSend(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [to-did-or-address] [amount] [sender-dap-did-full]",
		Short: "Create and sign a send tx using DIDs",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			toDidOrAddr := args[0]
			coinsStr := args[1]
			fromDxp := args[2]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			dxpDID, err := did.UnmarshalIxoDid(fromDxp)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithFromAddress(dxpDID.Address())
			msg := types.NewMsgSend(toDidOrAddr, coins, dxpDID.Did)
			return ante.NewDidTxBuild(cliCtx, msg, dxpDID).CompleteAndBroadcastTxCLI()
			//return dap.GenerateOrBroadcastMsgs(cliCtx, msg, ixoDid)
		},
	}
}

func GetCmdOracleTransfer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-transfer [from-did] [to-did-or-addr] [amount] [oracle-dap-did] [proof]",
		Short: "Create and sign an oracle-transfer tx using DIDs",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDid := args[0]
			toDidOrAddr := args[1]
			coinsStr := args[2]
			ixoDidStr := args[3]
			proof := args[4]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			ixoDid, err := did.UnmarshalIxoDid(ixoDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(ixoDid.Address())

			msg := types.NewMsgOracleTransfer(fromDid, toDidOrAddr, coins, ixoDid.Did, proof)

			//return dap.GenerateOrBroadcastMsgs(cliCtx, msg, ixoDid)
			return ante.NewDidTxBuild(cliCtx, msg, ixoDid).CompleteAndBroadcastTxCLI()
		},
	}
}

func GetCmdOracleMint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-mint [to-did-or-addr] [amount] [oracle-dap-did] [proof]",
		Short: "Create and sign an oracle-mint tx using DIDs",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			toDidOrAddr := args[0]
			coinsStr := args[1]
			ixoDidStr := args[2]
			proof := args[3]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			ixoDid, err := did.UnmarshalIxoDid(ixoDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(ixoDid.Address())

			msg := types.NewMsgOracleMint(toDidOrAddr, coins, ixoDid.Did, proof)

			//return dap.GenerateOrBroadcastMsgs(cliCtx, msg, ixoDid)

			return ante.NewDidTxBuild(cliCtx, msg, ixoDid).CompleteAndBroadcastTxCLI()

		},
	}
}

func GetCmdOracleBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-burn [from-did] [amount] [oracle-dap-did] [proof]",
		Short: "Create and sign an oracle-burn tx using DIDs",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDid := args[0]
			coinsStr := args[1]
			ixoDidStr := args[2]
			proof := args[3]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			ixoDid, err := did.UnmarshalIxoDid(ixoDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(ixoDid.Address())

			msg := types.NewMsgOracleBurn(fromDid, coins, ixoDid.Did, proof)

			//return dap.GenerateOrBroadcastMsgs(cliCtx, msg, ixoDid)

			return ante.NewDidTxBuild(cliCtx, msg, ixoDid).CompleteAndBroadcastTxCLI()
		},
	}
}
