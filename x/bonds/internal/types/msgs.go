package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	"strings"
)

const (
	TypeMsgCreateBond = "create_bond"
	TypeMsgEditBond   = "edit_bond"
	TypeMsgBuy        = "buy"
	TypeMsgSell       = "sell"
	TypeMsgSwap       = "swap"
	TypeMsgBurn       = "burn"
	TypeMsgMint       = "mint"
	TypeMsgTransfer   = "transfer"
)

type (
	MsgCreateBond struct {
		BondDid                types.Did      `json:"bond_did" yaml:"bond_did"`
		Token                  string         `json:"token" yaml:"token"`
		Name                   string         `json:"name" yaml:"name"`
		Description            string         `json:"description" yaml:"description"`
		FunctionType           string         `json:"function_type" yaml:"function_type"`
		FunctionParameters     FunctionParams `json:"function_parameters" yaml:"function_parameters"`
		CreatorDid             types.Did      `json:"creator_did" yaml:"creator_did"`
		CreatorPubKey          string         `json:"pub_key" yaml:"pub_key"`
		ReserveTokens          []string       `json:"reserve_tokens" yaml:"reserve_tokens"`
		TxFeePercentage        sdk.Dec        `json:"tx_fee_percentage" yaml:"tx_fee_percentage"`
		ExitFeePercentage      sdk.Dec        `json:"exit_fee_percentage" yaml:"exit_fee_percentage"`
		FeeAddress             sdk.AccAddress `json:"fee_address" yaml:"fee_address"`
		MaxSupply              sdk.Coin       `json:"max_supply" yaml:"max_supply"`
		OrderQuantityLimits    sdk.Coins      `json:"order_quantity_limits" yaml:"order_quantity_limits"`
		SanityRate             sdk.Dec        `json:"sanity_rate" yaml:"sanity_rate"`
		SanityMarginPercentage sdk.Dec        `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
		AllowSells             string         `json:"allow_sells" yaml:"allow_sells"`
		BatchBlocks            sdk.Uint       `json:"batch_blocks" yaml:"batch_blocks"`
	}
	MsgEditBond struct {
		BondDid                types.Did `json:"bond_did" yaml:"bond_did"`
		Token                  string    `json:"token" yaml:"token"`
		Name                   string    `json:"name" yaml:"name"`
		Description            string    `json:"description" yaml:"description"`
		OrderQuantityLimits    string    `json:"order_quantity_limits" yaml:"order_quantity_limits"`
		SanityRate             string    `json:"sanity_rate" yaml:"sanity_rate"`
		SanityMarginPercentage string    `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
		EditorDid              types.Did `json:"editor_did" yaml:"editor_did"`
		EditorPubKey           string    `json:"pub_key" yaml:"pub_key"`
	}
	MsgBuy struct {
		BuyerDid  types.Did `json:"buyer_did" yaml:"buyer_did"`
		PubKey    string    `json:"pub_key" yaml:"pub_key"`
		Amount    sdk.Coin  `json:"amount" yaml:"amount"`
		MaxPrices sdk.Coins `json:"max_prices" yaml:"max_prices"`
		BondDid   types.Did `json:"bond_did" yaml:"bond_did"`
	}

	MsgSwap struct {
		SwapperDid types.Did `json:"swapper_did" yaml:"swapper_did"`
		PubKey     string    `json:"pub_key" yaml:"pub_key"`
		BondDid    types.Did `json:"bond_did" yaml:"bond_did"`
		From       sdk.Coin  `json:"from" yaml:"from"`
		ToToken    string    `json:"to_token" yaml:"to_token"`
	}

	MsgMint struct {
		ID     types.Did      `json:"minter_did" yaml:"minter_did"`
		Minter sdk.AccAddress `json:"minter_address" yaml:"minter_address"`
		Amount sdk.Uint       `json:"amount" yaml:"amount"`
	}

	MsgBurn struct {
		ID     types.Did      `json:"burner_did" yaml:"burner_did"`
		Burner sdk.AccAddress `json:"burner_address" yaml:"burner_address"`
		Amount sdk.Uint       `json:"amount" yaml:"amount"`
	}

	MsgTransfer struct {
		ID     types.Did      `json:"transfer_did" yaml:"transfer_did"`
		From   sdk.AccAddress `json:"from_address" yaml:"from_address"`
		To     sdk.AccAddress `json:"to_address" yaml:"to_address"`
		Amount sdk.Uint       `json:"amount" yaml:"amount"`
	}
)

