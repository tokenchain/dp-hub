package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authCli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankCli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tmlibs/cli"
	"github.com/tokenchain/dp-hub/app"
	cli2 "github.com/tokenchain/dp-hub/client/cli"
	clientTx "github.com/tokenchain/dp-hub/client/tx"
	/*authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	distRest "github.com/cosmos/cosmos-sdk/x/distribution/client/rest"
	distcmd "github.com/cosmos/cosmos-sdk/x/distribution"
	distClient "github.com/cosmos/cosmos-sdk/x/distribution/client"*/

)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	config := sdk.GetConfig()
	config.SetCoinType(app.CoinType)
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "dpcli",
		Short: "dp Light-Client",
	}

	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return cli2.InitConfig(rootCmd)
	}

	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		lcd.ServeCommand(cdc, registerRoutes),
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)

	//mc := []sdk.{
	/*nsclient.NewModuleClient(storeNS, cdc),
	pricingclient.NewModuleClient(storePricing, cdc),
	stakingclient.NewModuleClient(st.StoreKey, cdc),*/
	//	distClient.NewModuleClient(distcmd.StoreKey, cdc),
	/*slashingclient.NewModuleClient(sl.StoreKey, cdc),
	mintclient.NewModuleClient(mint.StoreKey, cdc),*/
	//}

	executor := cli.PrepareMainCmd(rootCmd, "DXO", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authCli.GetAccountCmd(cdc),
		flags.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authCli.QueryTxsByEventsCmd(cdc),
		flags.LineBreak,
		cli2.QueryTxCmd(cdc),
	)

	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:     "tx",
		Aliases: []string{"tx"},
		Short:   "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankCli.SendTxCmd(cdc),
		flags.LineBreak,
		authCli.GetSignCommand(cdc),
		authCli.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authCli.GetBroadcastCommand(cdc),
		authCli.GetEncodeCommand(cdc),
		flags.LineBreak,
	)
	app.ModuleBasics.AddTxCommands(txCmd, cdc)
	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	//authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	clientTx.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}
