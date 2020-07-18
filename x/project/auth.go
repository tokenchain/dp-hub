package project

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/dap"
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
	"github.com/tokenchain/ixo-blockchain/x/did"
	aute2 "github.com/tokenchain/ixo-blockchain/x/did/ante"
	export2 "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"time"
)

/*
func GetPubKeyGetter(keeper Keeper, didKeeper did.Keeper) dap.PubKeyGetter {
	return func(ctx sdk.Context, msg dap.IxoMsg) (crypto.PubKey, error) {

		// Get signer PubKey
		var pubKey crypto.PubKey
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
}*/

func GetPubKeyGetter(keeper Keeper, didKeeper did.Keeper) aute2.PubKeyGetter {
	return func(ctx sdk.Context, msg aute2.IxoMsg) (pubKey crypto.PubKey, res error) {

		// Get signer PubKey
		var pubKeyEd25519 ed25519.PubKeyEd25519
		switch msg := msg.(type) {
		case MsgCreateProject:
			copy(pubKeyEd25519[:], base58.Decode(msg.GetPubKey()))
		case MsgWithdrawFunds:
			signerDid := msg.GetSignerDid()
			signerDoc, _ := didKeeper.GetDidDoc(ctx, signerDid)
			if signerDoc == nil {
				return pubKey, Unauthorized("signer did not found")
			}
			copy(pubKeyEd25519[:], base58.Decode(signerDoc.GetPubKey()))
		default:
			// For the remaining messages, the project is the signer
			projectDoc, err := keeper.GetProjectDoc(ctx, msg.GetSignerDid())
			if err != nil {
				return pubKey, IntErr("project did not found")
			}
			copy(pubKeyEd25519[:], base58.Decode(projectDoc.GetPubKey()))
		}
		return pubKeyEd25519, nil
	}
}

func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}

func IntErr(m string) error {
	return errors.Wrap(errors.ErrPanic, m)
}

// Identical to Cosmos DeductFees function, but tokens sent to project account
func deductProjectFundingFees(bankKeeper bank.Keeper, ctx sdk.Context, acc exported.Account, projectDid export2.Did, fees sdk.Coins) error {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return errors.Wrapf(errors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return errors.Wrapf(errors.ErrInsufficientFunds, "insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return errors.Wrapf(errors.ErrInsufficientFunds, "insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	projectAddr := aute2.DidToAddr(projectDid)
	err := bankKeeper.SendCoins(ctx, acc.GetAddress(), projectAddr, fees)
	if err != nil {
		return err
	}

	return nil
}

func getProjectCreationSignBytes(chainID string, ixoTx aute2.IxoTx, acc exported.Account, genesis bool) []byte {
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
	bk bank.Keeper, didKeeper did.Keeper,
	pubKeyGetter aute2.PubKeyGetter) sdk.AnteHandler {
	//} func xNewProjectCreationAnteHandler(ak auth.AccountKeeper, sk supply.Keeper, bk bank.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res error) {

		if addr := sk.GetModuleAddress(auth.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", auth.FeeCollectorName))
		}
		// var _ sdk.Tx = dap.IxoTx{}
		// all transactions must be of type ixo.IxoTx
		ixoTx, ok := tx.(aute2.IxoTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = auth.SetGasMeter(simulate, ctx, 0)
			return newCtx, export2.IntErr("tx must be ixo.IxoTx")
		}

		params := ak.GetParams(ctx)

		// Project creation uses an infinite gas meter
		newCtx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err
		}

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")

		// require that long memos get rejected
		vmd := ante.NewValidateMemoDecorator(ak)
		antehandler := sdk.ChainAnteDecorators(vmd)
		_, err := antehandler(ctx, tx, false)
		if err != nil {
			return newCtx, err
		}
		/*
			if res := tx.ValidateMemo(dap.IxoTx{Memo: ixoTx.Memo}, params); res != nil {
				return newCtx, res
			}*/

		// message must be of type MsgCreateProject
		msg, ok := ixoTx.GetMsgs()[0].(MsgCreateProject)
		if !ok {
			return newCtx, export2.IntErr("msg must be MsgCreateProject")
		}

		// Fetch signer (project itself). Account expected to not exist
		signerAddr := ixoTx.GetSigner()
		signerAcc, res := auth.GetSignerAcc(newCtx, ak, signerAddr)
		if res != nil {
			return newCtx, export2.IntErr("expected project account to not exist")
		}

		// confirm that fee is the exact amount expected
		expectedTotalFee := sdk.NewCoins(sdk.NewCoin(types.NativeToken, sdk.NewInt(MsgCreateProjectFee)))
		if !ixoTx.Fee.Amount.IsEqual(expectedTotalFee) {
			return newCtx, export2.ErrInvalidCoins("invalid fee")
		}

		// Calculate transaction fee and project funding
		transactionFee := sdk.NewCoins(sdk.NewCoin(types.NativeToken, sdk.NewInt(MsgCreateProjectTransactionFee)))
		projectFunding := expectedTotalFee.Sub(transactionFee) // panics if negative result

		// deduct the fees
		if !ixoTx.Fee.Amount.IsZero() {
			// fetch fee payer
			feePayerAddr := aute2.DidToAddr(msg.SenderDid)
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

		// Get project pubKey
		/*projectPubKey, res := pubKeyGetter(ctx, msg)
		if res != nil {
			return newCtx, res
		}*/
		// Fetch signer account (project itself); create if it does not exist
		signerAcc, res = auth.GetSignerAcc(ctx, ak, signerAddr)
		if res != nil {
			signerAcc = ak.NewAccountWithAddress(ctx, signerAddr)
			ak.SetAccount(ctx, signerAcc)
		}

		// check signature, return account with incremented nonce

		// check signature, return account with incremented nonce
		//ixoSig := auth.StdSignature{PubKey: projectPubKey, Signature: ixoTx.GetSignatures()[0].SignatureValue[:]}
		ixoSig := aute2.NewSignature(time.Now(), ixoTx.GetSignatures()[0].SignatureValue)
		isGenesis := ctx.BlockHeight() == 0
		signBytes := getProjectCreationSignBytes(newCtx.ChainID(), ixoTx, signerAcc, isGenesis)
		signerAcc, res = dap.ProcessSig(ctx, signerAcc, ixoSig, signBytes, simulate, ak.GetParams(ctx))
		if res != nil {
			return newCtx, res
		}

		ak.SetAccount(newCtx, signerAcc)
		//		sdk.Result{GasWanted: ixoTx.Fee.Gas}
		return newCtx, nil // continue...
	}
}
