package did

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/dap"
	types2 "github.com/tokenchain/ixo-blockchain/x/dap/types"
)

func GetPubKeyGetter(keeper Keeper) dap.PubKeyGetter {
	return func(ctx sdk.Context, msg types2.IxoMsg) ([32]byte, error) {

		// Get signer PubKey
		var pubKey [32]byte
		switch msg := msg.(type) {
		case types.MsgAddDid:
			copy(pubKey[:], base58.Decode(msg.DidDoc.PubKey))
		case types.MsgAddCredential:
			did := msg.GetSignerDid()
			didDoc, _ := keeper.GetDidDoc(ctx, did)
			if didDoc == nil {
				return pubKey, x.Unauthorized("Issuer did not found")
			}
			copy(pubKey[:], base58.Decode(didDoc.GetPubKey()))
		default:
			return pubKey, x.UnknownRequest("No match for message type.")
		}
		return pubKey, nil
	}
}
