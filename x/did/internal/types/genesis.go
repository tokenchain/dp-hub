package types

import (
	"github.com/tokenchain/dp-hub/x/ixo/types"
)

type GenesisState struct {
	DidDocs []types.DidDoc `json:"did_docs" yaml:"did_docs"`
}

func NewGenesisState(didDocs []types.DidDoc) GenesisState {
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
