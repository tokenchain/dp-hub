package project

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	er "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tokenchain/ixo-blockchain/x"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/ixo"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

func GetPubKeyGetter(keeper Keeper, didKeeper did.Keeper) ixo.PubKeyGetter {
	return func(ctx sdk.Context, msg types.IxoMsg) ([32]byte, error) {

		// Get signer PubKey
		var pubKey [32]byte
		switch msg := msg.(type) {
		case MsgCreateProject:
			copy(pubKey[:], base58.Decode(msg.GetPubKey()))
		case MsgUpdateProjectStatus:
			projectDoc, err := keeper.GetProjectDoc(ctx, msg.ProjectDid)
			if err != nil {
				return pubKey, er.Wrap(x.ErrInternalE, "project did not found")
			}
			copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
		case MsgCreateAgent:
			projectDoc, err := keeper.GetProjectDoc(ctx, msg.ProjectDid)
			if err != nil {
				return pubKey, x.IntErr("project did not found")
			}
			copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
		case MsgUpdateAgent:
			projectDoc, err := keeper.GetProjectDoc(ctx, msg.ProjectDid)
			if err != nil {
				return pubKey, x.IntErr("project did not found")
			}
			copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
		case MsgCreateClaim:
			projectDoc, err := keeper.GetProjectDoc(ctx, msg.ProjectDid)
			if err != nil {
				return pubKey, x.IntErr("project did not found")
			}
			copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
		case MsgCreateEvaluation:
			projectDoc, err := keeper.GetProjectDoc(ctx, msg.ProjectDid)
			if err != nil {
				return pubKey, x.IntErr("project did not found")
			}
			copy(pubKey[:], base58.Decode(projectDoc.GetPubKey()))
		case MsgWithdrawFunds:
			didDoc, _ := didKeeper.GetDidDoc(ctx, msg.Data.RecipientDid)
			if didDoc == nil {
				return pubKey, x.Unauthorized("signer did not found")
			}
			copy(pubKey[:], base58.Decode(didDoc.GetPubKey()))
		default:
			return pubKey, x.UnknownRequest("No match for message type.")
		}
		return pubKey, nil
	}
}

// Identical to Cosmos DeductFees function, but tokens sent to project account
func deductProjectFundingFees(bankKeeper bank.Keeper, ctx sdk.Context, acc exported.Account, projectDid types.Did, fees sdk.Coins) error {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return er.Wrapf(er.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return er.Wrapf(er.ErrInsufficientFunds, "insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return er.Wrapf(er.ErrInsufficientFunds, "insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	projectAddr := types.DidToAddr(projectDid)
	err := bankKeeper.SendCoins(ctx, acc.GetAddress(), projectAddr, fees)
	if err != nil {
		return err
	}

	return nil
}

func getProjectCreationSignBytes(chainID string, ixoTx types.IxoTx, acc exported.Account, genesis bool) []byte {
	var accNum uint64
	if !genesis {
		// Fixed account number used so that sign bytes do not depend on it
		accNum = uint64(0)
	}

	return auth.StdSignBytes(
		chainID, accNum, acc.GetSequence(), ixoTx.Fee, ixoTx.Msgs, ixoTx.Memo,
	)
}

func NewProjectCreationAnteHandler(ak auth.AccountKeeper, sk supply.Keeper,
	bk bank.Keeper, pubKeyGetter ixo.PubKeyGetter) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res error) {

		if addr := sk.GetModuleAddress(auth.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", auth.FeeCollectorName))
		}

		// all transactions must be of type ixo.IxoTx
		ixoTx, ok := tx.(types.IxoTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = auth.SetGasMeter(simulate, ctx, 0)
			return newCtx, x.IntErr("tx must be ixo.IxoTx")
		}

		params := ak.GetParams(ctx)

		// Project creation uses an infinite gas meter
		newCtx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err
		}

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")

		if res := ixo.ValidateMemo(types.IxoTx{Memo: ixoTx.Memo}, params); res != nil {
			return newCtx, res
		}

		// message must be of type MsgCreateProject
		msg, ok := ixoTx.GetMsgs()[0].(MsgCreateProject)
		if !ok {
			return newCtx, x.IntErr("msg must be MsgCreateProject")
		}

		// Fetch signer (project itself). Account expected to not exist
		signerAddr := ixoTx.GetSigner()
		signerAcc, res := auth.GetSignerAcc(newCtx, ak, signerAddr)
		if res != nil {
			return newCtx, x.IntErr("expected project account to not exist")
		}

		// confirm that fee is the exact amount expected
		expectedTotalFee := sdk.NewCoins(sdk.NewCoin(types.IxoNativeToken, sdk.NewInt(MsgCreateProjectFee)))
		if !ixoTx.Fee.Amount.IsEqual(expectedTotalFee) {
			return newCtx, x.ErrInvalidCoins("invalid fee")
		}

		// Calculate transaction fee and project funding
		transactionFee := sdk.NewCoins(sdk.NewCoin(types.IxoNativeToken, sdk.NewInt(MsgCreateProjectTransactionFee)))
		projectFunding := expectedTotalFee.Sub(transactionFee) // panics if negative result

		// deduct the fees
		if !ixoTx.Fee.Amount.IsZero() {
			// fetch fee payer
			feePayerAddr := types.DidToAddr(msg.SenderDid)
			feePayerAcc, res := auth.GetSignerAcc(ctx, ak, feePayerAddr)
			if res != nil {
				return newCtx, res
			}

			res = auth.DeductFees(sk, newCtx, feePayerAcc, transactionFee)
			if res != nil {
				return newCtx, res
			}

			res = deductProjectFundingFees(bk, newCtx, feePayerAcc, msg.ProjectDid, projectFunding)
			if res != nil {
				return newCtx, res
			}

			// reload the account as fees have been deducted
			feePayerAcc = ak.GetAccount(newCtx, feePayerAcc.GetAddress())
		}

		// Get pubKey
		pubKey, res := pubKeyGetter(ctx, msg)
		if res != nil {
			return newCtx, res
		}

		// Fetch signer account (project itself); create if it does not exist
		signerAcc, res = auth.GetSignerAcc(ctx, ak, signerAddr)
		if res != nil {
			signerAcc = ak.NewAccountWithAddress(ctx, signerAddr)
			ak.SetAccount(ctx, signerAcc)
		}

		// check signature, return account with incremented nonce
		ixoSig := ixoTx.GetSignatures()[0]
		isGenesis := ctx.BlockHeight() == 0
		signBytes := getProjectCreationSignBytes(newCtx.ChainID(), ixoTx, signerAcc, isGenesis)
		signerAcc, res = ixo.ProcessSig(newCtx, signerAcc, signBytes, pubKey, ixoSig, simulate, params)
		if res != nil {
			return newCtx, res
		}

		ak.SetAccount(newCtx, signerAcc)
		//		sdk.Result{GasWanted: ixoTx.Fee.Gas}
		return newCtx, nil // continue...
	}
}
