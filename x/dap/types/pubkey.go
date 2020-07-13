package types

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	aexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
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
		ak              auth.AccountKeeper
		stdSignature    auth.StdSignature
		account_address []byte
		pubkey          []byte
		dap_tx          IxoTx
		pgetter         PubKeyGetter
		signature       IxoSignature
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
func initializeAddress() sdk.AccAddress {
	var addressDefault [32]byte
	entropySeed, _ := bip39.NewEntropy(64)
	p, _, _ := ed25519.GenerateKey(bytes.NewReader(entropySeed))
	copy(addressDefault[:], p[:])
	return sdk.AccAddress(addressDefault[:])
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
func (sv SigVerification) GetSignerAccont(ctx sdk.Context) aexported.Account {
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
		sv.signature = NewSignature(time.Now(), sv.dap_tx.GetFirstSignatureValues())
		return sv, nil
	} else {
		return sv, ErrItemNotFound("tx signature not found")
	}
}

func (sv SigVerification) RetrievePubkey(ctx sdk.Context, tx sdk.Tx, simulate bool) (nsv SigVerification, pubKey crypto.PubKey, err error) {
	sigTx, ok := tx.(IxoTx)
	if !ok {
		return sv, nil, InvalidTxDecode()
	}
	//fmt.Println(sigTx)
	//fmt.Println(*sv.tx)
	// all messages must be of type IxoMsg
	msg, ok := sigTx.GetMsgs()[0].(IxoMsg)
	if !ok {
		//gInfo = sdk.GasInfo{}
		return sv, nil, IntErr("msg must be ixo.IxoMsg. dxp")
	}

	pubKey, err = sv.pgetter(ctx, msg)
	if err != nil {
		return sv, nil, err
	}

	address := sdk.AccAddress(pubKey.Address())
	signerAcc, err := auth.GetSignerAcc(ctx, sv.ak, address)
	//signer := sigTx.GetSigner()
	//acc := sv.ak.GetAccount(ctx, signer)
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
	var pubkey_orginal ed25519tm.PubKeyEd25519
	copy(pubkey_orginal[:], pubKey.Bytes()[5:])
	sv.pubkey = pubkey_orginal[:]
	//fmt.Println(sv.account_address)

	return sv, pubKey, nil
}

func NewDapPubKeyDecorator(ak auth.AccountKeeper, p PubKeyGetter) DapPubKeyDecoratorDecorator {
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

func NewSigVerificationDecorator(ak auth.AccountKeeper, p PubKeyGetter) SigVerificationDecorator {
	return SigVerificationDecorator{
		SigVerification: NewSigVerification(ak, p),
	}
}

// verify the signature and increment the sequence. If the account doesn't have
func (sv SigVerificationDecorator) SignMessage(signBytes []byte, privKey [64]byte) error {
	pk := ed25519tm.PrivKeyEd25519(privKey)
	signatureBytes, err := pk.Sign(signBytes)
	if err != nil {
		return err
	}
	//var setSignature [64]byte
	//copy(setSignature[:], signatureBytes)
	sv.signature = NewSignature(time.Now(), signatureBytes)
	return nil
}
func (sv SigVerificationDecorator) VerifyNow(pub []byte, message []byte, sign []byte) error {
	if len(sign) != ed25519.SignatureSize {
		return Unauthorizedf("signature size is not matched, expected size %d got %d !", ed25519.SignatureSize, len(sign))
	}

	if l := len(pub); l != ed25519.PublicKeySize {
		return Unauthorizedf("ed25519: bad public key length expected %d but got %d! ", ed25519.PublicKeySize, l)
	}

	if ed25519.Verify(pub, message, sign) {
		return nil
	} else {
		return Unauthorized("Signature Verification failed. dxp")
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
	signedMessageBytes := nsv2.dap_tx.GetSignBytes(ctx, nsv2.GetSignerAccont(ctx))
	fmt.Println("✅  check signature data ....")
	fmt.Println(nsv2.signature.SignatureValue[:])
	fmt.Println("✅  check signed message data ....")
	fmt.Println(signedMessageBytes)

	if !simulate && !pp.VerifyBytes(signedMessageBytes, nsv2.signature.SignatureValue[:]) {
		return ctx, Unauthorized("Signature Verification failed. dxp")
	}
	if !simulate {
		key := pp.(ed25519tm.PubKeyEd25519)
		if er := sv.VerifyNow(key[:], signedMessageBytes, nsv2.signature.SignatureValue[:]); er != nil {
			return ctx, er
		}
		fmt.Println("✅  SigVerificationDecorator pass ....")
	}

	return next(ctx, tx, simulate)
}
