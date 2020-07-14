package project

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	er "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tokenchain/ixo-blockchain/x"
	 "github.com/tokenchain/ixo-blockchain/x/dap"
	types2 "github.com/tokenchain/ixo-blockchain/x/dap/types"
	"github.com/tokenchain/ixo-blockchain/x/did/ante"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/payments"
)

const (
	IxoAccountFeesId               InternalAccountID = "IxoFees"
	IxoAccountPayFeesId            InternalAccountID = "IxoPayFees"
	InitiatingNodeAccountPayFeesId InternalAccountID = "InitiatingNodePayFees"
	ValidatingNodeSetAccountFeesId InternalAccountID = "ValidatingNodeSetFees"
)

func NewHandler(k Keeper, fk payments.Keeper, bk bank.Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateProject:
			return handleMsgCreateProject(ctx, k, msg)
		case MsgCreateAgent:
			return handleMsgCreateAgent(ctx, k, bk, msg)
		case MsgUpdateAgent:
			return handleMsgUpdateAgent(ctx, k, bk, msg)
		case MsgCreateClaim:
			return handleMsgCreateClaim(ctx, k, fk, bk, msg)
		case MsgCreateEvaluation:
			return handleMsgCreateEvaluation(ctx, k, fk, bk, msg)
		case MsgWithdrawFunds:
			return handleMsgWithdrawFunds(ctx, k, bk, msg)
		case MsgUpdateProjectStatus:
			return handleMsgUpdateProjectStatus(ctx, k, bk, msg)
		default:
			return nil, x.UnknownRequest("No match for message type.")
		}
	}
}

