package keeper

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tokenchain/dp-hub/x/project/internal/types"
)

func TestQueryProjectDoc(t *testing.T) {
	ctx, k, cdc, _, _ := CreateTestInput()
	codec.RegisterCrypto(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "", nil)

	require.False(t, k.ProjectDocExists(ctx, types.ValidCreateProjectMsg.GetProjectDid()))
	k.SetProjectDoc(ctx, &types.ValidCreateProjectMsg)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(k)
	res, err := querier(ctx, []string{"queryProjectDoc", types.ValidCreateProjectMsg.ProjectDid}, query)
	require.Nil(t, err)

	emptyRes, err := querier(ctx, []string{"queryProjectDoc", "InvalidProjectDid"}, query)
	require.Nil(t, emptyRes)
	require.NotNil(t, err)

	var projectDoc types.MsgCreateProject
	cdc.MustUnmarshalJSON(res, &projectDoc)
}

func TestQueryProjectAccounts(t *testing.T) {
	ctx, k, cdc, _, _ := CreateTestInput()
	codec.RegisterCrypto(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "", nil)

	require.False(t, k.ProjectDocExists(ctx, types.ValidCreateProjectMsg.GetProjectDid()))
	k.SetProjectDoc(ctx, &types.ValidCreateProjectMsg)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(k)
	_, err := querier(ctx, []string{QueryProjectDoc, types.ValidCreateProjectMsg.ProjectDid}, query)
	require.Nil(t, err)

	account, err := k.CreateNewAccount(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidAccId1)
	require.Nil(t, err)
	k.AddAccountToProjectAccounts(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidAccId1, account)

	res, err := querier(ctx, []string{QueryProjectAccounts, types.ValidCreateProjectMsg.ProjectDid}, query)
	require.Nil(t, err)

	var data interface{}
	require.Nil(t, json.Unmarshal(res, &data))

	accountMap := data.(map[string]interface{})
	_, errRes := json.Marshal(accountMap)
	require.Nil(t, errRes)

	account, err = k.CreateNewAccount(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidAccId1)
	require.NotNil(t, err)
}

func TestQueryTxs(t *testing.T) {
	ctx, k, cdc, _, _ := CreateTestInput()
	codec.RegisterCrypto(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "", nil)

	require.False(t, k.ProjectDocExists(ctx, types.ValidCreateProjectMsg.GetProjectDid()))
	k.SetProjectDoc(ctx, &types.ValidCreateProjectMsg)

	k.AddProjectWithdrawalTransaction(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidWithdrawalInfo)
	k.AddProjectWithdrawalTransaction(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidWithdrawalInfo)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(k)
	res, err := querier(ctx, []string{QueryProjectTx, types.ValidCreateProjectMsg.ProjectDid}, query)
	require.Nil(t, err)

	var txs []types.WithdrawalInfo
	cdc.MustUnmarshalJSON(res, &txs)

	_, err = querier(ctx, []string{QueryProjectTx, "InvalidDid"}, query)
	require.NotNil(t, err)

}
