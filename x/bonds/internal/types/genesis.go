package types

type GenesisState struct {
	Bonds   []Bond  `json:"bonds" yaml:"bonds"`
	Batches []Batch `json:"batches" yaml:"batches"`
	Params  Params  `json:"params" yaml:"params"`
}

func NewGenesisState(bonds []Bond, batches []Batch) GenesisState {
	return GenesisState{
		Bonds:   bonds,
		Batches: batches,
	}
}

//noinspection GoUnusedParameter
func ValidateGenesis(data GenesisState) error {
	return ValidateParams(data.Params)
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Bonds:   []Bond{},
		Batches: []Batch{},
		Params:  DefaultParams(),
	}
}

/*
func NewBondDollar() Bond {
	return NewBond(
		"GOLD",
		"United Reserve of Gold",
		"The equal value to the gold price - the commondity value",


		)
}
*/
/*

func ValidateGenesis(data GenesisState) error {
	currentId := store.ZeroEntityID

	for _, asset := range data.Assets {
		currentId = currentId.Inc()
		if !currentId.Equals(asset.ID) {
			return errors.New("Invalid Asset: ID must monotonically increase.")
		}
		if asset.Name == "" {
			return errors.New("Invalid Asset: Must specify a name.")
		}
		if asset.Symbol == "" {
			return errors.New("Invalid Asset: Must specify a symbol.")
		}
		if asset.TotalSupply.IsZero() {
			return errors.New("Invalid Asset: Must specify a non-zero total supply.")
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Assets: []types.Asset{
			{
				ID:                store.NewEntityID(1),
				Name:              "UEX Staking Token",
				Symbol:            "UEX",
				CirculatingSupply: sdk.NewUintFromString("40000000000000000000000000"),
				TotalSupply:       sdk.NewUintFromString("1000000000000000000000000000"),
			},
			{
				ID:                store.NewEntityID(2),
				Name:              "Test Token",
				Symbol:            "TEST",
				CirculatingSupply: sdk.NewUintFromString("40000000000000000000000000"),
				TotalSupply:       sdk.NewUintFromString("1000000000000000000000000000"),
			},
			{
				ID:                store.NewEntityID(3),
				Name:              "Reward Token",
				Symbol:            "AHH",
				CirculatingSupply: sdk.NewUintFromString("40000000000000000000000000"),
				TotalSupply:       sdk.NewUintFromString("1000000000000000000000000000"),
			},
		},
	}
}*/