func handleMsgCreateProject(ctx sdk.Context, k Keeper, msg MsgCreateProject) (*sdk.Result, error) {

	projectDid := msg.GetProjectDid()

	if _, err := createAccountInProjectAccounts(ctx, k, projectDid, IxoAccountFeesId); err != nil {
		return nil, err
	}
	if _, err := createAccountInProjectAccounts(ctx, k, projectDid, IxoAccountPayFeesId); err != nil {
		return nil, err
	}
	if _, err := createAccountInProjectAccounts(ctx, k, projectDid, InitiatingNodeAccountPayFeesId); err != nil {
		return nil, err
	}
	if _, err := createAccountInProjectAccounts(ctx, k, projectDid, ValidatingNodeSetAccountFeesId); err != nil {
		return nil, err
	}
	if _, err := createAccountInProjectAccounts(ctx, k, projectDid, InternalAccountID(msg.GetProjectDid())); err != nil {
		return nil, err
	}

	if k.ProjectDocExists(ctx, msg.GetProjectDid()) {
		return nil, x.ErrInvalidDid("Project already exists")
	}
	k.SetProjectDoc(ctx, &msg)
	k.SetProjectWithdrawalTransactions(ctx, msg.GetProjectDid(), nil)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUpdateProjectStatus(ctx sdk.Context, k Keeper, bk bank.Keeper,
	msg MsgUpdateProjectStatus) (result *sdk.Result, res error) {

	existingProjectDoc, err := getProjectDoc(ctx, k, msg.ProjectDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Project")
	}

	newStatus := msg.Data.Status
	if !newStatus.IsValidProgressionFrom(existingProjectDoc.GetStatus()) {
		return nil, x.UnknownRequest("Invalid Status Progression requested")
	}

	if newStatus == FundedStatus {
		projectAddr, err := getProjectAccount(ctx, k, existingProjectDoc.GetProjectDid())
		if err != nil {
			return nil, err
		}

		projectAcc := k.AccountKeeper.GetAccount(ctx, projectAddr)
		if projectAcc == nil {
			return nil, x.UnknownRequest("Could not find project account")
		}

		minimumFunding := k.GetParams(ctx).ProjectMinimumInitialFunding
		if projectAcc.GetCoins().AmountOf(dap.IxoNativeToken).LT(minimumFunding) {
			return nil, er.Wrapf(er.ErrInsufficientFunds, "Project has not reached minimum funding %s", minimumFunding)
		}
	}

	if newStatus == PaidoutStatus {
		result, err = payoutFees(ctx, k, bk, existingProjectDoc.GetProjectDid())
	}

	existingProjectDoc.SetStatus(newStatus)
	_, _ = k.UpdateProjectDoc(ctx, existingProjectDoc)

	return result, err
}

func payoutFees(ctx sdk.Context, k Keeper, bk bank.Keeper, projectDid exported.Did) (*sdk.Result, error) {
	var allevents []sdk.Event

	event, err := payAllFeesToAddress(ctx, k, bk, projectDid, IxoAccountPayFeesId, IxoAccountFeesId)
	if err != nil {
		return nil, x.ErrInvalidDid("Failed to send coins")
	} else {
		allevents = append(allevents, event...)
	}

	event, err = payAllFeesToAddress(ctx, k, bk, projectDid, InitiatingNodeAccountPayFeesId, IxoAccountFeesId)
	if err != nil {
		return nil, x.ErrInvalidDid("Failed to send coins")
	} else {
		allevents = append(allevents, event...)
	}

	event, err = payAllFeesToAddress(ctx, k, bk, projectDid, ValidatingNodeSetAccountFeesId, IxoAccountFeesId)
	if err != nil {
		return nil, x.ErrInvalidDid("Failed to send coins")
	} else {
		allevents = append(allevents, event...)
	}

	ixoDid := k.GetParams(ctx).IxoDid
	amount := getIxoAmount(ctx, k, bk, projectDid, IxoAccountFeesId)
	err = payoutAndRecon(ctx, k, bk, projectDid, IxoAccountFeesId, ixoDid, amount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: allevents}, err
}

func payAllFeesToAddress(ctx sdk.Context, k Keeper, bk bank.Keeper, projectDid exported.Did,
	sendingAddress InternalAccountID, receivingAddress InternalAccountID) (sdk.Events, error) {
	feesToPay := getIxoAmount(ctx, k, bk, projectDid, sendingAddress)

	if feesToPay.Amount.LT(sdk.ZeroInt()) {
		return nil, x.ErrInvalidDid("Negative fee to pay")
	}
	if feesToPay.Amount.IsZero() {
		return nil, nil
	}

	receivingAccount, err := getAccountInProjectAccounts(ctx, k, projectDid, receivingAddress)
	if err != nil {
		return sdk.Events{}, err
	}

	sendingAccount, _ := getAccountInProjectAccounts(ctx, k, projectDid, sendingAddress)

	return sdk.Events{}, bk.SendCoins(ctx, sendingAccount, receivingAccount, sdk.Coins{feesToPay})
}

func getIxoAmount(ctx sdk.Context, k Keeper, bk bank.Keeper, projectDid exported.Did, accountID InternalAccountID) sdk.Coin {
	found := checkAccountInProjectAccounts(ctx, k, projectDid, accountID)
	if found {
		accAddr, _ := getAccountInProjectAccounts(ctx, k, projectDid, accountID)
		coins := bk.GetCoins(ctx, accAddr)
		return sdk.NewCoin(dap.IxoNativeToken, coins.AmountOf(dap.IxoNativeToken))
	}
	return sdk.NewCoin(dap.IxoNativeToken, sdk.ZeroInt())
}

func handleMsgCreateAgent(ctx sdk.Context, k Keeper, bk bank.Keeper, msg MsgCreateAgent) (*sdk.Result, error) {
	// Check if project exists
	_, err := getProjectDoc(ctx, k, msg.ProjectDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Project")
	}
	// Create account in project accounts for the agent
	_, err = createAccountInProjectAccounts(ctx, k, msg.ProjectDid, InternalAccountID(msg.Data.AgentDid))
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUpdateAgent(ctx sdk.Context, k Keeper, bk bank.Keeper, msg MsgUpdateAgent) (*sdk.Result, error) {
	// Check if project exists
	_, err := getProjectDoc(ctx, k, msg.ProjectDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Project")
	}
	// TODO: implement agent update (or remove functionality)
	return &sdk.Result{}, nil
}

func handleMsgCreateClaim(ctx sdk.Context, k Keeper, fk payments.Keeper, bk bank.Keeper, msg MsgCreateClaim) (*sdk.Result, error) {
	// Check if project exists
	_, err := getProjectDoc(ctx, k, msg.ProjectDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Project")
	}
	// Process claim fees
	err = processFees(ctx, k, fk, bk, payments.FeeClaimTransaction, msg.ProjectDid)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateEvaluation(ctx sdk.Context, k Keeper, fk payments.Keeper, bk bank.Keeper, msg MsgCreateEvaluation) (*sdk.Result, error) {

	// Check if project exists
	projectDoc, err := getProjectDoc(ctx, k, msg.ProjectDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Project")
	}

	// Process evaluation fees
	err = processFees(
		ctx, k, fk, bk, payments.FeeEvaluationTransaction, msg.ProjectDid)
	if err != nil {
		return nil, err
	}

	// Process evaluator pay
	err = processEvaluatorPay(ctx, k, fk, bk, msg.ProjectDid,
		msg.SenderDid, projectDoc.GetEvaluatorPay())
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleMsgWithdrawFunds(ctx sdk.Context, k Keeper, bk bank.Keeper,
	msg MsgWithdrawFunds) (*sdk.Result, error) {

	withdrawFundsDoc := msg.Data
	projectDoc, err := getProjectDoc(ctx, k, withdrawFundsDoc.ProjectDid)
	if err != nil {
		return nil, x.UnknownRequest("Could not find Project")
	}

	if projectDoc.GetStatus() != PaidoutStatus {
		return nil, x.UnknownRequest("Project not in PAIDOUT Status")
	}

	projectDid := withdrawFundsDoc.ProjectDid
	recipientDid := withdrawFundsDoc.RecipientDid
	amount := withdrawFundsDoc.Amount

	// If this is a refund, recipient has to be the project creator
	if withdrawFundsDoc.IsRefund && (recipientDid != projectDoc.GetSenderDid()) {
		return nil, x.UnknownRequest("Only project creator can get a refund")
	}

	var fromAccountId InternalAccountID
	if withdrawFundsDoc.IsRefund {
		fromAccountId = InternalAccountID(projectDid)
	} else {
		fromAccountId = InternalAccountID(recipientDid)
	}

	amountCoin := sdk.NewCoin(dap.IxoNativeToken, amount)
	err = payoutAndRecon(ctx, k, bk, projectDid, fromAccountId, recipientDid, amountCoin)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func payoutAndRecon(ctx sdk.Context, k Keeper, bk bank.Keeper, projectDid exported.Did,
	fromAccountId InternalAccountID, recipientDid exported.Did, amount sdk.Coin) error {

	ixoBalance := getIxoAmount(ctx, k, bk, projectDid, fromAccountId)
	if ixoBalance.IsLT(amount) {
		return x.ErrInvalidDid("insufficient funds in specified account")
	}

	fromAccount, err := getAccountInProjectAccounts(ctx, k, projectDid, fromAccountId)
	if err != nil {
		return err
	}

	recipientAddr := ante.StringToAddr(recipientDid)
	err = bk.SendCoins(ctx, fromAccount, recipientAddr, sdk.Coins{amount})
	if err != nil {
		return err
	}

	var actionId [32]byte
	dec := sdk.OneDec() // TODO: should increment with each withdrawal (ref: #113)
	copy(actionId[:], dec.Bytes())

	addProjectWithdrawalTransaction(ctx, k, projectDid, actionId, recipientDid, amount)
	return nil
}

func getProjectDoc(ctx sdk.Context, k Keeper, projectDid exported.Did) (StoredProjectDoc, error) {
	ixoProjectDoc, err := k.GetProjectDoc(ctx, projectDid)
	if err != nil {
		return nil, err
	}

	return ixoProjectDoc.(StoredProjectDoc), nil
}

func processFees(ctx sdk.Context, k Keeper, fk payments.Keeper, bk bank.Keeper, feeType payments.FeeType, projectDid exported.Did) error {
	projectAddr, _ := getProjectAccount(ctx, k, projectDid)

	validatingNodeSetAddr, err := getAccountInProjectAccounts(ctx, k, projectDid, ValidatingNodeSetAccountFeesId)
	if err != nil {
		return err
	}

	ixoAddr, err := getAccountInProjectAccounts(ctx, k, projectDid, IxoAccountFeesId)
	if err != nil {
		return err
	}

	ixoFactor := fk.GetParams(ctx).IxoFactor
	nodePercentage := fk.GetParams(ctx).NodeFeePercentage

	var adjustedFeeAmount sdk.Dec
	switch feeType {
	case payments.FeeClaimTransaction:
		adjustedFeeAmount = fk.GetParams(ctx).ClaimFeeAmount.Mul(ixoFactor)
	case payments.FeeEvaluationTransaction:
		adjustedFeeAmount = fk.GetParams(ctx).EvaluationFeeAmount.Mul(ixoFactor)
	default:
		return x.UnknownRequest("Invalid Fee type.")
	}

	nodeAmount := adjustedFeeAmount.Mul(nodePercentage).RoundInt64()
	ixoAmount := adjustedFeeAmount.RoundInt64() - nodeAmount

	err = bk.SendCoins(ctx, projectAddr, validatingNodeSetAddr, sdk.Coins{sdk.NewInt64Coin(dap.IxoNativeToken, nodeAmount)})
	if err != nil {
		return err
	}

	err = bk.SendCoins(ctx, projectAddr, ixoAddr, sdk.Coins{sdk.NewInt64Coin(dap.IxoNativeToken, ixoAmount)})
	if err != nil {
		return err
	}

	return nil
}

func processEvaluatorPay(ctx sdk.Context, k Keeper, fk payments.Keeper,
	bk bank.Keeper, projectDid, senderDid exported.Did, evaluatorPay int64) error {

	if evaluatorPay == 0 {
		return nil
	}

	projectAddr, _ := getAccountInProjectAccounts(ctx, k, projectDid, InternalAccountID(projectDid))
	evaluatorAccAddr, _ := getAccountInProjectAccounts(ctx, k, projectDid, InternalAccountID(senderDid))

	nodeAddr, err := getAccountInProjectAccounts(ctx, k, projectDid, InitiatingNodeAccountPayFeesId)
	if err != nil {
		return err
	}

	ixoAddr, err := getAccountInProjectAccounts(ctx, k, projectDid, IxoAccountPayFeesId)
	if err != nil {
		return err
	}

	feePercentage := fk.GetParams(ctx).EvaluationPayFeePercentage
	nodeFeePercentage := fk.GetParams(ctx).EvaluationPayNodeFeePercentage

	totalEvaluatorPayAmount := sdk.NewDec(evaluatorPay).Mul(types2.IxoDecimals) // This is in IXO * 10^8
	evaluatorPayFeeAmount := totalEvaluatorPayAmount.Mul(feePercentage)
	evaluatorPayLessFees := totalEvaluatorPayAmount.Sub(evaluatorPayFeeAmount)
	nodePayFees := evaluatorPayFeeAmount.Mul(nodeFeePercentage)
	ixoPayFees := evaluatorPayFeeAmount.Sub(nodePayFees)

	err = bk.SendCoins(ctx, projectAddr, evaluatorAccAddr, sdk.Coins{sdk.NewInt64Coin(dap.IxoNativeToken, evaluatorPayLessFees.RoundInt64())})
	if err != nil {
		return err
	}

	err = bk.SendCoins(ctx, projectAddr, nodeAddr, sdk.Coins{sdk.NewInt64Coin(dap.IxoNativeToken, nodePayFees.RoundInt64())})
	if err != nil {
		return err
	}

	err = bk.SendCoins(ctx, projectAddr, ixoAddr, sdk.Coins{sdk.NewInt64Coin(dap.IxoNativeToken, ixoPayFees.RoundInt64())})
	if err != nil {
		return err
	}

	return nil
}

func checkAccountInProjectAccounts(ctx sdk.Context, k Keeper, projectDid exported.Did,
	accountId InternalAccountID) bool {
	accMap := k.GetAccountMap(ctx, projectDid)
	_, found := accMap[accountId]

	return found
}

func addProjectWithdrawalTransaction(ctx sdk.Context, k Keeper, projectDid exported.Did,
	actionID [32]byte, recipientDid exported.Did, amount sdk.Coin) {
	actionIDStr := "0x" + hex.EncodeToString(actionID[:])

	withdrawalInfo := WithdrawalInfo{
		ActionID:     actionIDStr,
		ProjectDid:   projectDid,
		RecipientDid: recipientDid,
		Amount:       amount,
	}

	k.AddProjectWithdrawalTransaction(ctx, projectDid, withdrawalInfo)
}

func createAccountInProjectAccounts(ctx sdk.Context, k Keeper, projectDid exported.Did, accountId InternalAccountID) (sdk.AccAddress, error) {
	acc, err := k.CreateNewAccount(ctx, projectDid, accountId)
	if err != nil {
		return nil, err
	}

	k.AddAccountToProjectAccounts(ctx, projectDid, accountId, acc)

	return acc.GetAddress(), nil
}

func getAccountInProjectAccounts(ctx sdk.Context, k Keeper, projectDid exported.Did,
	accountId InternalAccountID) (sdk.AccAddress, error) {
	accMap := k.GetAccountMap(ctx, projectDid)

	addr, found := accMap[accountId]
	if found {
		return addr, nil
	} else {
		return createAccountInProjectAccounts(ctx, k, projectDid, accountId)
	}
}

func getProjectAccount(ctx sdk.Context, k Keeper, projectDid exported.Did) (sdk.AccAddress, error) {
	return getAccountInProjectAccounts(ctx, k, projectDid, InternalAccountID(projectDid))
}
