package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tokenchain/ixo-blockchain/client/utils"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/tokenchain/ixo-blockchain/x/project/internal/keeper"
	"github.com/tokenchain/ixo-blockchain/x/project/internal/types"
)

func GetCmdProjectDoc(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-project-doc [did]",
		Short: "Query ProjectDoc for a DID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			didAddr := args[0]
			key := exported.Did(didAddr)

			res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s/%s", types.QuerierRoute, keeper.QueryProjectDoc, key)
			if err != nil {
				return err
			}

			var projectDoc types.MsgCreateProject
			err = cdc.UnmarshalJSON(res, &projectDoc)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(projectDoc, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}

func GetCmdProjectAccounts(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-project-accounts [did]",
		Short: "Get a Project accounts of a Project by Did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			projectDid := args[0]

			res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s/%s", types.QuerierRoute, keeper.QueryProjectAccounts, projectDid)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Project does not exist")
			}

			var f interface{}
			err = json.Unmarshal(res, &f)
			if err != nil {
				return err
			}
			accMap := f.(map[string]interface{})

			output, err := json.MarshalIndent(accMap, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}

func GetCmdProjectTxs(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-project-txs [project-did]",
		Short: "Get a Project txs for a projectDid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			projectDid := args[0]
			res, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s/%s", types.QuerierRoute, keeper.QueryProjectTx, projectDid)
			if err != nil {
				return err
			}

			var txs []types.WithdrawalInfo
			if len(res) == 0 {
				return errors.New("projectTxs does not exist for a projectDid")
			} else {
				err = cdc.UnmarshalJSON(res, &txs)
				if err != nil {
					return err
				}
			}

			output, err := json.MarshalIndent(txs, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}

func GetParamsRequestHandler(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s", types.QuerierRoute, keeper.QueryParams)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(bz, &params); err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
}
