package payments

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/tokenchain/ixo-blockchain/x/payments/client/cli"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tokenchain/ixo-blockchain/x/payments/client/rest"
	"github.com/tokenchain/ixo-blockchain/x/payments/internal/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	paymentsTxCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "payments transaction sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	paymentsTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdCreatePaymentTemplate(cdc),
		cli.GetCmdCreatePaymentContract(cdc),
		cli.GetCmdCreateSubscription(cdc),
		cli.GetCmdSetPaymentContractAuthorisation(cdc),
		cli.GetCmdGrantPaymentDiscount(cdc),
		cli.GetCmdRevokePaymentDiscount(cdc),
		cli.GetCmdEffectPayment(cdc),
	)...)

	return paymentsTxCmd
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	paymentsQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "payments query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	paymentsQueryCmd.AddCommand(client.GetCommands(
		cli.GetParamsRequestHandler(cdc),
		cli.GetCmdPaymentTemplate(cdc),
		cli.GetCmdPaymentContract(cdc),
		cli.GetCmdSubscription(cdc),
	)...)

	return paymentsQueryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper     keeper.Keeper
	bankKeeper bank.Keeper
}

func NewAppModule(keeper Keeper, bankKeeper bank.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		bankKeeper:     bankKeeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

func (AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper, am.bankKeeper)
}

func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}
