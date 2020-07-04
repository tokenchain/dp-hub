package types

type GenesisState struct {
	DidDocs []DidDoc `json:"did_docs" yaml:"did_docs"`
}

func NewGenesisState(didDocs []DidDoc) GenesisState {
	return GenesisState{
		DidDocs: didDocs,
	}
}

//noinspection GoUnusedParameter
func ValidateGenesis(data GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		DidDocs: nil,
	}
}
