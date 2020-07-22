package exported

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/pkg/errors"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	edgen "github.com/tokenchain/ixo-blockchain/x/did/ed25519"
	"strings"
	"unsafe"

	naclBox "golang.org/x/crypto/nacl/box"
	"io"
)

type SovrinSecret struct {
	Seed                 string `json:"seed" yaml:"seed"`
	SignKey              string `json:"signKey" yaml:"signKey"`
	EncryptionPrivateKey string `json:"encryptionPrivateKey" yaml:"encryptionPrivateKey"`
}

func (ss SovrinSecret) String() string {
	output, err := json.MarshalIndent(ss, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

type SovrinDid struct {
	Did                 string       `json:"did" yaml:"did"`
	VerifyKey           string       `json:"verifyKey" yaml:"verifyKey"`
	EncryptionPublicKey string       `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	Secret              SovrinSecret `json:"secret" yaml:"secret"`
}

func NewDapDid(did, verifykey, publickey, seed, signkey, privatekey, dpaddress, dppubkey, name string) IxoDid {
	return IxoDid{
		Did:                 did,
		VerifyKey:           verifykey,
		EncryptionPublicKey: publickey,
		Secret:              NewSecret(seed, signkey, privatekey),
		Dpinfo: DpInfo{
			DpAddress: dpaddress,
			PubKey:    dppubkey,
			Name:      name,
			Algo:      "ed25519",
		},
	}
}

// todo: not working and will be removed
func GenDidInfoExperiment(doc keys.Info, privateKey tmcrypto.PrivKey, x keys.SigningAlgo) IxoDid {

	_, privateKeyBytes, err := edgen.GenerateKey(bytes.NewReader(doc.GetPubKey().Bytes()[0:32]))
	publicKeyBytes2, _, err := edgen.GenerateKey(bytes.NewReader(privateKeyBytes[:]))
	if err != nil {
		panic(err)
	}
	//signKey := base58.Encode(privateKeyBytes[:32])
	hashedEntropy := sha256.Sum256(privateKey.Bytes())
	dpaddress := doc.GetAddress().String()

	privKey := PrivateKeyToSecp256k1(privateKey)

	sovDid := IxoDid{
		Did:                 dxpDidAddress(base58.Encode(doc.GetPubKey().Bytes()[:16])),
		VerifyKey:           base58.Encode([]byte(dpaddress)),
		EncryptionPublicKey: base58.Encode(publicKeyBytes2[:]),

		Secret: Secret{
			Seed:                 hex.EncodeToString(hashedEntropy[:]),
			SignKey:              strings.ToUpper(hex.EncodeToString(privKey[24:])),
			EncryptionPrivateKey: strings.ToUpper(hex.EncodeToString(privKey[:24])),
		},
	}

	//	addr, err := sdk.AccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t")

	return sovDid

}
func makePubKey(bt *[32]byte) (pubKey tmcrypto.PubKey) {
	var pubKeyRaw ed25519tm.PubKeyEd25519
	copy(pubKeyRaw[:], bt[:])
	return pubKeyRaw
}
func InfoToDidEd25519(docName string, derivedPriv []byte, debug bool) IxoDid {
	pubEd25519, priEd25519, _ := edgen.GenerateKey(bytes.NewReader(derivedPriv[0:32]))
	keyPairPublicKey, keyPairPrivateKey, _ := naclBox.GenerateKey(bytes.NewReader(priEd25519[:]))

	if debug {
		fmt.Println("========pair1 key  =========")
		fmt.Println("pub", len(pubEd25519), pubEd25519)
		fmt.Println("pri", len(priEd25519), priEd25519)
		fmt.Println("========pair2 key  =========")
		fmt.Println("pub", len(keyPairPublicKey), keyPairPublicKey)
		fmt.Println("pri", len(keyPairPrivateKey), keyPairPrivateKey)
		fmt.Println("========Derived private key  =========")
		fmt.Println(len(derivedPriv), derivedPriv)
		fmt.Println("================================")
	}
	sovDid := IxoDid{
		Dpinfo: DpInfo{
			DpAddress: VerifyKeyToAddrEd25519(base58.Encode(pubEd25519[:])).String(),
			PubKey:    sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, makePubKey(pubEd25519)),
			Name:      docName,
			Algo:      "ed25519",
			//DpAddress: doc.GetAddress().String(),
			//Algo:      "secp256k1",
		},
		Did:                 dxpDidAddress(base58.Encode(pubEd25519[:16])),
		VerifyKey:           base58.Encode(pubEd25519[:]),
		EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),
		Secret: Secret{
			Seed:                 hex.EncodeToString(derivedPriv[0:32]),
			SignKey:              base58.Encode(priEd25519[:32]),
			EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
		},
	}

	return sovDid
}

func (sd SovrinDid) String() string {
	output, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

func generateMnemonic() string {
	entropy, _ := bip39.NewEntropy(12)
	mnemonicWords, _ := bip39.NewMnemonic(entropy)
	return mnemonicWords
}

func fromJsonString(jsonSovrinDid string) (IxoDid, error) {
	var did IxoDid
	err := json.Unmarshal([]byte(jsonSovrinDid), &did)
	if err != nil {
		err := fmt.Errorf("Could not unmarshal did into struct. Error: %s", err.Error())
		return IxoDid{}, err
	}

	return did, nil
}

func mnemonicToDid(mnemonic string) IxoDid {
	seed := sha256.New()
	seed.Write([]byte(mnemonic))
	var seed32 [32]byte
	copy(seed32[:], seed.Sum(nil)[:32])
	return fromSeedToDid(seed32)
}

func MnToDid(mnemonic string, usr string) IxoDid {
	return mnemonicToDid(mnemonic)
}
func SeedToDid(seed []byte) IxoDid {
	var seed32 [32]byte
	copy(seed32[:], seed[:32])
	return fromSeedToDid(seed32)
}
func UnverifiedToAddr(ver string) sdk.AccAddress {
	return sdk.AccAddress(tmcrypto.AddressHash([]byte(ver)))
}
func UnmarshalDxpDid(jsonSovrinDid string) (IxoDid, error) {
	return fromJsonStringDp(jsonSovrinDid)
}
func BytesToString(data []byte) string {
	return string(data[:])
}
func BytesToStringUnsafe(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

/*
func VerifyKeyToAddrEd25519(verifyKey string) sdk.AccAddress {
	var pubKey ed25519.PublicKey
	copy(pubKey[:], base58.Decode(verifyKey))
	return sdk.AccAddress(pubKey)
}
*/
func VerifyKeyToAddrEd25519(verifyKey string) sdk.AccAddress {
	var pubKey ed25519tm.PubKeyEd25519
	copy(pubKey[:], base58.Decode(verifyKey))
	return sdk.AccAddress(pubKey.Address())
}
func VerifyKeyToPublicKeyEd25519(verifyKey string) tmcrypto.PubKey {
	var pubKey ed25519tm.PubKeyEd25519
	copy(pubKey[:], base58.Decode(verifyKey))
	return ed25519tm.PubKeyEd25519(pubKey)
}

func VerifyKeyToPublicKeyEd25519Mech32(verifyKey string) string {
	return sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, VerifyKeyToPublicKeyEd25519(verifyKey))
}
func VerifyKeyToAddr(verifyKey string) sdk.AccAddress {
	//var privkeyL []byte
	code := base58.Decode(verifyKey)
	//copy(privkeyL, code)
	t := BytesToString(code)
	//hash := tmcrypto.AddressHash(code)
	g, e := sdk.AccAddressFromBech32(t)
	//println(verifyKey)
	//println("Testing at this line")
	if e != nil {
		msg := fmt.Sprintf("cannot verify this key %s. ", verifyKey)
		panic(msg)
	}
	return g
}

func dxpDidAddress(document string) string {
	return fmt.Sprintf("%s:%s", DidPrefix, document)
}

func fromSeedToDid(seed [32]byte) IxoDid {
	publicKeyBytes, privateKeyBytes, err := edgen.GenerateKey(bytes.NewReader(seed[0:32]))
	if err != nil {
		panic(err)
	}
	//head part
	signKey := base58.Encode(privateKeyBytes[:32])
	//keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))
	keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKeyBytes[:]))

	sovDid := IxoDid{
		Did:                 dxpDidAddress(base58.Encode(publicKeyBytes[:16])),
		VerifyKey:           base58.Encode(publicKeyBytes[:]),
		EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),

		Secret: Secret{
			Seed:                 hex.EncodeToString(seed[0:32]),
			SignKey:              signKey,
			EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
		},
	}

	return sovDid
}

func SignMessage(message []byte, signKey string, verifyKey string) []byte {
	// Force the length to 64
	privateKey := make([]byte, edgen.PrivateKeySize)
	fullPrivKey := edgen.PrivateKey(privateKey)
	copy(fullPrivKey[:], getArrayFromKey(signKey))
	copy(fullPrivKey[32:], getArrayFromKey(verifyKey))

	return edgen.Sign(fullPrivKey, message)
}

func VerifySignedMessage(message []byte, signature []byte, verifyKey string) bool {
	publicKey := edgen.PublicKey{}
	copy(publicKey[:], getArrayFromKey(verifyKey))
	result := edgen.Verify(publicKey, message, signature)

	return result
}

func SignMessageDid(message []byte, did_doc IxoDid) []byte {
	var recover_privKey secp256k1.PrivKeySecp256k1
	p1, _ := hex.DecodeString(strings.ToLower(did_doc.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(did_doc.Secret.SignKey))
	copy(recover_privKey[:], p1)
	copy(recover_privKey[24:], p2)
	return recover_privKey[:]
}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)
	if start < 0 || end > length || start > end {
		return ""
	}
	if start == 0 && end == length {
		return source
	}
	return string(r[start:end])
}

func GetNonce() [24]byte {
	var nonce [24]byte
	if _, err := io.ReadFull(cryptoRand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}
func AddAccountEd25519ByDid(kb keys.Keybase, name string, doc IxoDid) error {
	accpub := VerifyKeyToPublicKeyEd25519Mech32(doc.VerifyKey)
	return AddAccountEd25519(kb, name, accpub)
}
func AddAccountEd25519(kb keys.Keybase, name string, pubkey string) error {
	_, err := kb.Get(name)
	if err == nil {
		//account exist
		return errors.Wrap(nil, "account exist")
	}
	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pubkey)
	if err != nil {
		//account exist
		return err
	}
	algo := keys.Ed25519

	_, err = kb.CreateOffline(name, pk, algo)
	if err != nil {
		return err
	}
	return nil
}

/*
func GetKeyPairFromSignKey(signKey string) ([32]byte, [32]byte) {
	publicKey, privateKey, err := naclbox.GenerateKey(bytes.NewReader(getArrayFromKey(signKey)))
	if err != nil {
		panic(err)
	}
	return *publicKey, *privateKey
}
*/
