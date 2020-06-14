package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

func GetCmdBondDoc(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "getBondDoc [did]",
		Short: "Query BondDoc for a DID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().
				WithCodec(cdc)

			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide a did")
			}

			didAddr := args[0]
			key := ixo.Did(didAddr)

			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute,
				keeper.QueryBondDoc, key), nil)
			if err != nil {
				return err
			}

			var bondDoc types.MsgCreateBond
			err = cdc.UnmarshalJSON(res, &bondDoc)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(bondDoc, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}
