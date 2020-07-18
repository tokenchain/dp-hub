package app


import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestIxodExport(t *testing.T) {
	db := db.NewMemDB()
	ixoApp := NewDarkpoolApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	setGenesis(ixoApp)

	// Making a new app object with the db, so that initchain hasn't been called
	NewDarkpoolApp := NewDarkpoolApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err := NewDarkpoolApp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := db.NewMemDB()
	app := NewDarkpoolApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	for acc := range maccPerms {
		require.True(t, app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(ixoApp *DpApp) error {

	genesisState := ModuleBasics.DefaultGenesis() // TODO: set to simapp.DefaultGenesisState() once simapp added
	stateBytes, err := codec.MarshalJSONIndent(ixoApp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	ixoApp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	ixoApp.Commit()
	return nil
}

