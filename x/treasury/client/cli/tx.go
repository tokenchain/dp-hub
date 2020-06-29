package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	types2 "github.com/tokenchain/ixo-blockchain/x/ixo/types"

	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/treasury/internal/types"
)

func GetCmdSend(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [to-did] [amount] [sender-sovrin-did]",
		Short: "Create and sign a send tx using DIDs",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			toDid := args[0]
			coinsStr := args[1]
			sovrinDidStr := args[2]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(types2.DidToAddr(sovrinDid.Did))

			msg := types.NewMsgSend(toDid, coins, sovrinDid)

			return ixo.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}

func GetCmdOracleTransfer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-transfer [from-did] [to-did] [amount] [oracle-sovrin-did] [proof]",
		Short: "Create and sign an oracle-transfer tx using DIDs",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDid := args[0]
			toDid := args[1]
			coinsStr := args[2]
			sovrinDidStr := args[3]
			proof := args[4]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(types2.DidToAddr(sovrinDid.Did))

			msg := types.NewMsgOracleTransfer(fromDid, toDid, coins, sovrinDid, proof)

			return ixo.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}

func GetCmdOracleMint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-mint [to-did] [amount] [oracle-sovrin-did] [proof]",
		Short: "Create and sign an oracle-mint tx using DIDs",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			toDid := args[0]
			coinsStr := args[1]
			sovrinDidStr := args[2]
			proof := args[3]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(types2.DidToAddr(sovrinDid.Did))

			msg := types.NewMsgOracleMint(toDid, coins, sovrinDid, proof)

			return ixo.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}

func GetCmdOracleBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-burn [from-did] [amount] [oracle-sovrin-did] [proof]",
		Short: "Create and sign an oracle-burn tx using DIDs",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDid := args[0]
			coinsStr := args[1]
			sovrinDidStr := args[2]
			proof := args[3]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid, err := types2.UnmarshalSovrinDid(sovrinDidStr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc).
				WithFromAddress(types2.DidToAddr(sovrinDid.Did))

			msg := types.NewMsgOracleBurn(fromDid, coins, sovrinDid, proof)

			return ixo.SignAndBroadcastTxCli(cliCtx, msg, sovrinDid)
		},
	}
}
