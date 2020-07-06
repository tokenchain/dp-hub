package treasury

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"

	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/dap"
)

func GetPubKeyGetter(didKeeper did.Keeper) dap.PubKeyGetter {
	return func(ctx sdk.Context, msg types.IxoMsg) ([32]byte, error) {

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
			return pubKey, x.UnknownRequest("No match for message type.")
		}

		// Check that sender's DID is ledgered
		senderDidDoc, _ := didKeeper.GetDidDoc(ctx, msg.GetSignerDid())
		if senderDidDoc == nil {
			return pubKey, x.Unauthorized("Sender did not found")
		}

		return pubKey, nil
	}
}
