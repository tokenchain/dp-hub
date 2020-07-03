package keeper

import (
	types "github.com/tokenchain/ixo-blockchain/x/did"
	common "github.com/tokenchain/ixo-blockchain/x/did/internal/types"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryDidDocs(t *testing.T) {
	ctx, k, cdc := CreateTestInput()
	cdc.RegisterInterface((*types.DidDoc)(nil), nil)
	err := k.SetDidDoc(ctx, &common.ValidDidDoc)
	require.Nil(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(k)
	res, err := querier(ctx, []string{"queryDidDoc", common.ValidDidDoc.Did}, query)
	require.Nil(t, err)

	var a types.BaseDidDoc
	if err := cdc.UnmarshalJSON(res, &a); err != nil {
		t.Log(err)
	}
	_, _ = cdc.MarshalJSONIndent(a, "", " ")
	resD, err := querier(ctx, []string{"queryAllDidDocs"}, query)
	require.Nil(t, err)

	var b []types.BaseDidDoc
	if err := cdc.UnmarshalJSON(resD, &b); err != nil {
		t.Log(err)
	}

	_, _ = cdc.MarshalJSONIndent(b, "", " ")

}
