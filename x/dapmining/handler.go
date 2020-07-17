package dapmining

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {

		default:
			errMsg := fmt.Sprintf("Unrecognized bonds Msg type: %v", msg.Type())
			return nil, exported.UnknownRequest(errMsg)
		}
	}
}
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
func InitGenesis(ctx sdk.Context, k Keeper, s GenesisState) {

}
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	s := GenesisState{}
	return s
}
