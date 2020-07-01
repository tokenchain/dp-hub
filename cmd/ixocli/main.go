package main

import (
	"github.com/cosmos/cosmos-sdk/client/flags"

	//	"github.com/cosmos/cosmos-sdk/x/mint"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authCli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankCli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tmlibs/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/tokenchain/ixo-blockchain/app"
	cli2 "github.com/tokenchain/ixo-blockchain/client/cli"
	tx2 "github.com/tokenchain/ixo-blockchain/client/tx"
	/*	distRest "github.com/cosmos/cosmos-sdk/x/distribution/client/rest"
		distcmd "github.com/cosmos/cosmos-sdk/x/distribution"
		distClient "github.com/cosmos/cosmos-sdk/x/distribution/client"*/

)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "dxocli",
		Short: "dxo Light-Client",
	}

	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
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
		Use:   "tx",
		Aliases: []string{"tx"},
		Short: "Transactions subcommands",
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

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	tx2.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	//distRest.RegisterRoutes(rs.CliCtx, rs.Mux, rs., distcmd.StoreKey)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}
