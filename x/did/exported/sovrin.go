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
	"github.com/tokenchain/ixo-blockchain/x/did/ed25519"
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

func (sd SovrinDid) String() string {
	output, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

func GenerateMnemonic() string {
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

func MnToDid(mnemonic string) IxoDid {
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
func InfoToDid(doc keys.Info, privateKey tmcrypto.PrivKey, x keys.SigningAlgo) IxoDid {

	_, privateKeyBytes, err := ed25519.GenerateKey(bytes.NewReader(doc.GetPubKey().Bytes()[0:32]))
	publicKeyBytes2, _, err := ed25519.GenerateKey(bytes.NewReader(privateKeyBytes[:]))
	if err != nil {
		panic(err)
	}
	//signKey := base58.Encode(privateKeyBytes[:32])
	hashedEntropy := sha256.Sum256(privateKey.Bytes())
	dpaddress := doc.GetAddress().String()

	//var privKey tmcryptoed25519.PrivKeyEd25519

	/*
		var privKey ed25519tm.PrivKeyEd25519
		copy(privKey[:], base58.Decode(ixoDid.Secret.SignKey))
		copy(privKey[32:], base58.Decode(ixoDid.VerifyKey))
	*/

	fmt.Println("private to bytes length =======")

	fmt.Println("algo type =======")
	fmt.Println(x)

	privKey := PrivateKeyToSecp256k1(privateKey)

	fmt.Println("byte length after marshal =======")
	fmt.Println(cap(privKey.Bytes()))
	fmt.Println(privKey.Bytes())

	fmt.Println(cap(privateKey.Bytes()))
	fmt.Println(privateKey.Bytes())

	fmt.Println("is the same", privateKey.Equals(privKey))

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

func InfoToDidEd25519(doc keys.Info, privateKey tmcrypto.PrivKey, derivedPriv []byte) IxoDid {

	pub, pri, err := ed25519.GenerateKey(bytes.NewReader(derivedPriv))
	if err != nil {
		panic(err)
	}
	//signKey := base58.Encode(privateKeyBytes[:32])
	hashedEntropy := sha256.Sum256(derivedPriv)
	dpaddress := doc.GetAddress().String()

	//var privKey tmcryptoed25519.PrivKeyEd25519

	/*
		var privKey ed25519tm.PrivKeyEd25519
		copy(privKey[:], base58.Decode(ixoDid.Secret.SignKey))
		copy(privKey[32:], base58.Decode(ixoDid.VerifyKey))
	*/

	//fmt.Println(">>> derivedPriv type, length =======")
	//fmt.Println(derivedPriv, len(derivedPriv))

	privKey := PrivateKeyToSecp256k1(privateKey)

	fmt.Println(">>> byte length after marshal =======")
	fmt.Println(len(privKey.Bytes()))
	fmt.Println(privKey.Bytes())

	fmt.Println(len(privateKey.Bytes()))
	fmt.Println(privateKey.Bytes())

	fmt.Println(">>> ed25519 keypair =======")
	fmt.Println("public", len(pub), pub)
	fmt.Println("private", len(pri), pri)

	//fmt.Println("is the same", privateKey.Equals(privKey))
	//48+32 = 80
	//Part1 := len(privateKey.Bytes())+len(pri[:32])
	//Part2 := len(pri[:32])
	//	allcap:=len(privateKey.Bytes())+len(pri[:32])
	var privateKeyFinal [37 + 32]byte
	copy(privateKeyFinal[:], privateKey.Bytes())
	copy(privateKeyFinal[37:], pri[:32])

	fmt.Println(">>> ed25519 keypair =======")
	fmt.Println(">>> secret private key", len(privateKeyFinal), privateKeyFinal)
	fmt.Println(">>> ed25519 final =======")

	// privateKeyFinal = EncryptionPrivateKey[37:] + SignKey[:]

	// revert_pri = privateKeyFinal[64:] +
	sovDid := IxoDid{
		Did:                 dxpDidAddress(base58.Encode(pub[:16])),
		VerifyKey:           base58.Encode([]byte(dpaddress)),
		EncryptionPublicKey: base58.Encode(pub[16:]),

		Secret: Secret{
			Seed:                 hex.EncodeToString(hashedEntropy[:]),
			SignKey:              strings.ToUpper(hex.EncodeToString(pri[32:])),
			EncryptionPrivateKey: strings.ToUpper(hex.EncodeToString(privateKeyFinal[:])),
		},
	}

	//	addr, err := sdk.AccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t")
	return sovDid

}

func dxpDidAddress(document string) string {
	return fmt.Sprintf("did:dxp:%s", document)
}

func fromSeedToDid(seed [32]byte) IxoDid {
	publicKeyBytes, privateKeyBytes, err := ed25519.GenerateKey(bytes.NewReader(seed[0:32]))
	if err != nil {
		panic(err)
	}
	//publicKey := []byte(publicKeyBytes)
	//privateKey := []byte(privateKeyBytes)
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

/*
func Gen() IxoDid {
	var seed [32]byte
	if _, err := io.ReadFull(cryptoRand.Reader, seed[:]); err != nil {
		panic(err)
	}
	did, _ := fromJsonString(seed)
	return did
}
*/
func SignMessage(message []byte, signKey string, verifyKey string) []byte {
	// Force the length to 64
	privateKey := make([]byte, ed25519.PrivateKeySize)
	fullPrivKey := ed25519.PrivateKey(privateKey)
	copy(fullPrivKey[:], getArrayFromKey(signKey))
	copy(fullPrivKey[32:], getArrayFromKey(verifyKey))

	return ed25519.Sign(fullPrivKey, message)
}

func VerifySignedMessage(message []byte, signature []byte, verifyKey string) bool {
	publicKey := ed25519.PublicKey{}
	copy(publicKey[:], getArrayFromKey(verifyKey))
	result := ed25519.Verify(publicKey, message, signature)

	return result
}

func SignMessageDid(message []byte, did_doc IxoDid) []byte {
	var recover_privKey secp256k1.PrivKeySecp256k1
	p1, _ := hex.DecodeString(strings.ToLower(did_doc.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(did_doc.Secret.SignKey))
	copy(recover_privKey[:], p1)
	copy(recover_privKey[24:], p2)
	//return ed25519.Sign(recover_privKey, message)
	return recover_privKey[:]
}
func RecoverDidToPrivateKeyClassic(did_doc IxoDid) ed25519tm.PrivKeyEd25519 {
	var privKey ed25519tm.PrivKeyEd25519
	copy(privKey[:], base58.Decode(did_doc.Secret.SignKey))
	copy(privKey[32:], base58.Decode(did_doc.VerifyKey))
	return privKey
}
func RecoverDidEd25519ToPrivateKey(did_ed_doc IxoDid) [64]byte {
	var recover_priv_key_ed [64]byte
	p1, _ := hex.DecodeString(strings.ToLower(did_ed_doc.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(did_ed_doc.Secret.SignKey))
	copy(recover_priv_key_ed[:], p1[37:])
	copy(recover_priv_key_ed[32:], p2)
	return recover_priv_key_ed
}

func RecoverDidSecpK1ToPrivateKey(did_secp_doc IxoDid) [32]byte {
	var recover_privKey secp256k1.PrivKeySecp256k1
	p1, _ := hex.DecodeString(strings.ToLower(did_secp_doc.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(did_secp_doc.Secret.SignKey))
	copy(recover_privKey[:], p1)
	copy(recover_privKey[24:], p2)
	return secp256k1.PrivKeySecp256k1(recover_privKey)
}

func GetNonce() [24]byte {
	var nonce [24]byte
	if _, err := io.ReadFull(cryptoRand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func getArrayFromKey(key string) []byte {
	return base58.Decode(key)
}

func GetKeyPairFromSignKey(signKey string) ([32]byte, [32]byte) {
	publicKey, privateKey, err := naclBox.GenerateKey(bytes.NewReader(getArrayFromKey(signKey)))
	if err != nil {
		panic(err)
	}
	return *publicKey, *privateKey
}
func PrivateKeyToSecp256k1(privKey tmcrypto.PrivKey) secp256k1.PrivKeySecp256k1 {
	var privKey_orginal secp256k1.PrivKeySecp256k1
	copy(privKey_orginal[:], privKey.Bytes()[5:])
	return privKey_orginal
}
func SecpPrivKey(bz []byte) secp256k1.PrivKeySecp256k1 {
	var bzArr [32]byte
	copy(bzArr[:], bz)
	return secp256k1.PrivKeySecp256k1(bzArr)
}
func PrivateKeyToEd25519(privKey tmcrypto.PrivKey) ed25519.PrivateKey {
	var privKey_orginal ed25519.PrivateKey
	copy(privKey_orginal[:], privKey.Bytes()[:])
	return privKey_orginal
}
func AddAccount(kb keys.Keybase, name string, pubkey string) error {
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
