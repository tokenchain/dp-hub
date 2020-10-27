package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/stretchr/testify/require"

	"github.com/tokenchain/dp-hub/x/project/internal/types"
)

func TestProjectDoc(t *testing.T) {
	ctx, k, _, _, _ := CreateTestInput()

	require.False(t, k.ProjectDocExists(ctx, types.ValidCreateProjectMsg.GetProjectDid()))
	k.SetProjectDoc(ctx, &types.ValidCreateProjectMsg)

	doc, err := k.GetProjectDoc(ctx, types.ValidCreateProjectMsg.ProjectDid)
	require.Nil(t, err)
	require.Equal(t, &types.ValidCreateProjectMsg, doc)

	resUpdated, err := k.UpdateProjectDoc(ctx, &types.ValidUpdateProjectMsg)
	require.Nil(t, err)

	expected, err := k.GetProjectDoc(ctx, types.ValidUpdateProjectMsg.ProjectDid)
	require.Equal(t, resUpdated, expected)

	_, err = k.GetProjectDoc(ctx, "Invalid Did")
	require.NotNil(t, err)
}

func TestKeeperAccountMap(t *testing.T) {
	ctx, k, cdc, _, _ := CreateTestInput()
	codec.RegisterCrypto(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "", nil)

	account, err := k.CreateNewAccount(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidAccId1)
	require.Nil(t, err)

	k.AddAccountToProjectAccounts(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidAccId1, account)

	accountMap := k.GetAccountMap(ctx, types.ValidCreateProjectMsg.ProjectDid)
	_, found := accountMap[types.ValidAccId1]
	require.True(t, found)

	account, err = k.CreateNewAccount(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidAccId1)
	require.NotNil(t, err)

}

func TestKeeperWithdrawalInfo(t *testing.T) {
	ctx, k, cdc, _, _ := CreateTestInput()
	codec.RegisterCrypto(cdc)

	withdrawals, err := k.GetProjectWithdrawalTransactions(ctx, "")
	require.NotNil(t, err)
	require.Equal(t, 0, len(withdrawals))

	k.AddProjectWithdrawalTransaction(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidWithdrawalInfo)
	k.AddProjectWithdrawalTransaction(ctx, types.ValidCreateProjectMsg.ProjectDid, types.ValidWithdrawalInfo)

	withdrawals, err = k.GetProjectWithdrawalTransactions(ctx, types.ValidCreateProjectMsg.ProjectDid)
	require.Nil(t, err)
	require.Equal(t, 2, len(withdrawals))
}
