package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genAccsCli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	genUtilCli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	oraclesCli "github.com/tokenchain/ixo-blockchain/x/oracles/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/tokenchain/ixo-blockchain/app"
)

const (
	flagInvCheckPeriod = "inv-check-period"
)

var (
	invCheckPeriod uint
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "dxod",
		Short:             "dxo Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(
		genUtilCli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genUtilCli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genUtilCli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{},
			app.DefaultNodeHome, app.DefaultCLIHome),
		genUtilCli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		genAccsCli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		oraclesCli.AddGenesisOracleCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
	)

	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "DXO", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewIxoApp(logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64,
	forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		nsApp := app.NewIxoApp(logger, db, traceStore, false, uint(2))
		err := nsApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	nsApp := app.NewIxoApp(logger, db, traceStore, true, uint(2))

	return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
