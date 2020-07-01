package bonds

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/bonds/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	types2 "github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

func GetPubKeyGetter(keeper Keeper, didKeeper did.Keeper) ixo.PubKeyGetter {
	return func(ctx sdk.Context, msg types2.IxoMsg) ([32]byte, error) {

		// Get signer PubKey and sender DID
		var pubKey [32]byte
		var senderDid types2.Did
		switch msg := msg.(type) {
		case types.MsgCreateBond:
			senderDid = msg.CreatorDid
			copy(pubKey[:], base58.Decode(msg.CreatorPubKey))
		case types.MsgEditBond:
			senderDid = msg.EditorDid
			copy(pubKey[:], base58.Decode(msg.EditorPubKey))
		case types.MsgBuy:
			senderDid = msg.BuyerDid
			copy(pubKey[:], base58.Decode(msg.PubKey))
		case types.MsgSell:
			senderDid = msg.SellerDid
			copy(pubKey[:], base58.Decode(msg.PubKey))
		case types.MsgSwap:
			senderDid = msg.SwapperDid
			copy(pubKey[:], base58.Decode(msg.PubKey))
		default:
			return pubKey, x.UnknownRequest("No match for message type")
		}

		// Check that sender's DID is ledgered
		senderDidDoc, _ := didKeeper.GetDidDoc(ctx, senderDid)
		if senderDidDoc == nil {
			return pubKey, x.Unauthorized("Sender did not found")
		}

		return pubKey, nil
	}
}
