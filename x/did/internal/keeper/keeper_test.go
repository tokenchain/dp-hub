package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

type (
	Did    = string
	DidDoc interface {
		SetDid(did Did) error
		GetDid() Did
		SetPubKey(pubkey string) error
		GetPubKey() string
		Address() sdk.AccAddress
		AddressUnverified() sdk.AccAddress
	}
)

func TestKeeper(t *testing.T) {
	ctx, k, cdc := CreateTestInput()
	cdc.RegisterInterface((*DidDoc)(nil), nil)
	_, err := k.GetDidDoc(ctx, types.EmptyDid)
	require.NotNil(t, err)

	err = k.SetDidDoc(ctx, &types.ValidDidDoc)
	require.Nil(t, err)

	_, err = k.GetDidDoc(ctx, types.ValidDidDoc.GetDid())
	require.Nil(t, err)
}
