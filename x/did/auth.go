package did

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/ed25519"
	"github.com/tendermint/tendermint/crypto"
	ed25519Keys "github.com/tendermint/tendermint/crypto/ed25519"
	sdkx "github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

func GetPubKeyGetter(keeper Keeper) dap.PubKeyGetter {
	return func(ctx sdk.Context, msg dap.IxoMsg) (pubKey crypto.PubKey, res sdk.Result) {

		// Get signer PubKey
		var pubKeyRaw [ed25519.PublicKeySize]byte
		switch msg := msg.(type) {
		case types.MsgAddDid:
			copy(pubKeyRaw[:], base58.Decode(msg.DidDoc.PubKey))
		default:
			// For the remaining messages, the did is the signer
			didDoc, _ := keeper.GetDidDoc(ctx, msg.GetSignerDid())
			if didDoc == nil {
				return pubKey, sdkx.Unauthorized("Issuer did not found").Result()
			}
			copy(pubKeyRaw[:], base58.Decode(didDoc.GetPubKey()))
		}
		return ed25519Keys.PubKeyEd25519(pubKeyRaw), sdk.Result{}
	}
}
