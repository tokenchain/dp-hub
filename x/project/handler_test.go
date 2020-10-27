package project

import (
	types2 "github.com/tokenchain/dp-hub/x/dap/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/stretchr/testify/require"

	"github.com/tokenchain/dp-hub/x/project/internal/keeper"
	"github.com/tokenchain/dp-hub/x/project/internal/types"
)

func TestHandler_CreateClaim(t *testing.T) {
	ctx, k, cdc, paymentsKeeper, bankKeeper := keeper.CreateTestInput()
	codec.RegisterCrypto(cdc)
	cdc.RegisterConcrete(types.MsgCreateProject{}, "project/CreateProject", nil)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/Account", nil)
	params := paymentsKeeper.GetParams(ctx)
	params.IxoFactor = sdk.OneDec()
	params.NodeFeePercentage = sdk.ZeroDec()
	params.ClaimFeeAmount = sdk.NewDec(6).Quo(sdk.NewDec(10)).Mul(types2.IxoDecimals)
	paymentsKeeper.SetParams(ctx, params)
	projectMsg := types.MsgCreateClaim{
		ProjectDid: "6iftm1hHdaU6LJGKayRMev",
		TxHash:     "txHash",
		SenderDid:  "senderDid",
		Data:       types.CreateClaimDoc{ClaimID: "claim1"},
	}
	res, _ := handleMsgCreateClaim(ctx, k, paymentsKeeper, bankKeeper, projectMsg)
	require.NotNil(t, res)
}

func TestHandler_ProjectMsg(t *testing.T) {
	ctx, k, cdc, _, _ := keeper.CreateTestInput()
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/Account", nil)

	res, err := handleMsgCreateProject(ctx, k, types.ValidCreateProjectMsg)
	require.NoError(t, err)
	res, err = handleMsgCreateProject(ctx, k, types.ValidCreateProjectMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
}
func Test_CreateEvaluation(t *testing.T) {
	ctx, k, cdc, fk, bk := keeper.CreateTestInput()

	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/Account", nil)

	params := fk.GetParams(ctx)
	params.IxoFactor = sdk.OneDec()
	params.NodeFeePercentage = sdk.NewDec(5).Quo(sdk.NewDec(10))
	params.ClaimFeeAmount = sdk.NewDec(6).Quo(sdk.NewDec(10)).Mul(types2.IxoDecimals)
	params.EvaluationFeeAmount = sdk.NewDec(4).Quo(sdk.NewDec(10)).Mul(types2.IxoDecimals) // 0.4
	params.EvaluationPayFeePercentage = sdk.ZeroDec()
	params.EvaluationPayNodeFeePercentage = sdk.NewDec(5).Quo(sdk.NewDec(10))
	fk.SetParams(ctx, params)

	evaluationMsg := types.MsgCreateEvaluation{
		TxHash:     "txHash",
		SenderDid:  "senderDid",
		ProjectDid: "6iftm1hHdaU6LJGKayRMev",
		Data: types.CreateEvaluationDoc{
			ClaimID: "claim1",
			Status:  types.PendingClaim,
		},
	}

	msg := types.MsgCreateProject{
		TxHash:     "",
		SenderDid:  "",
		ProjectDid: "6iftm1hHdaU6LJGKayRMev",
		PubKey:     "47mm6LCDAyJmqkbUbqGoZKZkBixjBgvDFRMwQRF9HWMU",
		Data: types.ProjectDoc{
			NodeDid:              "Tu2QWRHuDufywDALbBQ2r",
			RequiredClaims:       "requireClaims1",
			EvaluatorPayPerClaim: "10",
			ServiceEndpoint:      "https://togo.pds.ixo.network",
			CreatedOn:            "2018-05-21T15:53:18.484Z",
			CreatedBy:            "6Fu7FbbGoCJ8tX3vMMCss9",
			Status:               "CREATED",
		},
	}

	_, err := createAccountInProjectAccounts(ctx, k, msg.GetProjectDid(), IxoAccountFeesId)
	require.Nil(t, err)
	acc, err := createAccountInProjectAccounts(ctx, k, msg.GetProjectDid(), InternalAccountID(msg.GetProjectDid()))
	require.Nil(t, err)
	require.NotNil(t, acc)
	require.False(t, k.ProjectDocExists(ctx, msg.GetProjectDid()))
	k.SetProjectDoc(ctx, &msg)

	res, err := handleMsgCreateEvaluation(ctx, k, fk, bk, evaluationMsg)

	require.NoError(t, err)
	require.NotNil(t, res)

}

func Test_WithdrawFunds(t *testing.T) {
	ctx, k, cdc, _, bk := keeper.CreateTestInput()
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/Account", nil)

	msg := types.MsgWithdrawFunds{
		SenderDid: "6iftm1hHdaU6LJGKayRMev",
		Data: types.WithdrawFundsDoc{
			ProjectDid:   "6iftm1hHdaU6LJGKayRMev",
			RecipientDid: "6iftm1hHdaU6LJGKayRMev",
			Amount:       sdk.NewInt(100),
			IsRefund:     true,
		},
	}

	msg1 := types.MsgCreateProject{
		TxHash:     "",
		SenderDid:  "",
		ProjectDid: "6iftm1hHdaU6LJGKayRMev",
		PubKey:     "47mm6LCDAyJmqkbUbqGoZKZkBixjBgvDFRMwQRF9HWMU",
		Data: types.ProjectDoc{
			NodeDid:              "Tu2QWRHuDufywDALbBQ2r",
			RequiredClaims:       "requireClaims1",
			EvaluatorPayPerClaim: "10",
			ServiceEndpoint:      "https://togo.pds.ixo.network",
			CreatedOn:            "2018-05-21T15:53:18.484Z",
			CreatedBy:            "6Fu7FbbGoCJ8tX3vMMCss9",
			Status:               "PAIDOUT",
		},
	}

	_, err := createAccountInProjectAccounts(ctx, k, msg1.GetProjectDid(), IxoAccountFeesId)
	require.Nil(t, err)

	account, errf := createAccountInProjectAccounts(ctx, k, msg1.GetProjectDid(), InternalAccountID(msg1.GetProjectDid()))
	//	require.Nil(t, err)
	require.NoError(t, errf)
	require.NotNil(t, account)

	require.False(t, k.ProjectDocExists(ctx, msg1.GetProjectDid()))
	k.SetProjectDoc(ctx, &msg1)

	res, err := handleMsgWithdrawFunds(ctx, k, bk, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

}
