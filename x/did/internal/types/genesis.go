package types

import (
	"github.com/tokenchain/dp-hub/x/did/exported"
)

type GenesisState struct {
	DidDocs []exported.DidDoc `json:"did_docs" yaml:"did_docs"`
}

func NewGenesisState(didDocs []exported.DidDoc) GenesisState {
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
