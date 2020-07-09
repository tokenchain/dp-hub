package did

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
)

func GetPubKeyGetter(keeper Keeper) types.PubKeyGetter {
	return func(ctx sdk.Context, msg dap.IxoMsg) (pubKey crypto.PubKey, res error) {
		// Get signer PubKey
		var pubKeyEd25519 ed25519.PubKeyEd25519
		switch msg := msg.(type) {
		case MsgAddDid:
			copy(pubKeyEd25519[:], base58.Decode(msg.DidDoc.PubKey))
		default:
			// For the remaining messages, the did is the signer
			didDoc, _ := keeper.GetDidDoc(ctx, msg.GetSignerDid())
			if didDoc == nil {
				return pubKey, Unauthorized("Issuer did not found")
			}
			copy(pubKeyEd25519[:], base58.Decode(didDoc.GetPubKey()))
		}

		return pubKeyEd25519, nil
	}
}

func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}