var (
	_ types.IxoMsg = MsgCreateBond{}
	_ types.IxoMsg = MsgEditBond{}
	_ types.IxoMsg = MsgBuy{}
	_ types.IxoMsg = MsgSell{}
	_ types.IxoMsg = MsgSwap{}
	_ types.IxoMsg = MsgMint{}
	_ types.IxoMsg = MsgBurn{}
	_ types.IxoMsg = MsgTransfer{}
)

func NewMsgCreateBond(token, name, description string, creatorDid types.SovrinDid,
	functionType string, functionParameters FunctionParams, reserveTokens []string,
	txFeePercentage, exitFeePercentage sdk.Dec, feeAddress sdk.AccAddress, maxSupply sdk.Coin,
	orderQuantityLimits sdk.Coins, sanityRate, sanityMarginPercentage sdk.Dec,
	allowSell string, batchBlocks sdk.Uint, bondDid types.Did) MsgCreateBond {
	return MsgCreateBond{
		BondDid:                bondDid,
		Token:                  token,
		Name:                   name,
		Description:            description,
		CreatorDid:             creatorDid.Did,
		CreatorPubKey:          creatorDid.VerifyKey,
		FunctionType:           functionType,
		FunctionParameters:     functionParameters,
		ReserveTokens:          reserveTokens,
		TxFeePercentage:        txFeePercentage,
		ExitFeePercentage:      exitFeePercentage,
		FeeAddress:             feeAddress,
		MaxSupply:              maxSupply,
		OrderQuantityLimits:    orderQuantityLimits,
		SanityRate:             sanityRate,
		SanityMarginPercentage: sanityMarginPercentage,
		AllowSells:             strings.ToLower(allowSell),
		BatchBlocks:            batchBlocks,
	}
}

