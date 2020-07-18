package ante

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	aexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did"
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
	didexported "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"time"
)

// DidKeeper defines the did contract that must be fulfilled throughout the ixo module

type (
	SigVerification struct {
		ak              auth.AccountKeeper
		stdSignature    auth.StdSignature
		account_address []byte
		pubkey          []byte
		dap_tx          did.IxoTx
		dk              didexported.DidKeeper
		signature       did.IxoSignature
	}
	DapPubKeyDecoratorDecorator struct {
		SigVerification
	}
	SigVerificationDecorator struct {
		SigVerification
		signature did.IxoSignature
	}
)

func NewSigVerification(ak auth.AccountKeeper, p didexported.DidKeeper) SigVerification {
	return SigVerification{
		ak: ak,
		dk: p,
	}
}
func (sv SigVerification) GetAccountKeeper() auth.AccountKeeper {
	return sv.ak
}
func (sv SigVerification) GetDapTx() did.IxoTx {
	return sv.dap_tx
}
func (sv SigVerification) isGenesis(ctx sdk.Context) bool {
	return ctx.BlockHeight() == 0
}
func (sv SigVerification) GetSignerAccount(ctx sdk.Context) aexported.Account {
	signerAcc, err := auth.GetSignerAcc(ctx, sv.ak, sv.account_address)
	if err != nil {
		panic(" cannot get the GetSignerAccont")
	}
	return signerAcc
}
func (sv SigVerification) RefreshAccount(ctx sdk.Context) error {
	//address := sdk.AccAddress(sv.publicKey.Address())
	//_, err := auth.GetSignerAcc(ctx, sv.ak, address)
	//if err != nil {
	//		return IntErr("msg must be ixo.IxoMsg. dxp")
	//	}
	//	sv.account_address = address
	return nil
}
func (sv SigVerification) initializeSignatures() (nsv SigVerification, err error) {
	//f := sv.tx.GetFirstSignatureValues()
	if len(sv.dap_tx.Signatures) > 0 {
		sv.signature = did.NewSignature(time.Now(), sv.dap_tx.GetFirstSignature())
		return sv, nil
	} else {
		return sv, ErrItemNotFound("tx signature not found")
	}
}

func (sv SigVerification) pubkeyExtract(ctx sdk.Context, msg did.IxoMsg) (pubKey crypto.PubKey, err error) {
	// Get signer PubKey
	var pubKeyEd25519 ed25519tm.PubKeyEd25519
	switch msg := msg.(type) {
	case did.MsgAddDid:
		fmt.Println("add did doc message extract the pubkey .... ")
		copy(pubKeyEd25519[:], base58.Decode(msg.DidDoc.PubKey))
		//pubKeyEd25519 = did.RecoverDidToEd25519PubKey(msg.DidDoc.)
	default:
		// For the remaining messages, the did is the signer
		fmt.Println("get did doc from other subjects .... ")
		didDoc, _ := sv.dk.GetDidDoc(ctx, msg.GetSignerDid())
		if didDoc == nil {
			return pubKey, Unauthorized("Issuer did not found")
		}
		copy(pubKeyEd25519[:], base58.Decode(didDoc.GetPubKey()))

	}

	fmt.Println("all data extracted .... ")
	return pubKeyEd25519, nil
}
func (sv SigVerification) RetrievePubkey(ctx sdk.Context, tx sdk.Tx, simulate bool) (nsv SigVerification, pubKey crypto.PubKey, err error) {
	sigTx, ok := did.CastTypeSdkTx(tx)

	if !ok {
		return sv, nil, InvalidTxDecode()
	}
	fmt.Println("debug ===============SigVerification.RetrievePubkey, sigTx")
	fmt.Println(sigTx)
	//fmt.Println(*sv.tx)
	// all messages must be of type IxoMsg
	msg, ok := sigTx.GetMsgs()[0].(did.IxoMsg)
	if !ok {
		return sv, nil, IntErr("msg must be ixo.IxoMsg. dxp")
	}

	pubKey, err = sv.pubkeyExtract(ctx, msg)
	if err != nil {
		return sv, nil, err
	}
	fmt.Println("debug ===============SigVerification.RetrievePubkey, pgetter => pubKey")
	fmt.Println(pubKey)
	address := sdk.AccAddress(pubKey.Address())
	signerAcc, err := auth.GetSignerAcc(ctx, sv.ak, address)
	//signer := sigTx.GetSigner()
	//acc := sv.ak.GetAccount(ctx, signer)
	fmt.Println("debug ===============SigVerification.RetrievePubkey, address")
	fmt.Println(address)

	if signerAcc != nil {
		sv.account_address = address
		//copy(sv.account_address, address.Bytes())
		//fmt.Println("check-RetrievePubkey check value pass ....")
		//fmt.Println(address)

	} else {
		return sv, nil, UnknownAddress("the signer account address is not found.")
	}

	if simulate {
		// In simulate mode the transaction comes with no signatures, thus if the
		// account's pubkey is nil, both signature verification and gasKVStore.Set()
		// shall consume the largest amount, i.e. it takes more gas to verify
		// secp256k1 keys than ed25519 ones.

		if pubKey == nil {
			//copy(sv.publicKey[:], simSecp256k1Pubkey[:])
			return sv, nil, InvalidPubKey("there is no valid public key to use for simulation.")
		}
	}

	/*	err = signerAcc.SetPubKey(sv.publicKey)
		if err != nil {
			return InvalidPubKey(err.Error())
		}
	*/
	simSig := types.StdSignature{
		Signature: simSecp256k1Sig[:],
		PubKey:    pubKey,
	}

	sv.stdSignature = simSig
	sv.dap_tx = sigTx
	//var pubkey_orginal ed25519tm.PubKeyEd25519
	//copy(pubkey_orginal[:], pubKey.Bytes()[5:])
	//sv.pubkey = pubkey_orginal[:]
	//fmt.Println(sv.account_address)
	//sv.pubkey = ed25519tm.PubKeyEd25519(pubKey)

	return sv, pubKey, nil
}

