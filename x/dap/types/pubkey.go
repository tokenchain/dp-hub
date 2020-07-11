package types

import (
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"time"
)

// DidKeeper defines the did contract that must be fulfilled throughout the ixo module

type (
	DidKeeper interface {
		GetDidDoc(ctx sdk.Context, did exported.Did) (exported.DidDoc, error)
		SetDidDoc(ctx sdk.Context, did exported.DidDoc) (err error)
		AddDidDoc(ctx sdk.Context, did exported.DidDoc)
		AddCredentials(ctx sdk.Context, did exported.Did, credential exported.DidCredential) (err error)
		GetAllDidDocs(ctx sdk.Context) (didDocs []exported.DidDoc)
		GetAddDids(ctx sdk.Context) (dids []exported.Did)
	}
	SigVerification struct {
		ak        auth.AccountKeeper
		publicKey crypto.PubKey
		account   authexported.Account
		pgetter   PubKeyGetter
		tx        IxoTx
	}
	DapPubKeyDecoratorDecorator struct {
		SigVerification
	}
	SigVerificationDecorator struct {
		SigVerification
		signature IxoSignature
	}
)

func NewDefaultPubKeyGetter(didKeeper DidKeeper) PubKeyGetter {
	return func(ctx sdk.Context, msg IxoMsg) (pubKey crypto.PubKey, res error) {
		signerDidDoc, err := didKeeper.GetDidDoc(ctx, msg.GetSignerDid())
		if err != nil {
			return pubKey, err
		}
		var pubKeyRaw ed25519tm.PubKeyEd25519
		copy(pubKeyRaw[:], base58.Decode(signerDidDoc.GetPubKey()))
		return pubKeyRaw, nil
	}
}

func NewSigVerification(ak auth.AccountKeeper, p PubKeyGetter) SigVerification {
	return SigVerification{
		ak:      ak,
		pgetter: p,
	}
}
func (sv SigVerification) isGenesis(ctx sdk.Context) bool {
	return ctx.BlockHeight() == 0
}
func (sv SigVerification) retrievePubkey(ctx sdk.Context, tx sdk.Tx, simulate bool) error {

	sigTx, ok := tx.(IxoTx)
	if !ok {
		return InvalidTxDecode()
	}
	sv.tx = sigTx

	// all messages must be of type IxoMsg
	msg, ok := sigTx.GetMsgs()[0].(IxoMsg)
	if !ok {
		//gInfo = sdk.GasInfo{}
		return IntErr("msg must be ixo.IxoMsg. dxp")
	}

	signer := sigTx.GetSigner()
	acc := sv.ak.GetAccount(ctx, signer)
	if acc != nil {
		p := acc.GetPubKey()
		if p != nil {
			sv.publicKey = p
			return nil
		}
	} else {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "account is not found...")
	}

	sv.account = acc
	// Get pubKey

	pubKey, err := sv.pgetter(ctx, msg)
	if simulate {
		// In simulate mode the transaction comes with no signatures, thus if the
		// account's pubkey is nil, both signature verification and gasKVStore.Set()
		// shall consume the largest amount, i.e. it takes more gas to verify
		// secp256k1 keys than ed25519 ones.

		if pubKey == nil {
			sv.publicKey = simSecp256k1Pubkey
		}
	} else {

		if err != nil {
			return err
		}

		sv.publicKey = pubKey
	}

	err = acc.SetPubKey(sv.publicKey)
	if err != nil {
		return InvalidPubKey(err.Error())
	}
	return nil
}

func NewDapPubKeyDecorator(ak auth.AccountKeeper, p PubKeyGetter) DapPubKeyDecoratorDecorator {
	return DapPubKeyDecoratorDecorator{
		NewSigVerification(ak, p),
	}
}
func (__edp DapPubKeyDecoratorDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if e := __edp.retrievePubkey(ctx, tx, simulate); e != nil {
		return ctx, InvalidTxDecodePubkeyNotFound(e)
	}
	return next(ctx, tx, simulate)
}

func NewSigVerificationDecorator(ak auth.AccountKeeper, p PubKeyGetter) SigVerificationDecorator {
	return SigVerificationDecorator{
		SigVerification: NewSigVerification(ak, p),
	}
}
func (sv SigVerificationDecorator) initializeSignature() error {
	if len(sv.tx.GetSignatures()) > 0 {
		sig := sv.tx.GetSignatures()[0]
		//sv.signature = sig
		sv.signature = NewSignature(time.Now(), sig.SignatureValue)
		return nil
	} else {
		return ErrItemNotFound("tx signature not found")
	}
}

// verify the signature and increment the sequence. If the account doesn't have
func (sv SigVerificationDecorator) SignMessage(signBytes []byte, privKey [64]byte) error {
	pk := ed25519tm.PrivKeyEd25519(privKey)
	signatureBytes, err := pk.Sign(signBytes)
	if err != nil {
		return err
	}
	var setSignature [64]byte
	copy(setSignature[:], signatureBytes)
	sv.signature = NewSignature(time.Now(), setSignature)
	return nil
}
func (sv SigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if e := sv.retrievePubkey(ctx, tx, simulate); e != nil {
		return ctx, InvalidTxDecodePubkeyNotFound(e)
	}
	if e := sv.initializeSignature(); e != nil {
		return ctx, e
	}
	signBytes := sv.tx.GetSignBytes(ctx, sv.account)
	if !simulate && !sv.publicKey.VerifyBytes(signBytes, sv.signature.SignatureValue[:]) {
		return ctx, Unauthorized("Signature Verification failed. dxp")
	}
	return next(ctx, tx, simulate)
}
