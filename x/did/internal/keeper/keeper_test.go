package keeper

import (
	"github.com/tokenchain/ixo-blockchain/x/did"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

func TestKeeper(t *testing.T) {
	ctx, k, cdc := CreateTestInput()
	cdc.RegisterInterface((*did.DidDoc)(nil), nil)
	_, err := k.GetDidDoc(ctx, types.EmptyDid)
	require.NotNil(t, err)

	err = k.SetDidDoc(ctx, &types.ValidDidDoc)
	require.Nil(t, err)

	_, err = k.GetDidDoc(ctx, types.ValidDidDoc.GetDid())
	require.Nil(t, err)
}
