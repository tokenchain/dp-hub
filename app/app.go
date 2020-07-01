package app

import (
	"encoding/json"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/tokenchain/ixo-blockchain/x/bonddoc"
	"github.com/tokenchain/ixo-blockchain/x/bonds"
	"github.com/tokenchain/ixo-blockchain/x/did"
	dap "github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"github.com/tokenchain/ixo-blockchain/x/nameservice"
	"github.com/tokenchain/ixo-blockchain/x/oracles"
	"github.com/tokenchain/ixo-blockchain/x/payments"
	"github.com/tokenchain/ixo-blockchain/x/project"
	"github.com/tokenchain/ixo-blockchain/x/treasury"
	"io"
	"os"
)

const (
	appName              = "Darkpool"
	Bech32MainPrefix     = "dx0"
	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.dxocli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.dxod")

	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsClient.ProposalHandler, distribution.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		did.AppModuleBasic{},
		payments.AppModuleBasic{},
		project.AppModuleBasic{},
		bonddoc.AppModuleBasic{},
		bonds.AppModuleBasic{},
		treasury.AppModuleBasic{},
		oracles.AppModuleBasic{},
		nameservice.AppModule{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:            nil,
		distribution.ModuleName:          nil,
		mint.ModuleName:                  {supply.Minter},
		staking.BondedPoolName:           {supply.Burner, supply.Staking},
		staking.NotBondedPoolName:        {supply.Burner, supply.Staking},
		gov.ModuleName:                   {supply.Burner},
		bonds.BondsMintBurnAccount:       {supply.Minter, supply.Burner},
		bonds.BatchesIntermediaryAccount: nil,
		treasury.ModuleName:              {supply.Minter, supply.Burner},
		nameservice.ModuleName:           {supply.Minter, supply.Burner},
		payments.PayRemainderPool:        nil,
		payments.ModuleName:              nil,
	}

	// Reserved payments module ID prefixes
	paymentsReservedIdPrefixes = []string{}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc.Seal()
}

type DpApp struct {
	*bam.BaseApp
	cdc                *codec.Codec
	invCheckPeriod     uint
	keys               map[string]*sdk.KVStoreKey
	tKeys              map[string]*sdk.TransientStoreKey
	subspaces          map[string]params.Subspace // subspaces
	accountKeeper      auth.AccountKeeper
	bankKeeper         bank.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	distributionKeeper distribution.Keeper
	supplyKeeper       supply.Keeper
	paramsKeeper       params.Keeper

	govKeeper      gov.Keeper
	mintKeeper     mint.Keeper
	crisisKeeper   crisis.Keeper
	didKeeper      did.Keeper
	paymentsKeeper payments.Keeper
	projectKeeper  project.Keeper
	bonddocKeeper  bonddoc.Keeper
	bondsKeeper    bonds.Keeper
	oraclesKeeper  oracles.Keeper
	treasuryKeeper treasury.Keeper
	nsKeeper       nameservice.Keeper

	mm *module.Manager
	sm *module.SimulationManager // simulation manager
}

// verify app interface at compile time
//var _ simapp.App = (*DpApp)(nil)

func NewIxoApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *DpApp {

	cdc := MakeCodec()
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, types.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	//bApp.SetAppVersion(version.NewVersion("v0.1.1"))

	keys := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, distribution.StoreKey, slashing.StoreKey,
		params.StoreKey, gov.StoreKey, did.StoreKey, mint.StoreKey,
		project.StoreKey, bonds.StoreKey, bonddoc.StoreKey, treasury.StoreKey,
		oracles.StoreKey)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &DpApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[distribution.ModuleName] = app.paramsKeeper.Subspace(distribution.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace)
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[payments.ModuleName] = app.paramsKeeper.Subspace(payments.DefaultParamspace)
	app.subspaces[project.ModuleName] = app.paramsKeeper.Subspace(project.DefaultParamspace)

	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount)
	// The BankKeeper allows you perform sdk.Coins interactions
	//app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs(), )
	//app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs())
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs())
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	//stakingKeeper := staking.NewKeeper(app.cdc, keys[staking.StoreKey], app.supplyKeeper, app.subspaces[staking.ModuleName])
	stakingKeeper := staking.NewKeeper(app.cdc, keys[staking.StoreKey], app.supplyKeeper, app.subspaces[staking.ModuleName])
	app.mintKeeper = mint.NewKeeper(app.cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	//app.distributionKeeper = distribution.NewKeeper(app.cdc, keys[distribution.StoreKey], app.subspaces[distribution.ModuleName], &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.distributionKeeper = distribution.NewKeeper(app.cdc, keys[distribution.StoreKey], app.subspaces[distribution.ModuleName], &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	//app.slashingKeeper = slashing.NewKeeper(app.cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName])
	app.slashingKeeper = slashing.NewKeeper(app.cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName])
	app.crisisKeeper = crisis.NewKeeper(app.subspaces[crisis.ModuleName], invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distribution.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(app.distributionKeeper))

	//	app.govKeeper = gov.NewKeeper(app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.supplyKeeper, &stakingKeeper, govRouter)
	app.govKeeper = gov.NewKeeper(app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.supplyKeeper, &stakingKeeper, govRouter)

	app.stakingKeeper = *stakingKeeper.SetHooks(staking.NewMultiStakingHooks(app.distributionKeeper.Hooks(),
		app.slashingKeeper.Hooks()))

	app.didKeeper = did.NewKeeper(app.cdc, keys[did.StoreKey])
	app.paymentsKeeper = payments.NewKeeper(app.cdc, keys[payments.StoreKey], app.subspaces[payments.ModuleName] , app.bankKeeper, paymentsReservedIdPrefixes)
	app.projectKeeper = project.NewKeeper(app.cdc, keys[project.StoreKey], app.subspaces[project.ModuleName] , app.accountKeeper, app.paymentsKeeper)
	app.bonddocKeeper = bonddoc.NewKeeper(app.cdc, keys[bonddoc.StoreKey])
	app.bondsKeeper = bonds.NewKeeper(app.bankKeeper, app.supplyKeeper, app.accountKeeper, app.stakingKeeper, keys[bonds.StoreKey], app.cdc)
	app.oraclesKeeper = oracles.NewKeeper(app.cdc, keys[oracles.StoreKey])
	app.treasuryKeeper = treasury.NewKeeper(app.cdc, keys[treasury.StoreKey], app.bankKeeper, app.oraclesKeeper, app.supplyKeeper)
	app.nsKeeper = nameservice.NewKeeper(app.cdc, keys[nameservice.StoreKey], app.bankKeeper)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distribution.NewAppModule(app.distributionKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),

		did.NewAppModule(app.didKeeper),
		payments.NewAppModule(app.paymentsKeeper, app.bankKeeper),
		project.NewAppModule(app.projectKeeper, app.paymentsKeeper, app.bankKeeper),
		bonddoc.NewAppModule(app.bonddocKeeper),
		bonds.NewAppModule(app.bondsKeeper, app.accountKeeper),
		treasury.NewAppModule(app.treasuryKeeper),
		oracles.NewAppModule(app.oraclesKeeper),
	)

	app.mm.SetOrderBeginBlockers(mint.ModuleName, distribution.ModuleName, slashing.ModuleName, bonds.ModuleName)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName, bonds.ModuleName)

	app.mm.SetOrderInitGenesis(
		distribution.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		supply.ModuleName,
		crisis.ModuleName,
		genutil.ModuleName,
		did.ModuleName,
		project.ModuleName,
		payments.ModuleName,
		bonddoc.ModuleName,
		bonds.ModuleName,
		treasury.ModuleName,
		oracles.ModuleName,
		nameservice.ModuleName)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	app.mm.RegisterInvariants(&app.crisisKeeper)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(initAnteHandler(app))

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *DpApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *DpApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *DpApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState map[string]json.RawMessage
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *DpApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *DpApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *DpApp) Codec() *codec.Codec {
	return app.cdc
}
// SimulationManager implements the SimulationApp interface
func (app *DpApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *DpApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	if forZeroHeight {
		app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}
	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//      in favour of export at a block height
func (app *DpApp) prepForZeroHeightGenesis(ctx sdk.Context, jailWhiteList []string) {
	applyWhiteList := false

	//Check if there is a whitelist
	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err.Error())
		}
		whiteListMap[addr] = true
	}

	/* Handle fee distribution state. */

	// withdraw all validator commission
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
		_, err := app.distributionKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		if err != nil {
			panic(err.Error())
		}
		return false
	})

	// withdraw all delegator rewards
	dels := app.stakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, err := app.distributionKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
		if err != nil {
			panic(err.Error())
		}
	}

	// clear validator slash events
	app.distributionKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.distributionKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {

		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.distributionKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool := app.distributionKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.distributionKeeper.SetFeePool(ctx, feePool)

		app.distributionKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	// reinitialize all delegations
	for _, del := range dels {
		app.distributionKeeper.Hooks().BeforeDelegationCreated(ctx, del.DelegatorAddress, del.ValidatorAddress)
		app.distributionKeeper.Hooks().AfterDelegationModified(ctx, del.DelegatorAddress, del.ValidatorAddress)
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	app.stakingKeeper.IterateRedelegations(ctx, func(_ int64, red staking.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.stakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	app.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd staking.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.stakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[staking.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, staking.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := app.stakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		app.stakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_ = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	/* Handle slashing state. */

	// reset start height on signing infos
	app.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashing.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.slashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}

func initAnteHandler(app *DpApp) sdk.AnteHandler {
	didPubKeyGetter := did.GetPubKeyGetter(app.didKeeper)
	projectPubKeyGetter := project.GetPubKeyGetter(app.projectKeeper, app.didKeeper)
	bonddocPubKeyGetter := bonddoc.GetPubKeyGetter(app.bonddocKeeper)
	bondsPubKeyGetter := bonds.GetPubKeyGetter(app.bondsKeeper, app.didKeeper)
	treasuryPubKeyGetter := treasury.GetPubKeyGetter(app.didKeeper)
	paymentsPubKeyGetter := payments.GetPubKeyGetter(app.didKeeper)
	cosmosAnteHandler := auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer)
	didAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, didPubKeyGetter)
	projectAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, projectPubKeyGetter)
	bonddocAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, bonddocPubKeyGetter)
	bondsAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, bondsPubKeyGetter)
	treasuryAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, treasuryPubKeyGetter)
	//feesAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, paymentsPubKeyGetter)
	paymentsAnteHandler := dap.NewAnteHandler(app.accountKeeper, app.supplyKeeper, paymentsPubKeyGetter)

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (_ sdk.Context, _ error) {
		msg := tx.GetMsgs()[0]
		switch msg.Route() {
		case did.RouterKey:
			return didAnteHandler(ctx, tx, simulate)
		case project.RouterKey:
			return projectAnteHandler(ctx, tx, simulate)
		case bonddoc.RouterKey:
			return bonddocAnteHandler(ctx, tx, simulate)
		case bonds.RouterKey:
			return bondsAnteHandler(ctx, tx, simulate)
		case treasury.RouterKey:
			return treasuryAnteHandler(ctx, tx, simulate)
		case payments.RouterKey:
			return paymentsAnteHandler(ctx, tx, simulate)
		default:
			return cosmosAnteHandler(ctx, tx, simulate)
		}
	}
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

/*case fees.RouterKey:
return feesAnteHandler(ctx, tx, simulate)
*/
