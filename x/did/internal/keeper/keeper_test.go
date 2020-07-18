package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/did/test"
	"testing"

	"github.com/stretchr/testify/require"
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
	_, err := k.GetDidDoc(ctx, test.EmptyDid)
	require.NotNil(t, err)

	err = k.SetDidDoc(ctx, &test.ValidDidDoc)
	require.Nil(t, err)

	_, err = k.GetDidDoc(ctx, test.ValidDidDoc.GetDid())
	require.Nil(t, err)
}
