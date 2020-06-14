package cli

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/ixo/sovrin"
	"github.com/tokenchain/ixo-blockchain/x/treasury/internal/types"
)

func IxoSignAndBroadcast(cdc *codec.Codec, ctx context.CLIContext, msg sdk.Msg, sovrinDid sovrin.SovrinDid) error {
	privKey := [64]byte{}
	copy(privKey[:], base58.Decode(sovrinDid.Secret.SignKey))
	copy(privKey[32:], base58.Decode(sovrinDid.VerifyKey))

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	signature := ixo.SignIxoMessage(msgBytes, sovrinDid.Did, privKey)
	tx := ixo.NewIxoTxSingleMsg(msg, signature)

	bz, err := cdc.MarshalJSON(tx)
	if err != nil {
		panic(err)
	}

	res, err := ctx.BroadcastTx(bz)
	if err != nil {
		return err
	}

	fmt.Println(res.String())
	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.TxHash)
	return nil

}

func unmarshalSovrinDID(sovrinJson string) sovrin.SovrinDid {
	sovrinDid := sovrin.SovrinDid{}
	sovrinErr := json.Unmarshal([]byte(sovrinJson), &sovrinDid)
	if sovrinErr != nil {
		panic(sovrinErr)
	}

	return sovrinDid
}

func GetCmdSend(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [to-did] [amount] [sender-sovrin-did]",
		Short: "Create and sign a send tx using DIDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().
				WithCodec(cdc)

			if len(args) != 3 || len(args[0]) == 0 ||
				len(args[1]) == 0 || len(args[2]) == 0 {
				return errors.New("You must provide the receiver DID, " +
					"amount, and sender private key")
			}

			toDid := args[0]
			coinsStr := args[1]
			sovrinDidStr := args[2]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid := unmarshalSovrinDID(sovrinDidStr)
			msg := types.NewMsgSend(toDid, coins, sovrinDid)

			return IxoSignAndBroadcast(cdc, ctx, msg, sovrinDid)
		},
	}
}

func GetCmdOracleTransfer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-transfer [from-did] [to-did] [amount] [oracle-sovrin-did] [proof]",
		Short: "Create and sign an oracle-transfer tx using DIDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().
				WithCodec(cdc)

			if len(args) != 5 || len(args[0]) == 0 || len(args[1]) == 0 ||
				len(args[2]) == 0 || len(args[3]) == 0 || len(args[4]) == 0 {
				return errors.New("You must provide the sender and receiver " +
					"DID, amount, sender private key, and proof")
			}

			fromDid := args[0]
			toDid := args[1]
			coinsStr := args[2]
			sovrinDidStr := args[3]
			proof := args[4]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid := unmarshalSovrinDID(sovrinDidStr)
			msg := types.NewMsgOracleTransfer(fromDid, toDid, coins, sovrinDid, proof)

			return IxoSignAndBroadcast(cdc, ctx, msg, sovrinDid)
		},
	}
}

func GetCmdOracleMint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-mint [to-did] [amount] [oracle-sovrin-did] [proof]",
		Short: "Create and sign an oracle-mint tx using DIDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().
				WithCodec(cdc)

			if len(args) != 4 || len(args[0]) == 0 || len(args[1]) == 0 ||
				len(args[2]) == 0 || len(args[3]) == 0 {
				return errors.New("You must provide the recipient DID, " +
					"amount, oracle private key, and proof")
			}

			toDid := args[0]
			coinsStr := args[1]
			sovrinDidStr := args[2]
			proof := args[3]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid := unmarshalSovrinDID(sovrinDidStr)
			msg := types.NewMsgOracleMint(toDid, coins, sovrinDid, proof)

			return IxoSignAndBroadcast(cdc, ctx, msg, sovrinDid)
		},
	}
}

func GetCmdOracleBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracle-burn [from-did] [amount] [oracle-sovrin-did] [proof]",
		Short: "Create and sign an oracle-burn tx using DIDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().
				WithCodec(cdc)

			if len(args) != 4 || len(args[0]) == 0 || len(args[1]) == 0 ||
				len(args[2]) == 0 || len(args[3]) == 0 {
				return errors.New("You must provide the source DID, " +
					"amount, oracle private key, and proof")
			}

			fromDid := args[0]
			coinsStr := args[1]
			sovrinDidStr := args[2]
			proof := args[3]

			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			sovrinDid := unmarshalSovrinDID(sovrinDidStr)
			msg := types.NewMsgOracleBurn(fromDid, coins, sovrinDid, proof)

			return IxoSignAndBroadcast(cdc, ctx, msg, sovrinDid)
		},
	}
}