func (msg MsgCreateBond) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.BondDid) == "" {
		return x.ErrArgumentCannotBeEmpty("BondDid")
	} else if strings.TrimSpace(msg.Token) == "" {
		return x.ErrArgumentCannotBeEmpty("Token")
	} else if strings.TrimSpace(msg.Name) == "" {
		return x.ErrArgumentCannotBeEmpty("Name")
	} else if strings.TrimSpace(msg.Description) == "" {
		return x.ErrArgumentCannotBeEmpty("Description")
	} else if strings.TrimSpace(msg.CreatorDid) == "" {
		return x.ErrArgumentCannotBeEmpty("CreatorDid")
	} else if strings.TrimSpace(msg.CreatorPubKey) == "" {
		return x.ErrArgumentCannotBeEmpty("CreatorPubKey")
	} else if len(msg.ReserveTokens) == 0 {
		return x.ErrArgumentCannotBeEmpty("Reserve token")
	} else if msg.FeeAddress.Empty() {
		return x.ErrArgumentCannotBeEmpty("Fee address")
	} else if strings.TrimSpace(msg.FunctionType) == "" {
		return x.ErrArgumentCannotBeEmpty("Function type")
	} else if strings.TrimSpace(msg.AllowSells) == "" {
		return x.ErrArgumentCannotBeEmpty("AllowSells")
	}
	// Note: FunctionParameters can be empty

	// Check that bond token is a valid token name
	err := CheckCoinDenom(msg.Token)
	if err != nil {
		return x.ErrInvalidCoinDenomination(msg.Token)
	}

	// Validate function parameters
	if err := msg.FunctionParameters.Validate(msg.FunctionType); err != nil {
		return err
	}

	// Validate reserve tokens
	if err = CheckReserveTokenNames(msg.ReserveTokens, msg.Token); err != nil {
		return err
	} else if err = CheckNoOfReserveTokens(msg.ReserveTokens, msg.FunctionType); err != nil {
		return err
	}

	// Validate coins
	if !msg.MaxSupply.IsValid() {
		return x.IntErr("max supply is invalid")
	} else if !msg.OrderQuantityLimits.IsValid() {
		return x.IntErr("order quantity limits are invalid")
	}

	// Check that max supply denom matches token denom
	if msg.MaxSupply.Denom != msg.Token {
		return x.ErrMaxSupplyDenomDoesNotMatchTokenDenom()
	}

	// Check that Sanity values not negative
	if msg.SanityRate.IsNegative() {
		return x.ErrArgumentCannotBeNegative("SanityRate")
	} else if msg.SanityMarginPercentage.IsNegative() {
		return x.ErrArgumentCannotBeNegative("SanityMarginPercentage")
	}

	// Check that true or false
	if msg.AllowSells != TRUE && msg.AllowSells != FALSE {
		return x.ErrArgumentMissingOrNonBoolean("AllowSells")
	}

	// Check FeePercentages not negative and don't add up to 100
	if msg.TxFeePercentage.IsNegative() {
		return x.ErrArgumentCannotBeNegative("TxFeePercentage")
	} else if msg.ExitFeePercentage.IsNegative() {
		return x.ErrArgumentCannotBeNegative("ExitFeePercentage")
	} else if msg.TxFeePercentage.Add(msg.ExitFeePercentage).GTE(sdk.NewDec(100)) {
		return errors.Wrap(x.ErrFeesCannotBeOrExceed100Percent, "Sum of fees is or exceeds 100 percent")
	}

	// Check that not zero
	if msg.BatchBlocks.IsZero() {
		return x.ErrArgumentMustBePositive("BatchBlocks")
	} else if msg.MaxSupply.Amount.IsZero() {
		return x.ErrArgumentMustBePositive("MaxSupply")
	}

	// Note: uniqueness of reserve tokens checked when parsing

	// Check that DIDs valid
	if !types.IsValidDid(msg.BondDid) {
		return x.ErrInvalidDid("bond did is invalid")
	} else if !types.IsValidDid(msg.CreatorDid) {
		return x.ErrInvalidDid("creator did is invalid")
	}

	return nil
}

func (msg MsgCreateBond) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgCreateBond) GetSignerDid() types.Did { return msg.CreatorDid }
func (msg MsgCreateBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgCreateBond) Route() string { return RouterKey }

func (msg MsgCreateBond) Type() string { return TypeMsgCreateBond }

func NewMsgEditBond(token, name, description, orderQuantityLimits, sanityRate,
	sanityMarginPercentage string, editorDid types.SovrinDid, bondDid types.Did) MsgEditBond {
	return MsgEditBond{
		BondDid:                bondDid,
		Token:                  token,
		Name:                   name,
		Description:            description,
		OrderQuantityLimits:    orderQuantityLimits,
		SanityRate:             sanityRate,
		SanityMarginPercentage: sanityMarginPercentage,
		EditorDid:              editorDid.Did,
		EditorPubKey:           editorDid.VerifyKey,
	}
}

func (msg MsgEditBond) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.BondDid) == "" {
		return x.ErrArgumentCannotBeEmpty("BondDid")
	} else if strings.TrimSpace(msg.Token) == "" {
		return x.ErrArgumentCannotBeEmpty("Token")
	} else if strings.TrimSpace(msg.Name) == "" {
		return x.ErrArgumentCannotBeEmpty("Name")
	} else if strings.TrimSpace(msg.Description) == "" {
		return x.ErrArgumentCannotBeEmpty("Description")
	} else if strings.TrimSpace(msg.SanityRate) == "" {
		return x.ErrArgumentCannotBeEmpty("SanityRate")
	} else if strings.TrimSpace(msg.SanityMarginPercentage) == "" {
		return x.ErrArgumentCannotBeEmpty("SanityMarginPercentage")
	} else if strings.TrimSpace(msg.EditorDid) == "" {
		return x.ErrArgumentCannotBeEmpty("EditorDid")
	} else if strings.TrimSpace(msg.EditorPubKey) == "" {
		return x.ErrArgumentCannotBeEmpty("EditorPubKey")
	}
	// Note: order quantity limits can be blank

	// Check that at least one editable was edited. Fields that will not
	// be edited should be "DoNotModifyField", and not an empty string
	inputList := []string{
		msg.Name, msg.Description, msg.OrderQuantityLimits,
		msg.SanityRate, msg.SanityMarginPercentage,
	}
	atLeaseOneEdit := false
	for _, e := range inputList {
		if e != DoNotModifyField {
			atLeaseOneEdit = true
			break
		}
	}
	if !atLeaseOneEdit {
		return ErrDidNotEditAnything()
	}

	// Check that DIDs valid
	if !types.IsValidDid(msg.BondDid) {
		return x.ErrInvalidDid("bond did is invalid")
	} else if !types.IsValidDid(msg.EditorDid) {
		return x.ErrInvalidDid("editor did is invalid")
	}

	return nil
}

