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
	"golang.org/x/crypto/ed25519"
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
func dxpDidAddress(document string) string {
	return fmt.Sprintf("did:dxp:%s", document)
}

func fromSeedToDid(seed [32]byte) IxoDid {

	publicKeyBytes, privateKeyBytes, err := ed25519.GenerateKey(bytes.NewReader(seed[0:32]))
	if err != nil {
		panic(err)
	}
	publicKey := []byte(publicKeyBytes)
	privateKey := []byte(privateKeyBytes)

	signKey := base58.Encode(privateKey[:32])
	keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))

	sovDid := IxoDid{
		Did:                 dxpDidAddress(base58.Encode(publicKey[:16])),
		VerifyKey:           base58.Encode(publicKey),
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
