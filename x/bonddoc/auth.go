package bonddoc

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x"
	types2 "github.com/tokenchain/ixo-blockchain/x/ixo/types"

	"github.com/tokenchain/ixo-blockchain/x/bonddoc/internal/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
)

func GetPubKeyGetter(keeper Keeper) ixo.PubKeyGetter {
	return func(ctx sdk.Context, msg types2.DpMsg) ([32]byte, error) {

		// Get signer PubKey
		var pubKey [32]byte
		switch msg := msg.(type) {
		case types.MsgCreateBond:
			copy(pubKey[:], base58.Decode(msg.GetPubKey()))
		case types.MsgUpdateBondStatus:
			bondDid := msg.GetSignerDid()
			bondDoc, err := keeper.GetBondDoc(ctx, bondDid)
			if err != nil {
				return pubKey, sdkerrors.Wrapf(sdkerrors.ErrNoSignatures,"bond did not is not right %s", bondDid)
			}
			copy(pubKey[:], base58.Decode(bondDoc.GetPubKey()))
		default:
			return pubKey, x.UnknownRequest( "No match for message type.")
		}
		return pubKey, nil
	}
}