func (msg MsgEditBond) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgEditBond) GetSignerDid() types.Did { return msg.EditorDid }
func (msg MsgEditBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgEditBond) Route() string { return RouterKey }

func (msg MsgEditBond) Type() string { return TypeMsgEditBond }

func NewMsgBuy(buyerDid types.SovrinDid, amount sdk.Coin, maxPrices sdk.Coins,
	bondDid types.Did) MsgBuy {
	return MsgBuy{
		BuyerDid:  buyerDid.Did,
		PubKey:    buyerDid.VerifyKey,
		Amount:    amount,
		MaxPrices: maxPrices,
		BondDid:   bondDid,
	}
}

func (msg MsgBuy) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.BuyerDid) == "" {
		return x.ErrArgumentCannotBeEmpty("BuyerDid")
	} else if strings.TrimSpace(msg.PubKey) == "" {
		return x.ErrArgumentCannotBeEmpty("PubKey")
	} else if strings.TrimSpace(msg.BondDid) == "" {
		return x.ErrArgumentCannotBeEmpty("BondDid")
	}

	// Check that amount valid and non zero
	if !msg.Amount.IsValid() {
		return x.IntErr("amount is invalid")
	} else if msg.Amount.Amount.IsZero() {
		return x.ErrArgumentMustBePositive("Amount")
	}

	// Check that maxPrices valid
	if !msg.MaxPrices.IsValid() {
		return x.IntErr("maxprices is invalid")
	}

	// Check that DIDs valid
	if !types.IsValidDid(msg.BondDid) {
		return x.ErrInvalidDid("bond did is invalid")
	} else if !types.IsValidDid(msg.BuyerDid) {
		return x.ErrInvalidDid("buyer did is invalid")
	}

	return nil
}

func (msg MsgBuy) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgBuy) GetSignerDid() types.Did { return msg.BuyerDid }
func (msg MsgBuy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgBuy) Route() string { return RouterKey }

func (msg MsgBuy) Type() string { return TypeMsgBuy }

type MsgSell struct {
	SellerDid types.Did `json:"seller_did" yaml:"seller_did"`
	PubKey    string    `json:"pub_key" yaml:"pub_key"`
	Amount    sdk.Coin  `json:"amount" yaml:"amount"`
	BondDid   types.Did `json:"bond_did" yaml:"bond_did"`
}

func NewMsgSell(sellerDid types.SovrinDid, amount sdk.Coin, bondDid types.Did) MsgSell {
	return MsgSell{
		SellerDid: sellerDid.Did,
		PubKey:    sellerDid.VerifyKey,
		Amount:    amount,
		BondDid:   bondDid,
	}
}

func (msg MsgSell) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.SellerDid) == "" {
		return x.ErrArgumentCannotBeEmpty("SellerDid")
	} else if strings.TrimSpace(msg.PubKey) == "" {
		return x.ErrArgumentCannotBeEmpty("PubKey")
	} else if strings.TrimSpace(msg.BondDid) == "" {
		return x.ErrArgumentCannotBeEmpty("BondDid")
	}

	// Check that amount valid and non zero
	if !msg.Amount.IsValid() {
		return x.IntErr("amount is invalid")
	} else if msg.Amount.Amount.IsZero() {
		return x.ErrArgumentMustBePositive("Amount")
	}

	// Check that DIDs valid
	if !types.IsValidDid(msg.BondDid) {
		return x.ErrInvalidDid("bond did is invalid")
	} else if !types.IsValidDid(msg.SellerDid) {
		return x.ErrInvalidDid("seller did is invalid")
	}

	return nil
}