func NewDapPubKeyDecorator(ak auth.AccountKeeper, p didexported.DidKeeper) DapPubKeyDecoratorDecorator {
	return DapPubKeyDecoratorDecorator{
		NewSigVerification(ak, p),
	}
}
func (__edp DapPubKeyDecoratorDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if _, _, e := __edp.RetrievePubkey(ctx, tx, simulate); e != nil {
		return ctx, InvalidTxDecodePubkeyNotFound(e)
	}
	fmt.Println("✅  DapPubKeyDecoratorDecorator pass")
	return next(ctx, tx, simulate)
}

func NewSigVerificationDecorator(ak auth.AccountKeeper, p didexported.DidKeeper) SigVerificationDecorator {
	return SigVerificationDecorator{
		SigVerification: NewSigVerification(ak, p),
	}
}

func (sv SigVerificationDecorator) VerifyNow(pub []byte, message []byte, sign []byte) error {
	if len(sign) != ed25519.SignatureSize {
		return Unauthorizedf("signature size is not matched, expected size %d got %d !", ed25519.SignatureSize, len(sign))
	}
	if l := len(pub); l != ed25519.PublicKeySize {
		return Unauthorizedf("ed25519: bad public key length expected %d but got %d! ", ed25519.PublicKeySize, l)
	}
	/*
		fmt.Println("===> debug public key check:", base58.Encode(pub), len(pub), pub)
		fmt.Println("===> ⚠️ check signed message data ....")
		fmt.Println(message)
	*/
	if ed25519.Verify(pub, message, sign) {
		return nil
	} else {
		return Unauthorized("Signature Verification failed. dxp Z.")
	}
}

func (sv SigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	nsv, pp, e := sv.RetrievePubkey(ctx, tx, simulate)
	if e != nil {
		return ctx, InvalidTxDecodePubkeyNotFound(e)
	}
	nsv2, e := nsv.initializeSignatures()
	if e != nil {
		return ctx, e
	}

	signedMessageBytes := nsv2.dap_tx.GetSignBytes(ctx, nsv2.GetSignerAccount(ctx))
	/*

		fmt.Println("✅  check signature data ....")
		fmt.Println(nsv2.signature.SignatureValue[:])

	*/

	/*if !simulate && !pp.VerifyBytes(signedMessageBytes, nsv2.signature.SignatureValue[:]) {
		return ctx, Unauthorized("Signature Verification failed. dxp")
	}*/
	if !simulate {
		key := pp.(ed25519tm.PubKeyEd25519)
		if er := sv.VerifyNow(key[:], signedMessageBytes, nsv2.signature.SignatureValue[:]); er != nil {
			return ctx, er
		}
		fmt.Println("✅  SigVerificationDecorator pass ....")
	}

	acc := nsv2.GetSignerAccount(ctx)
	// increment account sequence
	if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
		return ctx, InvalidTxDecodeMsg(err.Error())
	}

	sv.ak.SetAccount(ctx, nsv2.GetSignerAccount(ctx))
	return next(ctx, tx, simulate)
}
