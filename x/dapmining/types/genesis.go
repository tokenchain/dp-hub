package types

import (
	"fmt"
	"time"
)

type GenesisState struct {
	RewardTime   time.Time `json:"snapshot_date" yaml:"snapshot_date"`
	TotalRewards uint64    `json:"total_rewards" yaml:"total_rewards"`
}

func NewGenesisState() GenesisState {
	return GenesisState{
		RewardTime:   time.Now(),
		TotalRewards: 0,
	}
}

//noinspection GoUnusedParameter
func ValidateGenesis(data GenesisState) error {
	if data.TotalRewards < 0 {
		return fmt.Errorf("Impossible to negative reward accounts %d! ", data.TotalRewards)
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		RewardTime:   time.Now(),
		TotalRewards: 0,
	}
}

/*

type GenesisState struct {
	WhoisRecords []Whois `json:"whois_records"`
}

func NewGenesisState(whoIsRecords []Whois) GenesisState {
	return GenesisState{WhoisRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WhoisRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
		if record.Price == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords: []Whois{},
	}
}*/
