package treasury

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"

	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

func GetPubKeyGetter(didKeeper did.Keeper) ixo.PubKeyGetter {
	return func(ctx sdk.Context, msg types.IxoMsg) ([32]byte, sdk.Result) {

		// Get signer PubKey
		var pubKey [32]byte
		switch msg := msg.(type) {
		case MsgSend:
			copy(pubKey[:], base58.Decode(msg.PubKey))
		case MsgOracleMint:
			copy(pubKey[:], base58.Decode(msg.PubKey))
		case MsgOracleBurn:
			copy(pubKey[:], base58.Decode(msg.PubKey))
		case MsgOracleTransfer:
			copy(pubKey[:], base58.Decode(msg.PubKey))
		default:
			return pubKey, sdk.ErrUnknownRequest("No match for message type.").Result()
		}

		// Check that sender's DID is ledgered
		senderDidDoc, _ := didKeeper.GetDidDoc(ctx, msg.GetSignerDid())
		if senderDidDoc == nil {
			return pubKey, sdk.ErrUnauthorized("Sender did not found").Result()
		}

		return pubKey, sdk.Result{}
	}
}
