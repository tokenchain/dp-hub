package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tokenchain/dp-hub/x/did"
	"github.com/tokenchain/dp-hub/x/did/exported"
	"github.com/tokenchain/dp-hub/x/oracles"
	"github.com/tokenchain/dp-hub/x/treasury/internal/types"
)

type Keeper struct {
	cdc           *codec.Codec
	storeKey      sdk.StoreKey
	bankKeeper    bank.Keeper
	oraclesKeeper oracles.Keeper
	supplyKeeper  supply.Keeper
	didKeeper     did.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper bank.Keeper,
	oraclesKeeper oracles.Keeper, supplyKeeper supply.Keeper, didKeeper did.Keeper) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		bankKeeper:    bankKeeper,
		oraclesKeeper: oraclesKeeper,
		supplyKeeper:  supplyKeeper,
		didKeeper:     didKeeper,
	}
}
func (k Keeper) Send(ctx sdk.Context, fromDid, toDidOrAddr string, amount sdk.Coins) error {
	fromDidDoc, err := k.didKeeper.GetDidDoc(ctx, fromDid)
	if err != nil {
		fmt.Println("error occurred: ", err)
		return err
	}
	fromAddress := fromDidDoc.Address()
	toAddress, err := k.StringToDx0Addr(ctx, toDidOrAddr)
	if err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoins(ctx, fromAddress, toAddress, amount); err != nil {
		return err
	}
	return nil
}
func (k Keeper) OracleTransfer(ctx sdk.Context, fromDid exported.Did, toDidOrAddr string, oracleDid exported.Did, amount sdk.Coins) error {
	// Check if oracle exists
	if !k.oraclesKeeper.OracleExists(ctx, oracleDid) {
		return exported.IntErr("oracle specified is not a registered oracle")
	}

	// Confirm that oracle has the required capabilities
	oracle := k.oraclesKeeper.MustGetOracle(ctx, oracleDid)
	for _, c := range amount {
		if !oracle.Capabilities.Includes(c.Denom) {
			return exported.IntErr(fmt.Sprintf(
				"oracle does not have capability to send %s", c.Denom))
		}

		// Get capability by token name
		capability := oracle.Capabilities.MustGet(c.Denom)
		if !capability.Capabilities.Includes(oracles.TransferCap) {
			return exported.IntErr(fmt.Sprintf(
				"oracle does not have capability to send %s", c.Denom))
		}
	}

	// Perform send
	return k.Send(ctx, fromDid, toDidOrAddr, amount)
}
func (k Keeper) OracleMint(ctx sdk.Context, oracleDid exported.Did, toDidOrAddr string, amount sdk.Coins) error {

	toAddress, err := k.StringToDx0Addr(ctx, toDidOrAddr)
	if err != nil {
		return err
	}

	// Check if oracle exists
	if !k.oraclesKeeper.OracleExists(ctx, oracleDid) {
		return exported.IntErr("oracle specified is not a registered oracle")
	}

	// Confirm that oracle has the required capabilities
	oracle := k.oraclesKeeper.MustGetOracle(ctx, oracleDid)
	for _, c := range amount {
		if !oracle.Capabilities.Includes(c.Denom) {
			return exported.IntErr(fmt.Sprintf(
				"oracle does not have capability to mint %s", c.Denom))
		}

		// Get capability by token name
		capability := oracle.Capabilities.MustGet(c.Denom)
		if !capability.Capabilities.Includes(oracles.MintCap) {
			return exported.IntErr(fmt.Sprintf(
				"oracle does not have capability to mint %s", c.Denom))
		}
	}

	// Mint coins to module account
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
		return err
	}

	// Send minted tokens to recipient
	if err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddress, amount); err != nil {
		return err
	}
	return nil
}
func (k Keeper) OracleBurn(ctx sdk.Context, oracleDid, fromDid exported.Did, amount sdk.Coins) error {
	// Get from address
	fromDidDoc, err := k.didKeeper.GetDidDoc(ctx, fromDid)
	if err != nil {
		return err
	}
	fromAddress := fromDidDoc.Address()

	// Check if oracle exists
	if !k.oraclesKeeper.OracleExists(ctx, oracleDid) {
		return exported.IntErr("oracle specified is not a registered oracle")
	}

	// Confirm that oracle has the required capabilities
	oracle := k.oraclesKeeper.MustGetOracle(ctx, oracleDid)
	for _, c := range amount {
		if !oracle.Capabilities.Includes(c.Denom) {
			return exported.IntErr(fmt.Sprintf(
				"oracle does not have capability to burn %s", c.Denom))
		}

		// Get capability by token name
		capability := oracle.Capabilities.MustGet(c.Denom)
		if !capability.Capabilities.Includes(oracles.BurnCap) {
			return exported.IntErr(fmt.Sprintf(
				"oracle does not have capability to burn %s", c.Denom))
		}
	}

	// Take tokens to burn from account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx,
		fromAddress, types.ModuleName, amount); err != nil {
		return err
	}

	// Burn coins from module account
	if err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
		return err
	}

	return nil
}
func (k Keeper) StringToDx0Addr(ctx sdk.Context, unknown_address_string string) (sdk.AccAddress, error) {
	// Get to address
	var toAddress sdk.AccAddress
	if exported.IsValidDid(unknown_address_string) {
		toDidDoc, err := k.didKeeper.GetDidDoc(ctx, unknown_address_string)
		if err != nil {
			return nil, err
		}
		toAddress = toDidDoc.Address()
	} else {
		parsedAddr, err := sdk.AccAddressFromBech32(unknown_address_string)
		if err != nil {
			return nil, exported.IntErr(err.Error())
		}
		toAddress = parsedAddr
	}
	return toAddress, nil
}