func (msg MsgSell) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgSell) GetSignerDid() types.Did { return msg.SellerDid }
func (msg MsgSell) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgSell) Route() string { return RouterKey }

func (msg MsgSell) Type() string { return TypeMsgSell }

func NewMsgSwap(swapperDid types.SovrinDid, from sdk.Coin, toToken string,
	bondDid types.Did) MsgSwap {
	return MsgSwap{
		SwapperDid: swapperDid.Did,
		PubKey:     swapperDid.VerifyKey,
		From:       from,
		ToToken:    toToken,
		BondDid:    bondDid,
	}
}

func (msg MsgSwap) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.SwapperDid) == "" {
		return x.ErrArgumentCannotBeEmpty("SwapperDid")
	} else if strings.TrimSpace(msg.PubKey) == "" {
		return x.ErrArgumentCannotBeEmpty("PubKey")
	} else if strings.TrimSpace(msg.BondDid) == "" {
		return x.ErrArgumentCannotBeEmpty("BondDid")
	} else if strings.TrimSpace(msg.ToToken) == "" {
		return x.ErrArgumentCannotBeEmpty("ToToken")
	}

	// Validate from amount
	if !msg.From.IsValid() {
		return x.IntErr("from amount is invalid")
	}

	// Validate to token
	err := CheckCoinDenom(msg.ToToken)
	if err != nil {
		return err
	}

	// Check if from and to the same token
	if msg.From.Denom == msg.ToToken {
		return ErrFromAndToCannotBeTheSameToken()
	}

	// Check that non zero
	if msg.From.Amount.IsZero() {
		return x.ErrArgumentMustBePositive("FromAmount")
	}

	// Note: From denom and amount must be valid since sdk.Coin

	// Check that DIDs valid
	if !types.IsValidDid(msg.BondDid) {
		return x.ErrInvalidDid("bond did is invalid")
	} else if !types.IsValidDid(msg.SwapperDid) {
		return x.ErrInvalidDid("swapper did is invalid")
	}

	return nil
}

func (msg MsgSwap) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgSwap) GetSignerDid() types.Did { return msg.SwapperDid }
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}

func (msg MsgSwap) Route() string { return RouterKey }

func (msg MsgSwap) Type() string { return TypeMsgSwap }

func NewMsgTransfer(id types.Did, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Uint) MsgTransfer {
	return MsgTransfer{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (msg MsgTransfer) ValidateBasic() error    { return nil }
func (msg MsgTransfer) GetSignerDid() types.Did { return msg.ID }
func (msg MsgTransfer) Type() string            { return TypeMsgTransfer }
func (msg MsgTransfer) Route() string           { return RouterKey }
func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}
func (msg MsgTransfer) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func NewMsgBurn(id types.Did, from sdk.AccAddress, amount sdk.Uint) MsgBurn {
	return MsgBurn{
		ID:     id,
		Burner: from,
		Amount: amount,
	}
}

func (msg MsgBurn) ValidateBasic() error    { return nil }
func (msg MsgBurn) GetSignerDid() types.Did { return msg.ID }
func (msg MsgBurn) Type() string            { return TypeMsgBurn }
func (msg MsgBurn) Route() string           { return RouterKey }
func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}
func (msg MsgBurn) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func NewMsgMint(id types.Did, from sdk.AccAddress, amount sdk.Uint) MsgMint {
	return MsgMint{
		ID:     id,
		Minter: from,
		Amount: amount,
	}
}

func (msg MsgMint) ValidateBasic() error    { return nil }
func (msg MsgMint) GetSignerDid() types.Did { return msg.ID }
func (msg MsgMint) Type() string            { return TypeMsgMint }
func (msg MsgMint) Route() string           { return RouterKey }
func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}
func (msg MsgMint) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}
